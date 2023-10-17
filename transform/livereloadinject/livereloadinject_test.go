// Copyright 2018 The Hugo Authors. All rights reserved.
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

package livereloadinject

import (
	"bytes"
	"net/url"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/transform"
)

func TestLiveReloadInject(t *testing.T) {
	c := qt.New(t)

	lrurl, err := url.Parse("http://localhost:1234/subpath")
	if err != nil {
		t.Errorf("Parsing test URL failed")
		return
	}
	expectBase := `<script src="/subpath/livereload.js?mindelay=10&amp;v=2&amp;port=1234&amp;path=subpath/livereload" data-no-instant defer></script>`
	apply := func(s string) string {
		out := new(bytes.Buffer)
		in := strings.NewReader(s)

		tr := transform.New(New(*lrurl))
		tr.Apply(out, in)

		return out.String()
	}

	c.Run("Head lower", func(c *qt.C) {
		c.Assert(apply("<html><head>after"), qt.Equals, "<html><head>"+expectBase+"after")
	})

	c.Run("Head upper", func(c *qt.C) {
		c.Assert(apply("<HTML><HEAD>after"), qt.Equals, "<HTML><HEAD>"+expectBase+"after")
	})

	c.Run("Head mixed", func(c *qt.C) {
		c.Assert(apply("<Html><Head>after"), qt.Equals, "<Html><Head>"+expectBase+"after")
	})

	c.Run("Html if head missing", func(c *qt.C) {
		c.Assert(apply("<!doctype><html>after"), qt.Equals, "<!doctype><html>"+expectBase+"after")
	})

	c.Run("Html with attr", func(c *qt.C) {
		c.Assert(apply(`<html lang="en">after`), qt.Equals, `<html lang="en">`+expectBase+"after")
	})

	c.Run("Html with newline", func(c *qt.C) {
		c.Assert(apply("<html\n>after"), qt.Equals, "<html\n>"+expectBase+"after")
	})

	c.Run("Do not mistake header for head", func(c *qt.C) {
		c.Assert(apply("<!doctype><html><header>"), qt.Equals, "<!doctype><html>"+expectBase+"<header>")
	})

	c.Run("Do not mistake custom elements for head", func(c *qt.C) {
		c.Assert(apply("<!doctype><html><head-custom>"), qt.Equals, "<!doctype><html>"+expectBase+"<head-custom>")
	})

	c.Run("Doctype lower", func(c *qt.C) {
		c.Assert(apply("<!doctype html>after"), qt.Equals, "<!doctype html>"+expectBase+"after")
	})

	c.Run("Doctype mixed", func(c *qt.C) {
		c.Assert(apply("<!Doctype Html>after"), qt.Equals, "<!Doctype Html>"+expectBase+"after")
	})

	c.Run("Fallback to before first element", func(c *qt.C) {
		c.Assert(apply("<h1>No match</h1>"), qt.Equals, expectBase+"<h1>No match</h1>")
	})

	c.Run("Do not fallback before BOM", func(c *qt.C) {
		c.Assert(apply("\uFEFF<h1>No match</h1>"), qt.Equals, "\uFEFF"+expectBase+"<h1>No match</h1>")
	})

	c.Run("Search from the start of the input", func(c *qt.C) {
		c.Assert(apply("<head>after<!--<head>-->"), qt.Equals, "<head>"+expectBase+"after<!--<head>-->")
	})

	c.Run("Do not search in title", func(c *qt.C) {
		c.Assert(apply("<html><title>The <head> element</title>"), qt.Equals, "<html>"+expectBase+"<title>The <head> element</title>")
	})

	c.Run("Do not search in comment", func(c *qt.C) {
		c.Assert(apply("<html><!--<head>-->"), qt.Equals, "<html>"+expectBase+"<html><!--<head>-->")
	})

	c.Run("Do not search in subelements", func(c *qt.C) {
		c.Assert(apply("<html><template><head></template>"), qt.Equals, "<html>"+expectBase+"<template><head></template>")
	})
}
