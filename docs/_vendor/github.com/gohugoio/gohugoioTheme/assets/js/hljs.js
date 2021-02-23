var hljs = require('highlight.js/lib/highlight.js');

hljs.registerLanguage('bash', require('highlight.js/lib/languages/bash'));
hljs.registerLanguage('css', require('highlight.js/lib/languages/css'));
hljs.registerLanguage('markdown', require('highlight.js/lib/languages/markdown'));
hljs.registerLanguage('diff', require('highlight.js/lib/languages/diff'));
// hljs.registerLanguage('go', require('highlight.js/lib/languages/go'));
hljs.registerLanguage('javascript', require('highlight.js/lib/languages/javascript'));
hljs.registerLanguage('json', require('highlight.js/lib/languages/json'));
hljs.registerLanguage('yaml', require('highlight.js/lib/languages/yaml'));
hljs.registerLanguage('xml', require('highlight.js/lib/languages/xml'));
hljs.registerLanguage('html', require('highlight.js/lib/languages/handlebars'));

hljs.registerLanguage("go", function(e) {
  var t = { keyword: "code output note warning break default func interface select case map struct chan else goto package switch const fallthrough if range end type continue for import return var go defer bool byte complex64 complex128 float32 float64 int8 int16 int32 int64 string uint8 uint16 uint32 uint64 int uint uintptr rune id autoplay Get", literal: "file download copy true false iota nil Pages with", built_in: "append cap close complex highlight copy imag len make new panic print println real recover delete Site Data tweet youtube ref relref vimeo instagram gist figure innershortcode" };
  return { aliases: ["golang","hugo"], k: t, i: "</", c: [e.CLCM, e.CBCM, { cN: "string", v: [e.QSM, { b: "'", e: "[^\\\\]'" }, { b: "`", e: "`" }] }, { cN: "number", v: [{ b: e.CNR + "[dflsi]", r: 1 }, e.CNM] }, { b: /:=/ }, { cN: "function", bK: "func", e: /\s*\{/, eE: !0, c: [e.TM, { cN: "params", b: /\(/, e: /\)/, k: t, i: /["']/ }] }] }
});

hljs.initHighlightingOnLoad();
