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
	"context"
	_ "embed"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	qt "github.com/frankban/quicktest"
)

//go:embed wasm/greet.wasm
var greetWasm []byte

type person struct {
	Name string `json:"name"`
}

func TestKatex(t *testing.T) {
	c := qt.New(t)

	opts := Options{
		PoolSize: 8,
		Runtime:  quickjsBinary,
		Main:     katexBinary,
	}

	d, err := Start[KatexInput, KatexOutput](opts)
	c.Assert(err, qt.IsNil)

	defer d.Close()

	runExpression := func(c *qt.C, id uint32, expression string) (Message[KatexOutput], error) {
		c.Helper()

		ctx := context.Background()

		input := KatexInput{
			Expression: expression,
			Options: KatexOptions{
				Output:       "html",
				DisplayMode:  true,
				ThrowOnError: true,
			},
		}

		message := Message[KatexInput]{
			Header: Header{
				Version: currentVersion,
				ID:      uint32(id),
			},
			Data: input,
		}

		return d.Execute(ctx, message)
	}

	c.Run("Simple", func(c *qt.C) {
		id := uint32(32)
		result, err := runExpression(c, id, "c = \\pm\\sqrt{a^2 + b^2}")
		c.Assert(err, qt.IsNil)
		c.Assert(result.GetID(), qt.Equals, id)
	})

	c.Run("Invalid expression", func(c *qt.C) {
		id := uint32(32)
		result, err := runExpression(c, id, "c & \\foo\\")
		c.Assert(err, qt.IsNotNil)
		c.Assert(result.GetID(), qt.Equals, id)
	})
}

func TestGreet(t *testing.T) {
	c := qt.New(t)
	opts := Options{
		PoolSize: 1,
		Runtime:  quickjsBinary,
		Main:     greetBinary,
		Infof:    t.Logf,
	}

	for i := 0; i < 2; i++ {
		func() {
			d, err := Start[person, greeting](opts)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				c.Assert(d.Close(), qt.IsNil)
			}()

			ctx := context.Background()

			inputMessage := Message[person]{
				Header: Header{
					Version: currentVersion,
				},
				Data: person{
					Name: "Person",
				},
			}

			for j := 0; j < 20; j++ {
				inputMessage.Header.ID = uint32(j + 1)
				g, err := d.Execute(ctx, inputMessage)
				if err != nil {
					t.Fatal(err)
				}
				if g.Data.Greeting != "Hello Person!" {
					t.Fatalf("got: %v", g)
				}
				if g.GetID() != inputMessage.GetID() {
					t.Fatalf("%d vs %d", g.GetID(), inputMessage.GetID())
				}
			}
		}()
	}
}

func TestGreetParallel(t *testing.T) {
	c := qt.New(t)

	opts := Options{
		Runtime:  quickjsBinary,
		Main:     greetBinary,
		PoolSize: 4,
	}
	d, err := Start[person, greeting](opts)
	c.Assert(err, qt.IsNil)
	defer func() {
		c.Assert(d.Close(), qt.IsNil)
	}()

	var wg sync.WaitGroup

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			ctx := context.Background()

			for j := 0; j < 5; j++ {
				base := i * 100
				id := uint32(base + j)

				inputPerson := person{
					Name: fmt.Sprintf("Person %d", id),
				}
				inputMessage := Message[person]{
					Header: Header{
						Version: currentVersion,
						ID:      id,
					},
					Data: inputPerson,
				}
				g, err := d.Execute(ctx, inputMessage)
				if err != nil {
					t.Error(err)
					return
				}

				c.Assert(g.Data.Greeting, qt.Equals, fmt.Sprintf("Hello Person %d!", id))
				c.Assert(g.GetID(), qt.Equals, inputMessage.GetID())

			}
		}(i)

	}

	wg.Wait()
}

func TestKatexParallel(t *testing.T) {
	c := qt.New(t)

	opts := Options{
		Runtime:  quickjsBinary,
		Main:     katexBinary,
		PoolSize: 6,
	}
	d, err := Start[KatexInput, KatexOutput](opts)
	c.Assert(err, qt.IsNil)
	defer func() {
		c.Assert(d.Close(), qt.IsNil)
	}()

	var wg sync.WaitGroup

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			ctx := context.Background()

			for j := 0; j < 1; j++ {
				base := i * 100
				id := uint32(base + j)

				input := katexInputTemplate
				inputMessage := Message[KatexInput]{
					Header: Header{
						Version: currentVersion,
						ID:      id,
					},
					Data: input,
				}

				result, err := d.Execute(ctx, inputMessage)
				if err != nil {
					t.Error(err)
					return
				}

				if result.GetID() != inputMessage.GetID() {
					t.Errorf("%d vs %d", result.GetID(), inputMessage.GetID())
					return
				}
			}
		}(i)

	}

	wg.Wait()
}

func BenchmarkExecuteKatex(b *testing.B) {
	opts := Options{
		Runtime: quickjsBinary,
		Main:    katexBinary,
	}
	d, err := Start[KatexInput, KatexOutput](opts)
	if err != nil {
		b.Fatal(err)
	}
	defer d.Close()

	ctx := context.Background()

	input := katexInputTemplate

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		message := Message[KatexInput]{
			Header: Header{
				Version: currentVersion,
				ID:      uint32(i + 1),
			},
			Data: input,
		}

		result, err := d.Execute(ctx, message)
		if err != nil {
			b.Fatal(err)
		}

		if result.GetID() != message.GetID() {
			b.Fatalf("%d vs %d", result.GetID(), message.GetID())
		}

	}
}

func BenchmarkKatexStartStop(b *testing.B) {
	optsTemplate := Options{
		Runtime:             quickjsBinary,
		Main:                katexBinary,
		CompilationCacheDir: b.TempDir(),
	}

	runBench := func(b *testing.B, opts Options) {
		for i := 0; i < b.N; i++ {
			d, err := Start[KatexInput, KatexOutput](opts)
			if err != nil {
				b.Fatal(err)
			}
			if err := d.Close(); err != nil {
				b.Fatal(err)
			}
		}
	}

	for _, poolSize := range []int{1, 8, 16} {

		name := fmt.Sprintf("PoolSize%d", poolSize)

		b.Run(name, func(b *testing.B) {
			opts := optsTemplate
			opts.PoolSize = poolSize
			runBench(b, opts)
		})

	}
}

var katexInputTemplate = KatexInput{
	Expression: "c = \\pm\\sqrt{a^2 + b^2}",
	Options:    KatexOptions{Output: "html", DisplayMode: true},
}

func BenchmarkExecuteKatexPara(b *testing.B) {
	optsTemplate := Options{
		Runtime: quickjsBinary,
		Main:    katexBinary,
	}

	runBench := func(b *testing.B, opts Options) {
		d, err := Start[KatexInput, KatexOutput](opts)
		if err != nil {
			b.Fatal(err)
		}
		defer d.Close()

		ctx := context.Background()

		b.ResetTimer()

		var id atomic.Uint32
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				message := Message[KatexInput]{
					Header: Header{
						Version: currentVersion,
						ID:      id.Add(1),
					},
					Data: katexInputTemplate,
				}

				result, err := d.Execute(ctx, message)
				if err != nil {
					b.Fatal(err)
				}
				if result.GetID() != message.GetID() {
					b.Fatalf("%d vs %d", result.GetID(), message.GetID())
				}
			}
		})
	}

	for _, poolSize := range []int{1, 8, 16} {
		name := fmt.Sprintf("PoolSize%d", poolSize)

		b.Run(name, func(b *testing.B) {
			opts := optsTemplate
			opts.PoolSize = poolSize
			runBench(b, opts)
		})
	}
}

func BenchmarkExecuteGreet(b *testing.B) {
	opts := Options{
		Runtime: quickjsBinary,
		Main:    greetBinary,
	}
	d, err := Start[person, greeting](opts)
	if err != nil {
		b.Fatal(err)
	}
	defer d.Close()

	ctx := context.Background()

	input := person{
		Name: "Person",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		message := Message[person]{
			Header: Header{
				Version: currentVersion,
				ID:      uint32(i + 1),
			},
			Data: input,
		}
		result, err := d.Execute(ctx, message)
		if err != nil {
			b.Fatal(err)
		}

		if result.GetID() != message.GetID() {
			b.Fatalf("%d vs %d", result.GetID(), message.GetID())
		}

	}
}

func BenchmarkExecuteGreetPara(b *testing.B) {
	opts := Options{
		Runtime:  quickjsBinary,
		Main:     greetBinary,
		PoolSize: 8,
	}

	d, err := Start[person, greeting](opts)
	if err != nil {
		b.Fatal(err)
	}
	defer d.Close()

	ctx := context.Background()

	inputTemplate := person{
		Name: "Person",
	}

	b.ResetTimer()

	var id atomic.Uint32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			message := Message[person]{
				Header: Header{
					Version: currentVersion,
					ID:      id.Add(1),
				},
				Data: inputTemplate,
			}

			result, err := d.Execute(ctx, message)
			if err != nil {
				b.Fatal(err)
			}
			if result.GetID() != message.GetID() {
				b.Fatalf("%d vs %d", result.GetID(), message.GetID())
			}
		}
	})
}

type greeting struct {
	Greeting string `json:"greeting"`
}

var (
	greetBinary = Binary{
		Name: "greet",
		Data: greetWasm,
	}

	katexBinary = Binary{
		Name: "renderkatex",
		Data: katexWasm,
	}

	quickjsBinary = Binary{
		Name: "javy_quickjs_provider_v2",
		Data: quickjsWasm,
	}
)
