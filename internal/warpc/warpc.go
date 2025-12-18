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
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bep/textandbinarywriter"

	"github.com/gohugoio/hugo/common/hstrings"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/maps"
	"golang.org/x/sync/errgroup"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/experimental"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

const (
	MessageKindJSON string = "json"
	MessageKindBlob string = "blob"
)

const currentVersion = 1

//go:embed wasm/quickjs.wasm
var quickjsWasm []byte

//go:embed wasm/webp.wasm
var webpWasm []byte

// Header is in both the request and response.
type Header struct {
	// Major version of the protocol.
	Version uint16 `json:"version"`

	// Unique ID for the request.
	// Note that this only needs to be unique within the current request set time window.
	ID uint32 `json:"id"`

	// Command is the command to execute.
	Command string `json:"command"`

	// RequestKinds is a list of kinds in this RPC request,
	// e.g. {"json", "blob"}, or {"json"}.
	RequestKinds []string `json:"requestKinds"`
	// ResponseKinds is a list of kinds expected in the response,
	// e.g. {"json", "blob"}, or {"json"}.
	ResponseKinds []string `json:"responseKinds"`

	// Set in the response if there was an error.
	Err string `json:"err"`

	// Warnings is a list of warnings that may be returned in the response.
	Warnings []string `json:"warnings,omitempty"`
}

func (m *Header) init() error {
	if m.ID == 0 {
		return errors.New("ID must not be 0 (note that this must be unique within the current request set time window)")
	}
	if m.Version == 0 {
		m.Version = currentVersion
	}
	if len(m.RequestKinds) == 0 {
		m.RequestKinds = []string{string(MessageKindJSON)}
	}
	if len(m.ResponseKinds) == 0 {
		m.ResponseKinds = []string{string(MessageKindJSON)}
	}
	if m.Version != currentVersion {
		return fmt.Errorf("unsupported version: %d", m.Version)
	}
	for range 2 {
		if len(m.RequestKinds) > 2 {
			return fmt.Errorf("invalid number of request kinds: %d", len(m.RequestKinds))
		}
		if len(m.ResponseKinds) > 2 {
			return fmt.Errorf("invalid number of response kinds: %d", len(m.ResponseKinds))
		}
		m.RequestKinds = hstrings.UniqueStringsReuse(m.RequestKinds)
		m.ResponseKinds = hstrings.UniqueStringsReuse(m.ResponseKinds)

	}

	return nil
}

type Message[T any] struct {
	Header Header `json:"header"`
	Data   T      `json:"data"`
}

func (m Message[T]) GetID() uint32 {
	return m.Header.ID
}

func (m *Message[T]) init() error {
	return m.Header.init()
}

type SourceProvider interface {
	GetSource() hugio.SizeReader
}

type DestinationProvider interface {
	GetDestination() io.Writer
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
	zeroR Message[R]

	id atomic.Uint32

	mu       sync.Mutex
	encodeMu sync.Mutex

	pending map[uint32]*call[Q, R]

	inOut   *inOut
	inGroup *errgroup.Group

	shutdown bool
	closing  bool
}

type inOut struct {
	sync.Mutex
	stdin        hugio.ReadWriteCloser
	stdout       io.WriteCloser
	stdoutBinary hugio.ReadWriteCloser
	dec          *json.Decoder
	enc          *json.Encoder
}

func (p *inOut) Close() error {
	if err := p.stdin.Close(); err != nil {
		return err
	}

	// This will also close the underlying writers.
	if err := p.stdout.Close(); err != nil {
		return err
	}
	return nil
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

	call, err := d.newCall(q)
	if err != nil {
		return d.zeroR, err
	}

	if err := d.send(call); err != nil {
		return d.zeroR, err
	}

	timer := getTimer(30 * time.Second)
	defer putTimer(timer)

	select {
	case call = <-call.donec:
	case <-p.donec:
		return d.zeroR, p.Err()
	case <-ctx.Done():
		return d.zeroR, ctx.Err()
	case <-timer.C:
		return d.zeroR, errors.New("timeout")
	}

	if call.err != nil {
		return d.zeroR, call.err
	}

	resp, err := call.response, p.Err()

	if err == nil && resp.Header.Err != "" {
		err = errors.New(resp.Header.Err)
	}

	return resp, err
}

func (d *dispatcher[Q, R]) newCall(q Message[Q]) (*call[Q, R], error) {
	if q.Header.ID == 0 {
		q.Header.ID = d.id.Add(1)
	}
	if err := q.init(); err != nil {
		return nil, err
	}
	responseKinds := maps.NewMap[string, bool]()
	for _, rk := range q.Header.ResponseKinds {
		responseKinds.Set(rk, true)
	}
	call := &call[Q, R]{
		donec:         make(chan *call[Q, R], 1),
		request:       q,
		responseKinds: responseKinds,
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
	d.mu.Lock()
	if d.closing || d.shutdown {
		d.mu.Unlock()
		return ErrShutdown
	}
	d.mu.Unlock()

	d.encodeMu.Lock()
	defer d.encodeMu.Unlock()
	err := d.inOut.enc.Encode(call.request)
	if err != nil {
		return err
	}
	if sp, ok := any(call.request.Data).(SourceProvider); ok {
		source := sp.GetSource()
		if source.Size() == 0 {
			return errors.New("source size must be greater than 0")
		}

		if err := textandbinarywriter.WriteBlobHeader(d.inOut.stdin, call.request.GetID(), uint32(source.Size())); err != nil {
			return err
		}
		_, err := io.Copy(d.inOut.stdin, source)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *dispatcher[Q, R]) inputBlobs() error {
	var inputErr error
	for {
		id, length, err := textandbinarywriter.ReadBlobHeader(d.inOut.stdoutBinary)
		if err != nil {
			inputErr = err
			break
		}

		lr := &io.LimitedReader{
			R: d.inOut.stdoutBinary,
			N: int64(length),
		}

		call := d.pendingCall(id)

		if err := call.handleBlob(lr); err != nil {
			inputErr = err
			break
		}
		if lr.N != 0 {
			inputErr = fmt.Errorf("blob %d: expected to read %d more bytes", id, lr.N)
			break
		}
		if err := call.responseKinds.WithWriteLock(
			func(m map[string]bool) error {
				if _, ok := m[MessageKindBlob]; !ok {
					return fmt.Errorf("unexpected blob response for %q call ID %d", call.request.Header.Command, id)
				}
				delete(m, MessageKindBlob)
				if len(m) == 0 {
					// Message exchange is complete.
					d.mu.Lock()
					delete(d.pending, id)
					d.mu.Unlock()
					call.done()
				}
				return nil
			}); err != nil {
			inputErr = err
			break
		}
	}

	return inputErr
}

func (d *dispatcher[Q, R]) inputJSON() error {
	var inputErr error

	for d.inOut.dec.More() {
		var r Message[R]
		if err := d.inOut.dec.Decode(&r); err != nil {
			inputErr = err
			break
		}

		call := d.pendingCall(r.GetID())

		if err := call.responseKinds.WithWriteLock(
			func(m map[string]bool) error {
				call.response = r
				if _, ok := m[MessageKindJSON]; !ok {
					return fmt.Errorf("unexpected JSON response for call ID %d", r.GetID())
				}
				delete(m, MessageKindJSON)
				if len(m) == 0 || r.Header.Err != "" {
					// Message exchange is complete.
					d.mu.Lock()
					delete(d.pending, r.GetID())
					d.mu.Unlock()
					call.done()
				}
				return nil
			}); err != nil {
			inputErr = err
			break
		}

	}

	// Terminate pending calls.
	d.shutdown = true
	if inputErr != nil {
		isEOF := inputErr == io.EOF || inputErr == io.ErrClosedPipe || strings.Contains(inputErr.Error(), "already closed")
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

	return inputErr
}

func (d *dispatcher[Q, R]) pendingCall(id uint32) *call[Q, R] {
	d.mu.Lock()
	defer d.mu.Unlock()
	c, ok := d.pending[id]
	if !ok {
		panic(fmt.Errorf("call with ID %d not found", id))
	}
	return c
}

type call[Q, R any] struct {
	request       Message[Q]
	response      Message[R]
	responseKinds *maps.Map[string, bool]
	err           error
	donec         chan *call[Q, R]
}

func (c *call[Q, R]) handleBlob(r io.Reader) error {
	dest := any(c.request.Data).(DestinationProvider).GetDestination()
	if dest == nil {
		panic("blob destination is not set")
	}
	_, err := io.Copy(dest, r)
	return err
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

type ProtocolType int

const (
	ProtocolJSON          ProtocolType = iota + 1
	ProtocolJSONAndBinary              = iota
)

type Options struct {
	Ctx context.Context

	Infof func(format string, v ...any)

	Warnf func(format string, v ...any)

	// E.g. quickjs wasm. May be omitted if not needed.
	Runtime Binary

	// The main module to instantiate.
	Main Binary

	CompilationCacheDir string
	PoolSize            int

	// Memory limit in MiB.
	Memory int
}

func (o *Options) init() error {
	if o.Infof == nil {
		o.Infof = func(format string, v ...any) {
		}
	}
	if o.Warnf == nil {
		o.Warnf = func(format string, v ...any) {
		}
	}
	return nil
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
	opts        Options

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
	if opts.Warnf == nil {
		opts.Warnf = func(format string, v ...any) {
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

	dispatchers := make([]*dispatcher[Q, R], opts.PoolSize)

	inOuts := make([]*inOut, opts.PoolSize)
	for i := range opts.PoolSize {
		var stdinPipe, stdoutBinary hugio.ReadWriteCloser
		var stdout io.WriteCloser
		var jsonr io.Reader

		stdinPipe = hugio.NewPipeReadWriteCloser()
		stdoutPipe := hugio.NewPipeReadWriteCloser()
		stdout = stdoutPipe

		var zero Q
		if _, ok := any(zero).(DestinationProvider); ok {
			stdoutBinary = hugio.NewPipeReadWriteCloser()
			jsonr = stdoutPipe
			stdout = textandbinarywriter.NewWriter(stdout, stdoutBinary)
		} else {
			jsonr = stdoutPipe
		}

		inOuts[i] = &inOut{
			stdin:        stdinPipe,
			stdout:       stdout,
			stdoutBinary: stdoutBinary,
			dec:          json.NewDecoder(jsonr),
			enc:          json.NewEncoder(stdinPipe),
		}
	}

	var (
		runtimeModule wazero.CompiledModule
		mainModule    wazero.CompiledModule
		err           error
	)

	if opts.Runtime.Data != nil {
		runtimeModule, err = r.CompileModule(experimental.WithCompilationWorkers(ctx, runtime.GOMAXPROCS(0)/4), opts.Runtime.Data)
		if err != nil {
			return nil, err
		}
	}

	mainModule, err = r.CompileModule(experimental.WithCompilationWorkers(ctx, runtime.GOMAXPROCS(0)/4), opts.Main.Data)
	if err != nil {
		return nil, err
	}

	toErr := func(what string, errBuff bytes.Buffer, err error) error {
		return fmt.Errorf("%s: %s: %w", what, errBuff.String(), err)
	}

	run := func() error {
		g, ctx := errgroup.WithContext(ctx)
		for _, c := range inOuts {
			g.Go(func() error {
				var errBuff bytes.Buffer
				stderr := io.MultiWriter(&errBuff, os.Stderr)
				ctx := context.WithoutCancel(ctx)
				// WithStartFunctions allows us to call the _start function manually.
				configBase := wazero.NewModuleConfig().WithStderr(stderr).WithStdout(c.stdout).WithStdin(c.stdin).WithStartFunctions()
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
		dispatchers: dispatchers,
		opts:        opts,

		errc:  make(chan error, 10),
		donec: make(chan struct{}),
	}

	go func() {
		// This will block until stdin is closed or it encounters an error.
		err := run()
		dp.SendIfErr(err)
		close(dp.donec)
	}()

	for i := range inOuts {
		ing, _ := errgroup.WithContext(ctx)
		d := &dispatcher[Q, R]{
			pending: make(map[uint32]*call[Q, R]),
			inOut:   inOuts[i],
			inGroup: ing,
		}
		d.inGroup.Go(func() error {
			return d.inputJSON()
		})

		if d.inOut.stdoutBinary != nil {
			d.inGroup.Go(func() error {
				return d.inputBlobs()
			})
		}
		dp.dispatchers[i] = d
	}

	dp.close = func() error {
		for _, d := range dp.dispatchers {
			d.closing = true
			if err := d.inOut.Close(); err != nil {
				return err
			}
		}

		for _, d := range dp.dispatchers {
			if err := d.inGroup.Wait(); err != nil {
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
	webp  *lazyDispatcher[WebpInput, WebpOutput]
}

func (d *Dispatchers) Katex() (Dispatcher[KatexInput, KatexOutput], error) {
	return d.katex.start()
}

func (d *Dispatchers) Webp() (Dispatcher[WebpInput, WebpOutput], error) {
	return d.webp.start()
}

func (d *Dispatchers) NewWepCodec(quality int, hint string) (*WebpCodec, error) {
	return &WebpCodec{
		d:       d.Webp,
		quality: quality,
		hint:    hint,
	}, nil
}

func (d *Dispatchers) Close() error {
	var errs []error
	if d.katex.started {
		if err := d.katex.dispatcher.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if d.webp.started {
		if err := d.webp.dispatcher.Close(); err != nil {
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
func AllDispatchers(katexOpts, webpOpts Options) *Dispatchers {
	if err := katexOpts.init(); err != nil {
		panic(err)
	}
	if err := webpOpts.init(); err != nil {
		panic(err)
	}
	if katexOpts.Runtime.Data == nil {
		katexOpts.Runtime = Binary{Name: "javy_quickjs_provider_v2", Data: quickjsWasm}
	}
	if katexOpts.Main.Data == nil {
		katexOpts.Main = Binary{Name: "renderkatex", Data: katexWasm}
	}

	webpOpts.Main = Binary{Name: "webp", Data: webpWasm}

	dispatchers := &Dispatchers{
		katex: &lazyDispatcher[KatexInput, KatexOutput]{opts: katexOpts},
		webp:  &lazyDispatcher[WebpInput, WebpOutput]{opts: webpOpts},
	}

	return dispatchers
}
