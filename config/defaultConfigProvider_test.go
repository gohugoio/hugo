// Copyright 2021 The Hugo Authors. All rights reserved.
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

package config

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/common/para"

	"github.com/gohugoio/hugo/common/maps"

	qt "github.com/frankban/quicktest"
)

func TestDefaultConfigProvider(t *testing.T) {
	c := qt.New(t)

	c.Run("Set and get", func(c *qt.C) {
		cfg := New()
		var k string
		var v any

		k, v = "foo", "bar"
		cfg.Set(k, v)
		c.Assert(cfg.Get(k), qt.Equals, v)
		c.Assert(cfg.Get(strings.ToUpper(k)), qt.Equals, v)
		c.Assert(cfg.GetString(k), qt.Equals, v)

		k, v = "foo", 42
		cfg.Set(k, v)
		c.Assert(cfg.Get(k), qt.Equals, v)
		c.Assert(cfg.GetInt(k), qt.Equals, v)

		c.Assert(cfg.Get(""), qt.DeepEquals, maps.Params{
			"foo": 42,
		})
	})

	c.Run("Set and get map", func(c *qt.C) {
		cfg := New()

		cfg.Set("foo", map[string]any{
			"bar": "baz",
		})

		c.Assert(cfg.Get("foo"), qt.DeepEquals, maps.Params{
			"bar": "baz",
		})

		c.Assert(cfg.GetStringMap("foo"), qt.DeepEquals, map[string]any{"bar": string("baz")})
		c.Assert(cfg.GetStringMapString("foo"), qt.DeepEquals, map[string]string{"bar": string("baz")})
	})

	c.Run("Set and get nested", func(c *qt.C) {
		cfg := New()

		cfg.Set("a", map[string]any{
			"B": "bv",
		})
		cfg.Set("a.c", "cv")

		c.Assert(cfg.Get("a"), qt.DeepEquals, maps.Params{
			"b": "bv",
			"c": "cv",
		})
		c.Assert(cfg.Get("a.c"), qt.Equals, "cv")

		cfg.Set("b.a", "av")
		c.Assert(cfg.Get("b"), qt.DeepEquals, maps.Params{
			"a": "av",
		})

		cfg.Set("b", map[string]any{
			"b": "bv",
		})

		c.Assert(cfg.Get("b"), qt.DeepEquals, maps.Params{
			"a": "av",
			"b": "bv",
		})

		cfg = New()

		cfg.Set("a", "av")

		cfg.Set("", map[string]any{
			"a": "av2",
			"b": "bv2",
		})

		c.Assert(cfg.Get(""), qt.DeepEquals, maps.Params{
			"a": "av2",
			"b": "bv2",
		})

		cfg = New()

		cfg.Set("a", "av")

		cfg.Set("", map[string]any{
			"b": "bv2",
		})

		c.Assert(cfg.Get(""), qt.DeepEquals, maps.Params{
			"a": "av",
			"b": "bv2",
		})

		cfg = New()

		cfg.Set("", map[string]any{
			"foo": map[string]any{
				"a": "av",
			},
		})

		cfg.Set("", map[string]any{
			"foo": map[string]any{
				"b": "bv2",
			},
		})

		c.Assert(cfg.Get("foo"), qt.DeepEquals, maps.Params{
			"a": "av",
			"b": "bv2",
		})
	})

	c.Run("Merge default strategy", func(c *qt.C) {
		cfg := New()

		cfg.Set("a", map[string]any{
			"B": "bv",
		})

		cfg.Merge("a", map[string]any{
			"B": "bv2",
			"c": "cv2",
		})

		c.Assert(cfg.Get("a"), qt.DeepEquals, maps.Params{
			"b": "bv",
			"c": "cv2",
		})

		cfg = New()

		cfg.Set("a", "av")

		cfg.Merge("", map[string]any{
			"a": "av2",
			"b": "bv2",
		})

		c.Assert(cfg.Get(""), qt.DeepEquals, maps.Params{
			"a": "av",
		})
	})

	c.Run("Merge shallow", func(c *qt.C) {
		cfg := New()

		cfg.Set("a", map[string]any{
			"_merge": "shallow",
			"B":      "bv",
			"c": map[string]any{
				"b": "bv",
			},
		})

		cfg.Merge("a", map[string]any{
			"c": map[string]any{
				"d": "dv2",
			},
			"e": "ev2",
		})

		c.Assert(cfg.Get("a"), qt.DeepEquals, maps.Params{
			"e":      "ev2",
			"_merge": maps.ParamsMergeStrategyShallow,
			"b":      "bv",
			"c": maps.Params{
				"b": "bv",
			},
		})
	})

	// Issue #8679
	c.Run("Merge typed maps", func(c *qt.C) {
		for _, left := range []any{
			map[string]string{
				"c": "cv1",
			},
			map[string]any{
				"c": "cv1",
			},
			map[any]any{
				"c": "cv1",
			},
		} {
			cfg := New()

			cfg.Set("", map[string]any{
				"b": left,
			})

			cfg.Merge("", maps.Params{
				"b": maps.Params{
					"c": "cv2",
					"d": "dv2",
				},
			})

			c.Assert(cfg.Get(""), qt.DeepEquals, maps.Params{
				"b": maps.Params{
					"c": "cv1",
					"d": "dv2",
				},
			})
		}

		for _, left := range []any{
			map[string]string{
				"b": "bv1",
			},
			map[string]any{
				"b": "bv1",
			},
			map[any]any{
				"b": "bv1",
			},
		} {
			for _, right := range []any{
				map[string]string{
					"b": "bv2",
					"c": "cv2",
				},
				map[string]any{
					"b": "bv2",
					"c": "cv2",
				},
				map[any]any{
					"b": "bv2",
					"c": "cv2",
				},
			} {
				cfg := New()

				cfg.Set("a", left)

				cfg.Merge("a", right)

				c.Assert(cfg.Get(""), qt.DeepEquals, maps.Params{
					"a": maps.Params{
						"b": "bv1",
						"c": "cv2",
					},
				})
			}
		}
	})

	// Issue #8701
	c.Run("Prevent _merge only maps", func(c *qt.C) {
		cfg := New()

		cfg.Set("", map[string]any{
			"B": "bv",
		})

		cfg.Merge("", map[string]any{
			"c": map[string]any{
				"_merge": "shallow",
				"d":      "dv2",
			},
		})

		c.Assert(cfg.Get(""), qt.DeepEquals, maps.Params{
			"b": "bv",
		})
	})

	c.Run("IsSet", func(c *qt.C) {
		cfg := New()

		cfg.Set("a", map[string]any{
			"B": "bv",
		})

		c.Assert(cfg.IsSet("A"), qt.IsTrue)
		c.Assert(cfg.IsSet("a.b"), qt.IsTrue)
		c.Assert(cfg.IsSet("z"), qt.IsFalse)
	})

	c.Run("Para", func(c *qt.C) {
		cfg := New()
		p := para.New(4)
		r, _ := p.Start(context.Background())

		setAndGet := func(k string, v int) error {
			vs := strconv.Itoa(v)
			cfg.Set(k, v)
			err := errors.New("get failed")
			if cfg.Get(k) != v {
				return err
			}
			if cfg.GetInt(k) != v {
				return err
			}
			if cfg.GetString(k) != vs {
				return err
			}
			if !cfg.IsSet(k) {
				return err
			}
			return nil
		}

		for i := 0; i < 20; i++ {
			i := i
			r.Run(func() error {
				const v = 42
				k := fmt.Sprintf("k%d", i)
				if err := setAndGet(k, v); err != nil {
					return err
				}

				m := maps.Params{
					"new": 42,
				}

				cfg.Merge("", m)

				return nil
			})
		}

		c.Assert(r.Wait(), qt.IsNil)
	})
}

func BenchmarkDefaultConfigProvider(b *testing.B) {
	type cfger interface {
		Get(key string) any
		Set(key string, value any)
		IsSet(key string) bool
	}

	newMap := func() map[string]any {
		return map[string]any{
			"a": map[string]any{
				"b": map[string]any{
					"c": 32,
					"d": 43,
				},
			},
			"b": 62,
		}
	}

	runMethods := func(b *testing.B, cfg cfger) {
		m := newMap()
		cfg.Set("mymap", m)
		cfg.Set("num", 32)
		if !(cfg.IsSet("mymap") && cfg.IsSet("mymap.a") && cfg.IsSet("mymap.a.b") && cfg.IsSet("mymap.a.b.c")) {
			b.Fatal("IsSet failed")
		}

		if cfg.Get("num") != 32 {
			b.Fatal("Get failed")
		}

		if cfg.Get("mymap.a.b.c") != 32 {
			b.Fatal("Get failed")
		}
	}

	b.Run("Custom", func(b *testing.B) {
		cfg := New()
		for i := 0; i < b.N; i++ {
			runMethods(b, cfg)
		}
	})
}
