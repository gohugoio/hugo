// Copyright 2024 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package warpc

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gohugoio/hugo/common/hugio"
	"golang.org/x/sync/errgroup"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/experimental"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

const currentVersion = 1

//go:embed wasm/quickjs.wasm
var quickjsWasm []byte

// Header is in both the request and response.
type Header struct {
	// Major version of the protocol.
	Version uint16 `json:"version"`

	// Unique ID for the request.
	// Note that this only needs to be unique within the current request set time window.
	ID uint32 `json:"id"`

	// Set in the response if there was an error.
	Err string `json:"err"`
}

type Message[T any] struct {
	Header Header `json:"header"`
	Data   T      `json:"data"`
}

func (m Message[T]) GetID() uint32 {
	return m.Header.ID
}

type Dispatcher[Q, R any] interface {
	Execute(ctx context.Context, q Message[Q]) (Message[R], error)
	Close() error
}

func (p *dispatcherPool[Q, R]) getDispatcher() *dispatcher[Q, R] {
	i := int(p.counter.Add(1)) % len(p.dispatchers)
	return p.dispatchers[i]
}

func (p *dispatcherPool[Q, R]) Close() error {
	return p.close()
}

type dispatcher[Q, R any] struct {
	zero Message[R]

	mu    sync.RWMutex
	encMu sync.Mutex

	pending map[uint32]*call[Q, R]

	inOut *inOut

	shutdown bool
	closing  bool
}

type inOut struct {
	sync.Mutex
	stdin  hugio.ReadWriteCloser
	stdout hugio.ReadWriteCloser
	dec    *json.Decoder
	enc    *json.Encoder
}

var ErrShutdown = fmt.Errorf("dispatcher is shutting down")

var timerPool = sync.Pool{}

func getTimer(d time.Duration) *time.Timer {
	if v := timerPool.Get(); v != nil {
		timer := v.(*time.Timer)
		timer.Reset(d)
		return timer
	}
	return time.NewTimer(d)
}

func putTimer(t *time.Timer) {
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}
	timerPool.Put(t)
}

// Execute sends a request to the dispatcher and waits for the response.
func (p *dispatcherPool[Q, R]) Execute(ctx context.Context, q Message[Q]) (Message[R], error) {
	d := p.getDispatcher()
	if q.GetID() == 0 {
		return d.zero, errors.New("ID must not be 0 (note that this must be unique within the current request set time window)")
	}

	call, err := d.newCall(q)
	if err != nil {
		return d.zero, err
	}

	if err := d.send(call); err != nil {
		return d.zero, err
	}

	timer := getTimer(30 * time.Second)
	defer putTimer(timer)

	select {
	case call = <-call.donec:
	case <-p.donec:
		return d.zero, p.Err()
	case <-ctx.Done():
		return d.zero, ctx.Err()
	case <-timer.C:
		return d.zero, errors.New("timeout")
	}

	if call.err != nil {
		return d.zero, call.err
	}

	resp, err := call.response, p.Err()
	if err == nil && resp.Header.Err != "" {
		err = errors.New(resp.Header.Err)
	}
	return resp, err
}

func (d *dispatcher[Q, R]) newCall(q Message[Q]) (*call[Q, R], error) {
	call := &call[Q, R]{
		donec:   make(chan *call[Q, R], 1),
		request: q,
	}

	if d.shutdown || d.closing {
		call.err = ErrShutdown
		call.done()
		return call, nil
	}

	d.mu.Lock()
	d.pending[q.GetID()] = call
	d.mu.Unlock()

	return call, nil
}

func (d *dispatcher[Q, R]) send(call *call[Q, R]) error {
	d.mu.RLock()
	if d.closing || d.shutdown {
		d.mu.RUnlock()
		return ErrShutdown
	}
	d.mu.RUnlock()

	d.encMu.Lock()
	defer d.encMu.Unlock()
	err := d.inOut.enc.Encode(call.request)
	if err != nil {
		return err
	}
	return nil
}

func (d *dispatcher[Q, R]) input() {
	var inputErr error

	for d.inOut.dec.More() {
		var r Message[R]
		if err := d.inOut.dec.Decode(&r); err != nil {
			inputErr = fmt.Errorf("decoding response: %w", err)
			break
		}

		d.mu.Lock()
		call, found := d.pending[r.GetID()]
		if !found {
			d.mu.Unlock()
			panic(fmt.Errorf("call with ID %d not found", r.GetID()))
		}
		delete(d.pending, r.GetID())
		d.mu.Unlock()
		call.response = r
		call.done()
	}

	// Terminate pending calls.
	d.shutdown = true
	if inputErr != nil {
		isEOF := inputErr == io.EOF || strings.Contains(inputErr.Error(), "already closed")
		if isEOF {
			if d.closing {
				inputErr = ErrShutdown
			} else {
				inputErr = io.ErrUnexpectedEOF
			}
		}
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	for _, call := range d.pending {
		call.err = inputErr
		call.done()
	}
}

type call[Q, R any] struct {
	request  Message[Q]
	response Message[R]
	err      error
	donec    chan *call[Q, R]
}

func (call *call[Q, R]) done() {
	select {
	case call.donec <- call:
	default:
	}
}

// Binary represents a WebAssembly binary.
type Binary struct {
	// The name of the binary.
	// For 	quickjs, this must match the instance import name, "javy_quickjs_provider_v2".
	// For the main module, we only use this for caching.
	Name string

	// THe wasm binary.
	Data []byte
}

type Options struct {
	Ctx context.Context

	Infof func(format string, v ...any)

	// E.g. quickjs wasm. May be omitted if not needed.
	Runtime Binary

	// The main module to instantiate.
	Main Binary

	CompilationCacheDir string
	PoolSize            int

	// Memory limit in MiB.
	Memory int
}

type CompileModuleContext struct {
	Opts    Options
	Runtime wazero.Runtime
}

type CompiledModule struct {
	// Runtime (e.g. QuickJS) may be nil if not needed (e.g. embedded in Module).
	Runtime wazero.CompiledModule

	// If Runtime is not nil, this should be the name of the instance.
	RuntimeName string

	// The main module to instantiate.
	// This will be insantiated multiple times in a pool,
	// so it does not need a name.
	Module wazero.CompiledModule
}

// Start creates a new dispatcher pool.
func Start[Q, R any](opts Options) (Dispatcher[Q, R], error) {
	if opts.Main.Data == nil {
		return nil, errors.New("Main.Data must be set")
	}
	if opts.Main.Name == "" {
		return nil, errors.New("Main.Name must be set")
	}

	if opts.Runtime.Data != nil && opts.Runtime.Name == "" {
		return nil, errors.New("Runtime.Name must be set")
	}

	if opts.PoolSize == 0 {
		opts.PoolSize = 1
	}

	return newDispatcher[Q, R](opts)
}

type dispatcherPool[Q, R any] struct {
	counter     atomic.Uint32
	dispatchers []*dispatcher[Q, R]
	close       func() error

	errc  chan error
	donec chan struct{}
}

func (p *dispatcherPool[Q, R]) SendIfErr(err error) {
	if err != nil {
		p.errc <- err
	}
}

func (p *dispatcherPool[Q, R]) Err() error {
	select {
	case err := <-p.errc:
		return err
	default:
		return nil
	}
}

func newDispatcher[Q, R any](opts Options) (*dispatcherPool[Q, R], error) {
	if opts.Ctx == nil {
		opts.Ctx = context.Background()
	}

	if opts.Infof == nil {
		opts.Infof = func(format string, v ...any) {
			// noop
		}
	}

	if opts.Memory <= 0 {
		// 32 MiB
		opts.Memory = 32
	}

	ctx := opts.Ctx

	// Page size is 64KB.
	numPages := opts.Memory * 1024 / 64
	runtimeConfig := wazero.NewRuntimeConfig().WithMemoryLimitPages(uint32(numPages))

	if opts.CompilationCacheDir != "" {
		compilationCache, err := wazero.NewCompilationCacheWithDir(opts.CompilationCacheDir)
		if err != nil {
			return nil, err
		}
		runtimeConfig = runtimeConfig.WithCompilationCache(compilationCache)
	}

	// Create a new WebAssembly Runtime.
	r := wazero.NewRuntimeWithConfig(opts.Ctx, runtimeConfig)

	// Instantiate WASI, which implements system I/O such as console output.
	if _, err := wasi_snapshot_preview1.Instantiate(ctx, r); err != nil {
		return nil, err
	}

	inOuts := make([]*inOut, opts.PoolSize)
	for i := 0; i < opts.PoolSize; i++ {
		var stdin, stdout hugio.ReadWriteCloser

		stdin = hugio.NewPipeReadWriteCloser()
		stdout = hugio.NewPipeReadWriteCloser()

		inOuts[i] = &inOut{
			stdin:  stdin,
			stdout: stdout,
			dec:    json.NewDecoder(stdout),
			enc:    json.NewEncoder(stdin),
		}
	}

	var (
		runtimeModule wazero.CompiledModule
		mainModule    wazero.CompiledModule
		err           error
	)

	if opts.Runtime.Data != nil {
		runtimeModule, err = r.CompileModule(ctx, opts.Runtime.Data)
		if err != nil {
			return nil, err
		}
	}

	mainModule, err = r.CompileModule(ctx, opts.Main.Data)
	if err != nil {
		return nil, err
	}

	toErr := func(what string, errBuff bytes.Buffer, err error) error {
		return fmt.Errorf("%s: %s: %w", what, errBuff.String(), err)
	}

	run := func() error {
		g, ctx := errgroup.WithContext(ctx)
		for _, c := range inOuts {
			c := c
			g.Go(func() error {
				var errBuff bytes.Buffer
				ctx := context.WithoutCancel(ctx)
				configBase := wazero.NewModuleConfig().WithStderr(&errBuff).WithStdout(c.stdout).WithStdin(c.stdin).WithStartFunctions()
				if opts.Runtime.Data != nil {
					// This needs to be anonymous, it will be resolved in the import resolver below.
					runtimeInstance, err := r.InstantiateModule(ctx, runtimeModule, configBase.WithName(""))
					if err != nil {
						return toErr("quickjs", errBuff, err)
					}
					ctx = experimental.WithImportResolver(ctx,
						func(name string) api.Module {
							if name == opts.Runtime.Name {
								return runtimeInstance
							}
							return nil
						},
					)
				}

				mainInstance, err := r.InstantiateModule(ctx, mainModule, configBase.WithName(""))
				if err != nil {
					return toErr(opts.Main.Name, errBuff, err)
				}
				if _, err := mainInstance.ExportedFunction("_start").Call(ctx); err != nil {
					return toErr(opts.Main.Name, errBuff, err)
				}

				// The console.log in the Javy/quickjs WebAssembly module will write to stderr.
				// In non-error situations, write that to the provided infof logger.
				if errBuff.Len() > 0 {
					opts.Infof("%s", errBuff.String())
				}

				return nil
			})
		}
		return g.Wait()
	}

	dp := &dispatcherPool[Q, R]{
		dispatchers: make([]*dispatcher[Q, R], len(inOuts)),

		errc:  make(chan error, 10),
		donec: make(chan struct{}),
	}

	go func() {
		// This will block until stdin is closed or it encounters an error.
		err := run()
		dp.SendIfErr(err)
		close(dp.donec)
	}()

	for i := 0; i < len(inOuts); i++ {
		d := &dispatcher[Q, R]{
			pending: make(map[uint32]*call[Q, R]),
			inOut:   inOuts[i],
		}
		go d.input()
		dp.dispatchers[i] = d
	}

	dp.close = func() error {
		for _, d := range dp.dispatchers {
			d.closing = true
			if err := d.inOut.stdin.Close(); err != nil {
				return err
			}
			if err := d.inOut.stdout.Close(); err != nil {
				return err
			}
		}

		// We need to wait for the WebAssembly instances to finish executing before we can close the runtime.
		<-dp.donec

		if err := r.Close(ctx); err != nil {
			return err
		}

		// Return potential late compilation errors.
		return dp.Err()
	}

	return dp, dp.Err()
}

type lazyDispatcher[Q, R any] struct {
	opts Options

	dispatcher Dispatcher[Q, R]
	startOnce  sync.Once
	started    bool
	startErr   error
}

func (d *lazyDispatcher[Q, R]) start() (Dispatcher[Q, R], error) {
	d.startOnce.Do(func() {
		start := time.Now()
		d.dispatcher, d.startErr = Start[Q, R](d.opts)
		d.started = true
		d.opts.Infof("started dispatcher in %s", time.Since(start))
	})
	return d.dispatcher, d.startErr
}

// Dispatchers holds all the dispatchers for the warpc package.
type Dispatchers struct {
	katex *lazyDispatcher[KatexInput, KatexOutput]
}

func (d *Dispatchers) Katex() (Dispatcher[KatexInput, KatexOutput], error) {
	return d.katex.start()
}

func (d *Dispatchers) Close() error {
	var errs []error
	if d.katex.started {
		if err := d.katex.dispatcher.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("%v", errs)
}

// AllDispatchers creates all the dispatchers for the warpc package.
// Note that the individual dispatchers are started lazily.
// Remember to call Close on the returned Dispatchers when done.
func AllDispatchers(katexOpts Options) *Dispatchers {
	if katexOpts.Runtime.Data == nil {
		katexOpts.Runtime = Binary{Name: "javy_quickjs_provider_v2", Data: quickjsWasm}
	}
	if katexOpts.Main.Data == nil {
		katexOpts.Main = Binary{Name: "renderkatex", Data: katexWasm}
	}

	if katexOpts.Infof == nil {
		katexOpts.Infof = func(format string, v ...any) {
			// noop
		}
	}

	return &Dispatchers{
		katex: &lazyDispatcher[KatexInput, KatexOutput]{opts: katexOpts},
	}
}
