// Keywords and other highlight in this script have been modified for better selection of Hugo-specific keywords
/*! highlight.js v9.9.0 | BSD3 License | git.io/hljslicense */
! function(e) {
  var n = "object" == typeof window && window || "object" == typeof self && self;
  "undefined" != typeof exports ? e(exports) : n && (n.hljs = e({}), "function" == typeof define && define.amd && define([], function() {
    return n.hljs
  }))
}(function(e) {
  function n(e) {
    return e.replace(/[&<>]/gm, function(e) {
      return I[e]
    })
  }

  function t(e) {
    return e.nodeName.toLowerCase()
  }

  function r(e, n) {
    var t = e && e.exec(n);
    return t && 0 === t.index
  }

  function i(e) {
    return k.test(e)
  }

  function a(e) {
    var n, t, r, a, o = e.className + " ";
    if (o += e.parentNode ? e.parentNode.className : "", t = B.exec(o)) return R(t[1]) ? t[1] : "no-highlight";
    for (o = o.split(/\s+/), n = 0, r = o.length; r > n; n++)
      if (a = o[n], i(a) || R(a)) return a
  }

  function o(e, n) {
    var t, r = {};
    for (t in e) r[t] = e[t];
    if (n)
      for (t in n) r[t] = n[t];
    return r
  }

  function u(e) {
    var n = [];
    return function r(e, i) {
      for (var a = e.firstChild; a; a = a.nextSibling) 3 === a.nodeType ? i += a.nodeValue.length : 1 === a.nodeType && (n.push({ event: "start", offset: i, node: a }), i = r(a, i), t(a).match(/br|hr|img|input/) || n.push({ event: "stop", offset: i, node: a }));
      return i
    }(e, 0), n
  }

  function c(e, r, i) {
    function a() {
      return e.length && r.length ? e[0].offset !== r[0].offset ? e[0].offset < r[0].offset ? e : r : "start" === r[0].event ? e : r : e.length ? e : r
    }

    function o(e) {
      function r(e) {
        return " " + e.nodeName + '="' + n(e.value) + '"'
      }
      l += "<" + t(e) + w.map.call(e.attributes, r).join("") + ">"
    }

    function u(e) { l += "</" + t(e) + ">" }

    function c(e) {
      ("start" === e.event ? o : u)(e.node)
    }
    for (var s = 0, l = "", f = []; e.length || r.length;) {
      var g = a();
      if (l += n(i.substring(s, g[0].offset)), s = g[0].offset, g === e) {
        f.reverse().forEach(u);
        do c(g.splice(0, 1)[0]), g = a(); while (g === e && g.length && g[0].offset === s);
        f.reverse().forEach(o)
      } else "start" === g[0].event ? f.push(g[0].node) : f.pop(), c(g.splice(0, 1)[0])
    }
    return l + n(i.substr(s))
  }

  function s(e) {
    function n(e) {
      return e && e.source || e
    }

    function t(t, r) {
      return new RegExp(n(t), "m" + (e.cI ? "i" : "") + (r ? "g" : ""))
    }

    function r(i, a) {
      if (!i.compiled) {
        if (i.compiled = !0, i.k = i.k || i.bK, i.k) {
          var u = {},
            c = function(n, t) {
              e.cI && (t = t.toLowerCase()), t.split(" ").forEach(function(e) {
                var t = e.split("|");
                u[t[0]] = [n, t[1] ? Number(t[1]) : 1]
              })
            };
          "string" == typeof i.k ? c("keyword", i.k) : E(i.k).forEach(function(e) { c(e, i.k[e]) }), i.k = u
        }
        i.lR = t(i.l || /\w+/, !0), a && (i.bK && (i.b = "\\b(" + i.bK.split(" ").join("|") + ")\\b"), i.b || (i.b = /\B|\b/), i.bR = t(i.b), i.e || i.eW || (i.e = /\B|\b/), i.e && (i.eR = t(i.e)), i.tE = n(i.e) || "", i.eW && a.tE && (i.tE += (i.e ? "|" : "") + a.tE)), i.i && (i.iR = t(i.i)), null == i.r && (i.r = 1), i.c || (i.c = []);
        var s = [];
        i.c.forEach(function(e) { e.v ? e.v.forEach(function(n) { s.push(o(e, n)) }) : s.push("self" === e ? i : e) }), i.c = s, i.c.forEach(function(e) { r(e, i) }), i.starts && r(i.starts, a);
        var l = i.c.map(function(e) {
          return e.bK ? "\\.?(" + e.b + ")\\.?" : e.b
        }).concat([i.tE, i.i]).map(n).filter(Boolean);
        i.t = l.length ? t(l.join("|"), !0) : {
          exec: function() {
            return null
          }
        }
      }
    }
    r(e)
  }

  function l(e, t, i, a) {
    function o(e, n) {
      var t, i;
      for (t = 0, i = n.c.length; i > t; t++)
        if (r(n.c[t].bR, e)) return n.c[t]
    }

    function u(e, n) {
      if (r(e.eR, n)) {
        for (; e.endsParent && e.parent;) e = e.parent;
        return e
      }
      return e.eW ? u(e.parent, n) : void 0
    }

    function c(e, n) {
      return !i && r(n.iR, e)
    }

    function g(e, n) {
      var t = N.cI ? n[0].toLowerCase() : n[0];
      return e.k.hasOwnProperty(t) && e.k[t]
    }

    function h(e, n, t, r) {
      var i = r ? "" : y.classPrefix,
        a = '<span class="' + i,
        o = t ? "" : C;
      return a += e + '">', a + n + o
    }

    function p() {
      var e, t, r, i;
      if (!E.k) return n(B);
      for (i = "", t = 0, E.lR.lastIndex = 0, r = E.lR.exec(B); r;) i += n(B.substring(t, r.index)), e = g(E, r), e ? (M += e[1], i += h(e[0], n(r[0]))) : i += n(r[0]), t = E.lR.lastIndex, r = E.lR.exec(B);
      return i + n(B.substr(t))
    }

    function d() {
      var e = "string" == typeof E.sL;
      if (e && !x[E.sL]) return n(B);
      var t = e ? l(E.sL, B, !0, L[E.sL]) : f(B, E.sL.length ? E.sL : void 0);
      return E.r > 0 && (M += t.r), e && (L[E.sL] = t.top), h(t.language, t.value, !1, !0)
    }

    function b() { k += null != E.sL ? d() : p(), B = "" }

    function v(e) { k += e.cN ? h(e.cN, "", !0) : "", E = Object.create(e, { parent: { value: E } }) }

    function m(e, n) {
      if (B += e, null == n) return b(), 0;
      var t = o(n, E);
      if (t) return t.skip ? B += n : (t.eB && (B += n), b(), t.rB || t.eB || (B = n)), v(t, n), t.rB ? 0 : n.length;
      var r = u(E, n);
      if (r) {
        var i = E;
        i.skip ? B += n : (i.rE || i.eE || (B += n), b(), i.eE && (B = n));
        do E.cN && (k += C), E.skip || (M += E.r), E = E.parent; while (E !== r.parent);
        return r.starts && v(r.starts, ""), i.rE ? 0 : n.length
      }
      if (c(n, E)) throw new Error('Illegal lexeme "' + n + '" for mode "' + (E.cN || "<unnamed>") + '"');
      return B += n, n.length || 1
    }
    var N = R(e);
    if (!N) throw new Error('Unknown language: "' + e + '"');
    s(N);
    var w, E = a || N,
      L = {},
      k = "";
    for (w = E; w !== N; w = w.parent) w.cN && (k = h(w.cN, "", !0) + k);
    var B = "",
      M = 0;
    try {
      for (var I, j, O = 0;;) {
        if (E.t.lastIndex = O, I = E.t.exec(t), !I) break;
        j = m(t.substring(O, I.index), I[0]), O = I.index + j
      }
      for (m(t.substr(O)), w = E; w.parent; w = w.parent) w.cN && (k += C);
      return { r: M, value: k, language: e, top: E }
    } catch (T) {
      if (T.message && -1 !== T.message.indexOf("Illegal")) return { r: 0, value: n(t) };
      throw T
    }
  }

  function f(e, t) {
    t = t || y.languages || E(x);
    var r = { r: 0, value: n(e) },
      i = r;
    return t.filter(R).forEach(function(n) {
      var t = l(n, e, !1);
      t.language = n, t.r > i.r && (i = t), t.r > r.r && (i = r, r = t)
    }), i.language && (r.second_best = i), r
  }

  function g(e) {
    return y.tabReplace || y.useBR ? e.replace(M, function(e, n) {
      return y.useBR && "\n" === e ? "<br>" : y.tabReplace ? n.replace(/\t/g, y.tabReplace) : void 0
    }) : e
  }

  function h(e, n, t) {
    var r = n ? L[n] : t,
      i = [e.trim()];
    return e.match(/\bhljs\b/) || i.push("hljs"), -1 === e.indexOf(r) && i.push(r), i.join(" ").trim()
  }

  function p(e) {
    var n, t, r, o, s, p = a(e);
    i(p) || (y.useBR ? (n = document.createElementNS("http://www.w3.org/1999/xhtml", "div"), n.innerHTML = e.innerHTML.replace(/\n/g, "").replace(/<br[ \/]*>/g, "\n")) : n = e, s = n.textContent, r = p ? l(p, s, !0) : f(s), t = u(n), t.length && (o = document.createElementNS("http://www.w3.org/1999/xhtml", "div"), o.innerHTML = r.value, r.value = c(t, u(o), s)), r.value = g(r.value), e.innerHTML = r.value, e.className = h(e.className, p, r.language), e.result = { language: r.language, re: r.r }, r.second_best && (e.second_best = { language: r.second_best.language, re: r.second_best.r }))
  }

  function d(e) { y = o(y, e) }

  function b() {
    if (!b.called) {
      b.called = !0;
      var e = document.querySelectorAll("pre code");
      w.forEach.call(e, p)
    }
  }

  function v() { addEventListener("DOMContentLoaded", b, !1), addEventListener("load", b, !1) }

  function m(n, t) {
    var r = x[n] = t(e);
    r.aliases && r.aliases.forEach(function(e) { L[e] = n })
  }

  function N() {
    return E(x)
  }

  function R(e) {
    return e = (e || "").toLowerCase(), x[e] || x[L[e]]
  }
  var w = [],
    E = Object.keys,
    x = {},
    L = {},
    k = /^(no-?highlight|plain|text)$/i,
    B = /\blang(?:uage)?-([\w-]+)\b/i,
    M = /((^(<[^>]+>|\t|)+|(?:\n)))/gm,
    C = "</span>",
    y = { classPrefix: "hljs-", tabReplace: null, useBR: !1, languages: void 0 },
    I = { "&": "&amp;", "<": "&lt;", ">": "&gt;" };
  return e.highlight = l, e.highlightAuto = f, e.fixMarkup = g, e.highlightBlock = p, e.configure = d, e.initHighlighting = b, e.initHighlightingOnLoad = v, e.registerLanguage = m, e.listLanguages = N, e.getLanguage = R, e.inherit = o, e.IR = "[a-zA-Z]\\w*", e.UIR = "[a-zA-Z_]\\w*", e.NR = "\\b\\d+(\\.\\d+)?", e.CNR = "(-?)(\\b0[xX][a-fA-F0-9]+|(\\b\\d+(\\.\\d*)?|\\.\\d+)([eE][-+]?\\d+)?)", e.BNR = "\\b(0b[01]+)", e.RSR = "!|!=|!==|%|%=|&|&&|&=|\\*|\\*=|\\+|\\+=|,|-|-=|/=|/|:|;|<<|<<=|<=|<|===|==|=|>>>=|>>=|>=|>>>|>>|>|\\?|\\[|\\{|\\(|\\^|\\^=|\\||\\|=|\\|\\||~", e.BE = { b: "\\\\[\\s\\S]", r: 0 }, e.ASM = { cN: "string", b: "'", e: "'", i: "\\n", c: [e.BE] }, e.QSM = { cN: "string", b: '"', e: '"', i: "\\n", c: [e.BE] }, e.PWM = { b: /\b(a|an|the|are|I'm|isn't|don't|doesn't|won't|but|just|should|pretty|simply|enough|gonna|going|wtf|so|such|will|you|your|like)\b/ }, e.C = function(n, t, r) {
    var i = e.inherit({ cN: "comment", b: n, e: t, c: [] }, r || {});
    return i.c.push(e.PWM), i.c.push({ cN: "doctag", b: "(?:TODO|FIXME|NOTE|BUG|XXX):", r: 0 }), i
  }, e.CLCM = e.C("//", "$"), e.CBCM = e.C("/\\*", "\\*/"), e.HCM = e.C("#", "$"), e.NM = { cN: "number", b: e.NR, r: 0 }, e.CNM = { cN: "number", b: e.CNR, r: 0 }, e.BNM = { cN: "number", b: e.BNR, r: 0 }, e.CSSNM = { cN: "number", b: e.NR + "(%|em|ex|ch|rem|vw|vh|vmin|vmax|cm|mm|in|pt|pc|px|deg|grad|rad|turn|s|ms|Hz|kHz|dpi|dpcm|dppx)?", r: 0 }, e.RM = { cN: "regexp", b: /\//, e: /\/[gimuy]*/, i: /\n/, c: [e.BE, { b: /\[/, e: /\]/, r: 0, c: [e.BE] }] }, e.TM = { cN: "title", b: e.IR, r: 0 }, e.UTM = { cN: "title", b: e.UIR, r: 0 }, e.METHOD_GUARD = { b: "\\.\\s*" + e.UIR, r: 0 }, e
});
hljs.registerLanguage("xml", function(s) {
  var e = "[A-Za-z0-9\\._:-]+",
    t = { eW: !0, i: /</, r: 0, c: [{ cN: "attr", b: e, r: 0 }, { b: /=\s*/, r: 0, c: [{ cN: "string", endsParent: !0, v: [{ b: /"/, e: /"/ }, { b: /'/, e: /'/ }, { b: /[^\s"'=<>`]+/ }] }] }] };
  return { aliases: ["html", "xhtml", "rss", "atom", "xjb", "xsd", "xsl", "plist"], cI: !0, c: [{ cN: "meta", b: "<!DOCTYPE", e: ">", r: 10, c: [{ b: "\\[", e: "\\]" }] }, s.C("<!--", "-->", { r: 10 }), { b: "<\\!\\[CDATA\\[", e: "\\]\\]>", r: 10 }, { b: /<\?(php)?/, e: /\?>/, sL: "php", c: [{ b: "/\\*", e: "\\*/", skip: !0 }] }, { cN: "tag", b: "<style(?=\\s|>|$)", e: ">", k: { name: "style" }, c: [t], starts: { e: "</style>", rE: !0, sL: ["css", "xml"] } }, { cN: "tag", b: "<script(?=\\s|>|$)", e: ">", k: { name: "script" }, c: [t], starts: { e: "</script>", rE: !0, sL: ["actionscript", "javascript", "handlebars", "xml"] } }, { cN: "meta", v: [{ b: /<\?xml/, e: /\?>/, r: 10 }, { b: /<\?\w+/, e: /\?>/ }] }, { cN: "tag", b: "</?", e: "/?>", c: [{ cN: "name", b: /[^\/><\s]+/, r: 0 }, t] }] }
});
hljs.registerLanguage("go", function(e) {
  var t = { keyword: "break default func interface select case map struct chan else goto package switch const fallthrough if range end type continue for import return var go defer bool byte complex64 complex128 float32 float64 int8 int16 int32 int64 string uint8 uint16 uint32 uint64 int uint uintptr rune id autoplay Get", literal: "true false iota nil Pages with", built_in: "append cap close complex highlight copy imag len make new panic print println real recover delete Site Data tweet speakerdeck youtube ref relref vimeo instagram gist figure innershortcode" };
  return { aliases: ["golang"], k: t, i: "</", c: [e.CLCM, e.CBCM, { cN: "string", v: [e.QSM, { b: "'", e: "[^\\\\]'" }, { b: "`", e: "`" }] }, { cN: "number", v: [{ b: e.CNR + "[dflsi]", r: 1 }, e.CNM] }, { b: /:=/ }, { cN: "function", bK: "func", e: /\s*\{/, eE: !0, c: [e.TM, { cN: "params", b: /\(/, e: /\)/, k: t, i: /["']/ }] }] }
});
hljs.registerLanguage("markdown", function(e) {
  return { aliases: ["md", "mkdown", "mkd"], c: [{ cN: "section", v: [{ b: "^#{1,6}", e: "$" }, { b: "^.+?\\n[=-]{2,}$" }] }, { b: "<", e: ">", sL: "xml", r: 0 }, { cN: "bullet", b: "^([*+-]|(\\d+\\.))\\s+" }, { cN: "strong", b: "[*_]{2}.+?[*_]{2}" }, { cN: "emphasis", v: [{ b: "\\*.+?\\*" }, { b: "_.+?_", r: 0 }] }, { cN: "quote", b: "^>\\s+", e: "$" }, { cN: "code", v: [{ b: "^```w*s*$", e: "^```s*$" }, { b: "`.+?`" }, { b: "^( {4}|  )", e: "$", r: 0 }] }, { b: "^[-\\*]{3,}", e: "$" }, { b: "\\[.+?\\][\\(\\[].*?[\\)\\]]", rB: !0, c: [{ cN: "string", b: "\\[", e: "\\]", eB: !0, rE: !0, r: 0 }, { cN: "link", b: "\\]\\(", e: "\\)", eB: !0, eE: !0 }, { cN: "symbol", b: "\\]\\[", e: "\\]", eB: !0, eE: !0 }], r: 10 }, { b: /^\[[^\n]+\]:/, rB: !0, c: [{ cN: "symbol", b: /\[/, e: /\]/, eB: !0, eE: !0 }, { cN: "link", b: /:\s*/, e: /$/, eB: !0 }] }] }
});
hljs.registerLanguage("handlebars", function(e) {
  var a = { "builtin-name": "each in with if else unless bindattr action collection debugger log outlet template unbound view yield" };
  return { aliases: ["hbs", "html.hbs", "html.handlebars"], cI: !0, sL: "xml", c: [e.C("{{!(--)?", "(--)?}}"), { cN: "template-tag", b: /\{\{[#\/]/, e: /\}\}/, c: [{ cN: "name", b: /[a-zA-Z\.-]+/, k: a, starts: { eW: !0, r: 0, c: [e.QSM] } }] }, { cN: "template-variable", b: /\{\{/, e: /\}\}/, k: a }] }
});
hljs.registerLanguage("apache", function(e) {
  var r = { cN: "number", b: "[\\$%]\\d+" };
  return { aliases: ["apacheconf"], cI: !0, c: [e.HCM, { cN: "section", b: "</?", e: ">" }, { cN: "attribute", b: /\w+/, r: 0, k: { nomarkup: "order deny allow setenv rewriterule rewriteengine rewritecond documentroot sethandler errordocument loadmodule options header listen serverroot servername" }, starts: { e: /$/, r: 0, k: { literal: "on off all" }, c: [{ cN: "meta", b: "\\s\\[", e: "\\]$" }, { cN: "variable", b: "[\\$%]\\{", e: "\\}", c: ["self", r] }, r, e.QSM] } }], i: /\S/ }
});
hljs.registerLanguage("ini", function(e) {
  var b = { cN: "string", c: [e.BE], v: [{ b: "'''", e: "'''", r: 10 }, { b: '"""', e: '"""', r: 10 }, { b: '"', e: '"' }, { b: "'", e: "'" }] };
  return { aliases: ["toml"], cI: !0, i: /\S/, c: [e.C(";", "$"), e.HCM, { cN: "section", b: /^\s*\[+/, e: /\]+/ }, { b: /^[a-z0-9\[\]_-]+\s*=\s*/, e: "$", rB: !0, c: [{ cN: "attr", b: /[a-z0-9\[\]_-]+/ }, { b: /=/, eW: !0, r: 0, c: [{ cN: "literal", b: /\bon|off|true|false|yes|no\b/ }, { cN: "variable", v: [{ b: /\$[\w\d"][\w\d_]*/ }, { b: /\$\{(.*?)}/ }] }, b, { cN: "number", b: /([\+\-]+)?[\d]+_[\d_]+/ }, e.NM] }] }] }
});
hljs.registerLanguage("css", function(e) {
  var c = "[a-zA-Z-][a-zA-Z0-9_-]*",
    t = { b: /[A-Z\_\.\-]+\s*:/, rB: !0, e: ";", eW: !0, c: [{ cN: "attribute", b: /\S/, e: ":", eE: !0, starts: { eW: !0, eE: !0, c: [{ b: /[\w-]+\(/, rB: !0, c: [{ cN: "built_in", b: /[\w-]+/ }, { b: /\(/, e: /\)/, c: [e.ASM, e.QSM] }] }, e.CSSNM, e.QSM, e.ASM, e.CBCM, { cN: "number", b: "#[0-9A-Fa-f]+" }, { cN: "meta", b: "!important" }] } }] };
  return { cI: !0, i: /[=\/|'\$]/, c: [e.CBCM, { cN: "selector-id", b: /#[A-Za-z0-9_-]+/ }, { cN: "selector-class", b: /\.[A-Za-z0-9_-]+/ }, { cN: "selector-attr", b: /\[/, e: /\]/, i: "$" }, { cN: "selector-pseudo", b: /:(:)?[a-zA-Z0-9\_\-\+\(\)"'.]+/ }, { b: "@(font-face|page)", l: "[a-z-]+", k: "font-face page" }, { b: "@", e: "[{;]", i: /:/, c: [{ cN: "keyword", b: /\w+/ }, { b: /\s/, eW: !0, eE: !0, r: 0, c: [e.ASM, e.QSM, e.CSSNM] }] }, { cN: "selector-tag", b: c, r: 0 }, { b: "{", e: "}", i: /\S/, c: [e.CBCM, t] }] }
});
hljs.registerLanguage("asciidoc", function(e) {
  return { aliases: ["adoc"], c: [e.C("^/{4,}\\n", "\\n/{4,}$", { r: 10 }), e.C("^//", "$", { r: 0 }), { cN: "title", b: "^\\.\\w.*$" }, { b: "^[=\\*]{4,}\\n", e: "\\n^[=\\*]{4,}$", r: 10 }, { cN: "section", r: 10, v: [{ b: "^(={1,5}) .+?( \\1)?$" }, { b: "^[^\\[\\]\\n]+?\\n[=\\-~\\^\\+]{2,}$" }] }, { cN: "meta", b: "^:.+?:", e: "\\s", eE: !0, r: 10 }, { cN: "meta", b: "^\\[.+?\\]$", r: 0 }, { cN: "quote", b: "^_{4,}\\n", e: "\\n_{4,}$", r: 10 }, { cN: "code", b: "^[\\-\\.]{4,}\\n", e: "\\n[\\-\\.]{4,}$", r: 10 }, { b: "^\\+{4,}\\n", e: "\\n\\+{4,}$", c: [{ b: "<", e: ">", sL: "xml", r: 0 }], r: 10 }, { cN: "bullet", b: "^(\\*+|\\-+|\\.+|[^\\n]+?::)\\s+" }, { cN: "symbol", b: "^(NOTE|TIP|IMPORTANT|WARNING|CAUTION):\\s+", r: 10 }, { cN: "strong", b: "\\B\\*(?![\\*\\s])", e: "(\\n{2}|\\*)", c: [{ b: "\\\\*\\w", r: 0 }] }, { cN: "emphasis", b: "\\B'(?!['\\s])", e: "(\\n{2}|')", c: [{ b: "\\\\'\\w", r: 0 }], r: 0 }, { cN: "emphasis", b: "_(?![_\\s])", e: "(\\n{2}|_)", r: 0 }, { cN: "string", v: [{ b: "``.+?''" }, { b: "`.+?'" }] }, { cN: "code", b: "(`.+?`|\\+.+?\\+)", r: 0 }, { cN: "code", b: "^[ \\t]", e: "$", r: 0 }, { b: "^'{3,}[ \\t]*$", r: 10 }, { b: "(link:)?(http|https|ftp|file|irc|image:?):\\S+\\[.*?\\]", rB: !0, c: [{ b: "(link|image:?):", r: 0 }, { cN: "link", b: "\\w", e: "[^\\[]+", r: 0 }, { cN: "string", b: "\\[", e: "\\]", eB: !0, eE: !0, r: 0 }], r: 10 }] }
});
hljs.registerLanguage("ruby", function(e) {
  var b = "[a-zA-Z_]\\w*[!?=]?|[-+~]\\@|<<|>>|=~|===?|<=>|[<>]=?|\\*\\*|[-/+%^&*~`|]|\\[\\]=?",
    r = { keyword: "and then defined module in return redo if BEGIN retry end for self when next until do begin unless END rescue else break undef not super class case require yield alias while ensure elsif or include attr_reader attr_writer attr_accessor", literal: "true false nil" },
    c = { cN: "doctag", b: "@[A-Za-z]+" },
    a = { b: "#<", e: ">" },
    s = [e.C("#", "$", { c: [c] }), e.C("^\\=begin", "^\\=end", { c: [c], r: 10 }), e.C("^__END__", "\\n$")],
    n = { cN: "subst", b: "#\\{", e: "}", k: r },
    t = { cN: "string", c: [e.BE, n], v: [{ b: /'/, e: /'/ }, { b: /"/, e: /"/ }, { b: /`/, e: /`/ }, { b: "%[qQwWx]?\\(", e: "\\)" }, { b: "%[qQwWx]?\\[", e: "\\]" }, { b: "%[qQwWx]?{", e: "}" }, { b: "%[qQwWx]?<", e: ">" }, { b: "%[qQwWx]?/", e: "/" }, { b: "%[qQwWx]?%", e: "%" }, { b: "%[qQwWx]?-", e: "-" }, { b: "%[qQwWx]?\\|", e: "\\|" }, { b: /\B\?(\\\d{1,3}|\\x[A-Fa-f0-9]{1,2}|\\u[A-Fa-f0-9]{4}|\\?\S)\b/ }, { b: /<<(-?)\w+$/, e: /^\s*\w+$/ }] },
    i = { cN: "params", b: "\\(", e: "\\)", endsParent: !0, k: r },
    d = [t, a, { cN: "class", bK: "class module", e: "$|;", i: /=/, c: [e.inherit(e.TM, { b: "[A-Za-z_]\\w*(::\\w+)*(\\?|\\!)?" }), { b: "<\\s*", c: [{ b: "(" + e.IR + "::)?" + e.IR }] }].concat(s) }, { cN: "function", bK: "def", e: "$|;", c: [e.inherit(e.TM, { b: b }), i].concat(s) }, { b: e.IR + "::" }, { cN: "symbol", b: e.UIR + "(\\!|\\?)?:", r: 0 }, { cN: "symbol", b: ":(?!\\s)", c: [t, { b: b }], r: 0 }, { cN: "number", b: "(\\b0[0-7_]+)|(\\b0x[0-9a-fA-F_]+)|(\\b[1-9][0-9_]*(\\.[0-9_]+)?)|[0_]\\b", r: 0 }, { b: "(\\$\\W)|((\\$|\\@\\@?)(\\w+))" }, { cN: "params", b: /\|/, e: /\|/, k: r }, { b: "(" + e.RSR + "|unless)\\s*", c: [a, { cN: "regexp", c: [e.BE, n], i: /\n/, v: [{ b: "/", e: "/[a-z]*" }, { b: "%r{", e: "}[a-z]*" }, { b: "%r\\(", e: "\\)[a-z]*" }, { b: "%r!", e: "![a-z]*" }, { b: "%r\\[", e: "\\][a-z]*" }] }].concat(s), r: 0 }].concat(s);
  n.c = d, i.c = d;
  var l = "[>?]>",
    o = "[\\w#]+\\(\\w+\\):\\d+:\\d+>",
    u = "(\\w+-)?\\d+\\.\\d+\\.\\d(p\\d+)?[^>]+>",
    w = [{ b: /^\s*=>/, starts: { e: "$", c: d } }, { cN: "meta", b: "^(" + l + "|" + o + "|" + u + ")", starts: { e: "$", c: d } }];
  return { aliases: ["rb", "gemspec", "podspec", "thor", "irb"], k: r, i: /\/\*/, c: s.concat(w).concat(d) }
});
hljs.registerLanguage("yaml", function(e) {
  var a = { literal: "{ } true false yes no Yes No True False null" },
    b = "^[ \\-]*",
    r = "[a-zA-Z_][\\w\\-]*",
    t = { cN: "attr", v: [{ b: b + r + ":" }, { b: b + '"' + r + '":' }, { b: b + "'" + r + "':" }] },
    c = { cN: "template-variable", v: [{ b: "{{", e: "}}" }, { b: "%{", e: "}" }] },
    l = { cN: "string", r: 0, v: [{ b: /'/, e: /'/ }, { b: /"/, e: /"/ }], c: [e.BE, c] };
  return { cI: !0, aliases: ["yml", "YAML", "yaml"], c: [t, { cN: "meta", b: "^---s*$", r: 10 }, { cN: "string", b: "[\\|>] *$", rE: !0, c: l.c, e: t.v[0].b }, { b: "<%[%=-]?", e: "[%-]?%>", sL: "ruby", eB: !0, eE: !0, r: 0 }, { cN: "type", b: "!!" + e.UIR }, { cN: "meta", b: "&" + e.UIR + "$" }, { cN: "meta", b: "\\*" + e.UIR + "$" }, { cN: "bullet", b: "^ *-", r: 0 }, l, e.HCM, e.CNM], k: a }
});
hljs.registerLanguage("powershell", function(e) {
  var t = { b: "`[\\s\\S]", r: 0 },
    o = { cN: "variable", v: [{ b: /\$[\w\d][\w\d_:]*/ }] },
    r = { cN: "literal", b: /\$(null|true|false)\b/ },
    n = { cN: "string", v: [{ b: /"/, e: /"/ }, { b: /@"/, e: /^"@/ }], c: [t, o, { cN: "variable", b: /\$[A-z]/, e: /[^A-z]/ }] },
    a = { cN: "string", v: [{ b: /'/, e: /'/ }, { b: /@'/, e: /^'@/ }] },
    i = { cN: "doctag", v: [{ b: /\.(synopsis|description|example|inputs|outputs|notes|link|component|role|functionality)/ }, { b: /\.(parameter|forwardhelptargetname|forwardhelpcategory|remotehelprunspace|externalhelp)\s+\S+/ }] },
    s = e.inherit(e.C(null, null), { v: [{ b: /#/, e: /$/ }, { b: /<#/, e: /#>/ }], c: [i] });
  return { aliases: ["ps"], l: /-?[A-z\.\-]+/, cI: !0, k: { keyword: "if else foreach return function do while until elseif begin for trap data dynamicparam end break throw param continue finally in switch exit filter try process catch", built_in: "Add-Computer Add-Content Add-History Add-JobTrigger Add-Member Add-PSSnapin Add-Type Checkpoint-Computer Clear-Content Clear-EventLog Clear-History Clear-Host Clear-Item Clear-ItemProperty Clear-Variable Compare-Object Complete-Transaction Connect-PSSession Connect-WSMan Convert-Path ConvertFrom-Csv ConvertFrom-Json ConvertFrom-SecureString ConvertFrom-StringData ConvertTo-Csv ConvertTo-Html ConvertTo-Json ConvertTo-SecureString ConvertTo-Xml Copy-Item Copy-ItemProperty Debug-Process Disable-ComputerRestore Disable-JobTrigger Disable-PSBreakpoint Disable-PSRemoting Disable-PSSessionConfiguration Disable-WSManCredSSP Disconnect-PSSession Disconnect-WSMan Disable-ScheduledJob Enable-ComputerRestore Enable-JobTrigger Enable-PSBreakpoint Enable-PSRemoting Enable-PSSessionConfiguration Enable-ScheduledJob Enable-WSManCredSSP Enter-PSSession Exit-PSSession Export-Alias Export-Clixml Export-Console Export-Counter Export-Csv Export-FormatData Export-ModuleMember Export-PSSession ForEach-Object Format-Custom Format-List Format-Table Format-Wide Get-Acl Get-Alias Get-AuthenticodeSignature Get-ChildItem Get-Command Get-ComputerRestorePoint Get-Content Get-ControlPanelItem Get-Counter Get-Credential Get-Culture Get-Date Get-Event Get-EventLog Get-EventSubscriber Get-ExecutionPolicy Get-FormatData Get-Host Get-HotFix Get-Help Get-History Get-IseSnippet Get-Item Get-ItemProperty Get-Job Get-JobTrigger Get-Location Get-Member Get-Module Get-PfxCertificate Get-Process Get-PSBreakpoint Get-PSCallStack Get-PSDrive Get-PSProvider Get-PSSession Get-PSSessionConfiguration Get-PSSnapin Get-Random Get-ScheduledJob Get-ScheduledJobOption Get-Service Get-TraceSource Get-Transaction Get-TypeData Get-UICulture Get-Unique Get-Variable Get-Verb Get-WinEvent Get-WmiObject Get-WSManCredSSP Get-WSManInstance Group-Object Import-Alias Import-Clixml Import-Counter Import-Csv Import-IseSnippet Import-LocalizedData Import-PSSession Import-Module Invoke-AsWorkflow Invoke-Command Invoke-Expression Invoke-History Invoke-Item Invoke-RestMethod Invoke-WebRequest Invoke-WmiMethod Invoke-WSManAction Join-Path Limit-EventLog Measure-Command Measure-Object Move-Item Move-ItemProperty New-Alias New-Event New-EventLog New-IseSnippet New-Item New-ItemProperty New-JobTrigger New-Object New-Module New-ModuleManifest New-PSDrive New-PSSession New-PSSessionConfigurationFile New-PSSessionOption New-PSTransportOption New-PSWorkflowExecutionOption New-PSWorkflowSession New-ScheduledJobOption New-Service New-TimeSpan New-Variable New-WebServiceProxy New-WinEvent New-WSManInstance New-WSManSessionOption Out-Default Out-File Out-GridView Out-Host Out-Null Out-Printer Out-String Pop-Location Push-Location Read-Host Receive-Job Register-EngineEvent Register-ObjectEvent Register-PSSessionConfiguration Register-ScheduledJob Register-WmiEvent Remove-Computer Remove-Event Remove-EventLog Remove-Item Remove-ItemProperty Remove-Job Remove-JobTrigger Remove-Module Remove-PSBreakpoint Remove-PSDrive Remove-PSSession Remove-PSSnapin Remove-TypeData Remove-Variable Remove-WmiObject Remove-WSManInstance Rename-Computer Rename-Item Rename-ItemProperty Reset-ComputerMachinePassword Resolve-Path Restart-Computer Restart-Service Restore-Computer Resume-Job Resume-Service Save-Help Select-Object Select-String Select-Xml Send-MailMessage Set-Acl Set-Alias Set-AuthenticodeSignature Set-Content Set-Date Set-ExecutionPolicy Set-Item Set-ItemProperty Set-JobTrigger Set-Location Set-PSBreakpoint Set-PSDebug Set-PSSessionConfiguration Set-ScheduledJob Set-ScheduledJobOption Set-Service Set-StrictMode Set-TraceSource Set-Variable Set-WmiInstance Set-WSManInstance Set-WSManQuickConfig Show-Command Show-ControlPanelItem Show-EventLog Sort-Object Split-Path Start-Job Start-Process Start-Service Start-Sleep Start-Transaction Start-Transcript Stop-Computer Stop-Job Stop-Process Stop-Service Stop-Transcript Suspend-Job Suspend-Service Tee-Object Test-ComputerSecureChannel Test-Connection Test-ModuleManifest Test-Path Test-PSSessionConfigurationFile Trace-Command Unblock-File Undo-Transaction Unregister-Event Unregister-PSSessionConfiguration Unregister-ScheduledJob Update-FormatData Update-Help Update-List Update-TypeData Use-Transaction Wait-Event Wait-Job Wait-Process Where-Object Write-Debug Write-Error Write-EventLog Write-Host Write-Output Write-Progress Write-Verbose Write-Warning Add-MDTPersistentDrive Disable-MDTMonitorService Enable-MDTMonitorService Get-MDTDeploymentShareStatistics Get-MDTMonitorData Get-MDTOperatingSystemCatalog Get-MDTPersistentDrive Import-MDTApplication Import-MDTDriver Import-MDTOperatingSystem Import-MDTPackage Import-MDTTaskSequence New-MDTDatabase Remove-MDTMonitorData Remove-MDTPersistentDrive Restore-MDTPersistentDrive Set-MDTMonitorData Test-MDTDeploymentShare Test-MDTMonitorData Update-MDTDatabaseSchema Update-MDTDeploymentShare Update-MDTLinkedDS Update-MDTMedia Update-MDTMedia Add-VamtProductKey Export-VamtData Find-VamtManagedMachine Get-VamtConfirmationId Get-VamtProduct Get-VamtProductKey Import-VamtData Initialize-VamtData Install-VamtConfirmationId Install-VamtProductActivation Install-VamtProductKey Update-VamtProduct", nomarkup: "-ne -eq -lt -gt -ge -le -not -like -notlike -match -notmatch -contains -notcontains -in -notin -replace" }, c: [t, e.NM, n, a, r, o, s] }
});
hljs.registerLanguage("scss", function(e) {
  var t = "[a-zA-Z-][a-zA-Z0-9_-]*",
    i = { cN: "variable", b: "(\\$" + t + ")\\b" },
    r = { cN: "number", b: "#[0-9A-Fa-f]+" };
  ({ cN: "attribute", b: "[A-Z\\_\\.\\-]+", e: ":", eE: !0, i: "[^\\s]", starts: { eW: !0, eE: !0, c: [r, e.CSSNM, e.QSM, e.ASM, e.CBCM, { cN: "meta", b: "!important" }] } });
  return { cI: !0, i: "[=/|']", c: [e.CLCM, e.CBCM, { cN: "selector-id", b: "\\#[A-Za-z0-9_-]+", r: 0 }, { cN: "selector-class", b: "\\.[A-Za-z0-9_-]+", r: 0 }, { cN: "selector-attr", b: "\\[", e: "\\]", i: "$" }, { cN: "selector-tag", b: "\\b(a|abbr|acronym|address|area|article|aside|audio|b|base|big|blockquote|body|br|button|canvas|caption|cite|code|col|colgroup|command|datalist|dd|del|details|dfn|div|dl|dt|em|embed|fieldset|figcaption|figure|footer|form|frame|frameset|(h[1-6])|head|header|hgroup|hr|html|i|iframe|img|input|ins|kbd|keygen|label|legend|li|link|map|mark|meta|meter|nav|noframes|noscript|object|ol|optgroup|option|output|p|param|pre|progress|q|rp|rt|ruby|samp|script|section|select|small|span|strike|strong|style|sub|sup|table|tbody|td|textarea|tfoot|th|thead|time|title|tr|tt|ul|var|video)\\b", r: 0 }, { b: ":(visited|valid|root|right|required|read-write|read-only|out-range|optional|only-of-type|only-child|nth-of-type|nth-last-of-type|nth-last-child|nth-child|not|link|left|last-of-type|last-child|lang|invalid|indeterminate|in-range|hover|focus|first-of-type|first-line|first-letter|first-child|first|enabled|empty|disabled|default|checked|before|after|active)" }, { b: "::(after|before|choices|first-letter|first-line|repeat-index|repeat-item|selection|value)" }, i, { cN: "attribute", b: "\\b(z-index|word-wrap|word-spacing|word-break|width|widows|white-space|visibility|vertical-align|unicode-bidi|transition-timing-function|transition-property|transition-duration|transition-delay|transition|transform-style|transform-origin|transform|top|text-underline-position|text-transform|text-shadow|text-rendering|text-overflow|text-indent|text-decoration-style|text-decoration-line|text-decoration-color|text-decoration|text-align-last|text-align|tab-size|table-layout|right|resize|quotes|position|pointer-events|perspective-origin|perspective|page-break-inside|page-break-before|page-break-after|padding-top|padding-right|padding-left|padding-bottom|padding|overflow-y|overflow-x|overflow-wrap|overflow|outline-width|outline-style|outline-offset|outline-color|outline|orphans|order|opacity|object-position|object-fit|normal|none|nav-up|nav-right|nav-left|nav-index|nav-down|min-width|min-height|max-width|max-height|mask|marks|margin-top|margin-right|margin-left|margin-bottom|margin|list-style-type|list-style-position|list-style-image|list-style|line-height|letter-spacing|left|justify-content|initial|inherit|ime-mode|image-orientation|image-resolution|image-rendering|icon|hyphens|height|font-weight|font-variant-ligatures|font-variant|font-style|font-stretch|font-size-adjust|font-size|font-language-override|font-kerning|font-feature-settings|font-family|font|float|flex-wrap|flex-shrink|flex-grow|flex-flow|flex-direction|flex-basis|flex|filter|empty-cells|display|direction|cursor|counter-reset|counter-increment|content|column-width|column-span|column-rule-width|column-rule-style|column-rule-color|column-rule|column-gap|column-fill|column-count|columns|color|clip-path|clip|clear|caption-side|break-inside|break-before|break-after|box-sizing|box-shadow|box-decoration-break|bottom|border-width|border-top-width|border-top-style|border-top-right-radius|border-top-left-radius|border-top-color|border-top|border-style|border-spacing|border-right-width|border-right-style|border-right-color|border-right|border-radius|border-left-width|border-left-style|border-left-color|border-left|border-image-width|border-image-source|border-image-slice|border-image-repeat|border-image-outset|border-image|border-color|border-collapse|border-bottom-width|border-bottom-style|border-bottom-right-radius|border-bottom-left-radius|border-bottom-color|border-bottom|border|background-size|background-repeat|background-position|background-origin|background-image|background-color|background-clip|background-attachment|background-blend-mode|background|backface-visibility|auto|animation-timing-function|animation-play-state|animation-name|animation-iteration-count|animation-fill-mode|animation-duration|animation-direction|animation-delay|animation|align-self|align-items|align-content)\\b", i: "[^\\s]" }, { b: "\\b(whitespace|wait|w-resize|visible|vertical-text|vertical-ideographic|uppercase|upper-roman|upper-alpha|underline|transparent|top|thin|thick|text|text-top|text-bottom|tb-rl|table-header-group|table-footer-group|sw-resize|super|strict|static|square|solid|small-caps|separate|se-resize|scroll|s-resize|rtl|row-resize|ridge|right|repeat|repeat-y|repeat-x|relative|progress|pointer|overline|outside|outset|oblique|nowrap|not-allowed|normal|none|nw-resize|no-repeat|no-drop|newspaper|ne-resize|n-resize|move|middle|medium|ltr|lr-tb|lowercase|lower-roman|lower-alpha|loose|list-item|line|line-through|line-edge|lighter|left|keep-all|justify|italic|inter-word|inter-ideograph|inside|inset|inline|inline-block|inherit|inactive|ideograph-space|ideograph-parenthesis|ideograph-numeric|ideograph-alpha|horizontal|hidden|help|hand|groove|fixed|ellipsis|e-resize|double|dotted|distribute|distribute-space|distribute-letter|distribute-all-lines|disc|disabled|default|decimal|dashed|crosshair|collapse|col-resize|circle|char|center|capitalize|break-word|break-all|bottom|both|bolder|bold|block|bidi-override|below|baseline|auto|always|all-scroll|absolute|table|table-cell)\\b" }, { b: ":", e: ";", c: [i, r, e.CSSNM, e.QSM, e.ASM, { cN: "meta", b: "!important" }] }, { b: "@", e: "[{;]", k: "mixin include extend for if else each while charset import debug media page content font-face namespace warn", c: [i, e.QSM, e.ASM, r, e.CSSNM, { b: "\\s[A-Za-z0-9_.-]+", r: 0 }] }] }
});
hljs.registerLanguage("json", function(e) {
  var i = { literal: "true false null" },
    n = [e.QSM, e.CNM],
    r = { e: ",", eW: !0, eE: !0, c: n, k: i },
    t = { b: "{", e: "}", c: [{ cN: "attr", b: /"/, e: /"/, c: [e.BE], i: "\\n" }, e.inherit(r, { b: /:/ })], i: "\\S" },
    c = { b: "\\[", e: "\\]", c: [e.inherit(r)], i: "\\S" };
  return n.splice(n.length, 0, t, c), { c: n, k: i, i: "\\S" }
});
hljs.registerLanguage("bash", function(e) {
  var t = { cN: "variable", v: [{ b: /\$[\w\d#@][\w\d_]*/ }, { b: /\$\{(.*?)}/ }] },
    s = { cN: "string", b: /"/, e: /"/, c: [e.BE, t, { cN: "variable", b: /\$\(/, e: /\)/, c: [e.BE] }] },
    a = { cN: "string", b: /'/, e: /'/ };
  return { aliases: ["sh", "zsh", "git"], l: /-?[a-z\._]+/, k: { keyword: "hugo if \| then else elif fi for while in do done case esac function yoursite.com", literal: "true -b branch false posts events authors", built_in: "draft break cd checkout continue eval exec exit export getopts hash pwd readonly new return shift test times trap umask unset alias bind builtin caller command declare echo enable help let local logout mapfile printf read readarray source type typeset ulimit unalias set shopt autoload bg bindkey bye cap chdir clone comparguments compcall compctl compdescribe compfiles compgroups compquote comptags comptry compvalues dirs disable disown echotc echoti emulate fc fg float functions getcap getln history integer jobs kill limit log noglob popd print pushd pushln rehash sched setcap setopt stat suspend ttyctl unfunction unhash unlimit unsetopt vared wait whence where which yoursite zcompile zformat zftp zle zmodload zparseopts zprof zpty zregexparse zsocket zstyle ztcp", _: "-ne -eq -lt -gt -f -d -e -s -l -a" }, c: [{ cN: "meta", b: /^#![^\n]+sh\s*$/, r: 10 }, { cN: "function", b: /\w[\w\d_]*\s*\(\s*\)\s*\{/, rB: !0, c: [e.inherit(e.TM, { b: /\w[\w\d_]*/ })], r: 0 }, e.HCM, s, a, t] }
});
hljs.registerLanguage("http", function(e) {
  var t = "HTTP/[0-9\\.]+";
  return { aliases: ["https"], i: "\\S", c: [{ b: "^" + t, e: "$", c: [{ cN: "number", b: "\\b\\d{3}\\b" }] }, { b: "^[A-Z]+ (.*?) " + t + "$", rB: !0, e: "$", c: [{ cN: "string", b: " ", e: " ", eB: !0, eE: !0 }, { b: t }, { cN: "keyword", b: "[A-Z]+" }] }, { cN: "attribute", b: "^\\w", e: ": ", eE: !0, i: "\\n|\\s|=", starts: { e: "$", r: 0 } }, { b: "\\n\\n", starts: { sL: [], eW: !0 } }] }
});
hljs.registerLanguage("javascript", function(e) {
  var r = "[A-Za-z$_][0-9A-Za-z$_]*",
    t = { keyword: "in of if for while finally var new function do return void else break catch instanceof with throw case default try this switch continue typeof delete let yield const export super debugger as async await static import from as", literal: "true false null undefined NaN Infinity", built_in: "eval isFinite isNaN parseFloat parseInt decodeURI decodeURIComponent encodeURI encodeURIComponent escape unescape Object Function Boolean Error EvalError InternalError RangeError ReferenceError StopIteration SyntaxError TypeError URIError Number Math Date String RegExp Array Float32Array Float64Array Int16Array Int32Array Int8Array Uint16Array Uint32Array Uint8Array Uint8ClampedArray ArrayBuffer DataView JSON Intl arguments require module console window document Symbol Set Map WeakSet WeakMap Proxy Reflect Promise" },
    a = { cN: "number", v: [{ b: "\\b(0[bB][01]+)" }, { b: "\\b(0[oO][0-7]+)" }, { b: e.CNR }], r: 0 },
    n = { cN: "subst", b: "\\$\\{", e: "\\}", k: t, c: [] },
    c = { cN: "string", b: "`", e: "`", c: [e.BE, n] };
  n.c = [e.ASM, e.QSM, c, a, e.RM];
  var s = n.c.concat([e.CBCM, e.CLCM]);
  return { aliases: ["js", "jsx"], k: t, c: [{ cN: "meta", r: 10, b: /^\s*['"]use (strict|asm)['"]/ }, { cN: "meta", b: /^#!/, e: /$/ }, e.ASM, e.QSM, c, e.CLCM, e.CBCM, a, { b: /[{,]\s*/, r: 0, c: [{ b: r + "\\s*:", rB: !0, r: 0, c: [{ cN: "attr", b: r, r: 0 }] }] }, { b: "(" + e.RSR + "|\\b(case|return|throw)\\b)\\s*", k: "return throw case", c: [e.CLCM, e.CBCM, e.RM, { cN: "function", b: "(\\(.*?\\)|" + r + ")\\s*=>", rB: !0, e: "\\s*=>", c: [{ cN: "params", v: [{ b: r }, { b: /\(\s*\)/ }, { b: /\(/, e: /\)/, eB: !0, eE: !0, k: t, c: s }] }] }, { b: /</, e: /(\/\w+|\w+\/)>/, sL: "xml", c: [{ b: /<\w+\s*\/>/, skip: !0 }, { b: /<\w+/, e: /(\/\w+|\w+\/)>/, skip: !0, c: [{ b: /<\w+\s*\/>/, skip: !0 }, "self"] }] }], r: 0 }, { cN: "function", bK: "function", e: /\{/, eE: !0, c: [e.inherit(e.TM, { b: r }), { cN: "params", b: /\(/, e: /\)/, eB: !0, eE: !0, c: s }], i: /\[|%/ }, { b: /\$[(.]/ }, e.METHOD_GUARD, { cN: "class", bK: "class", e: /[{;=]/, eE: !0, i: /[:"\[\]]/, c: [{ bK: "extends" }, e.UTM] }, { bK: "constructor", e: /\{/, eE: !0 }], i: /#(?!!)/ }
});
hljs.initHighlightingOnLoad();
