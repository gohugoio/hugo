/*! nouislider - 9.1.0 - 2016-12-10 16:00:32 */

! function(a) {
    "function" == typeof define && define.amd ? define([], a) : "object" == typeof exports ? module.exports = a() : window.noUiSlider = a()
}(function() {
    "use strict";

    function a(a, b) {
        var c = document.createElement("div");
        return j(c, b), a.appendChild(c), c
    }

    function b(a) {
        return a.filter(function(a) {
            return !this[a] && (this[a] = !0)
        }, {})
    }

    function c(a, b) {
        return Math.round(a / b) * b
    }

    function d(a, b) {
        var c = a.getBoundingClientRect(),
            d = a.ownerDocument,
            e = d.documentElement,
            f = m();
        return /webkit.*Chrome.*Mobile/i.test(navigator.userAgent) && (f.x = 0), b ? c.top + f.y - e.clientTop : c.left + f.x - e.clientLeft
    }

    function e(a) {
        return "number" == typeof a && !isNaN(a) && isFinite(a)
    }

    function f(a, b, c) {
        c > 0 && (j(a, b), setTimeout(function() {
            k(a, b)
        }, c))
    }

    function g(a) {
        return Math.max(Math.min(a, 100), 0)
    }

    function h(a) {
        return Array.isArray(a) ? a : [a]
    }

    function i(a) {
        a = String(a);
        var b = a.split(".");
        return b.length > 1 ? b[1].length : 0
    }

    function j(a, b) {
        a.classList ? a.classList.add(b) : a.className += " " + b
    }

    function k(a, b) {
        a.classList ? a.classList.remove(b) : a.className = a.className.replace(new RegExp("(^|\\b)" + b.split(" ").join("|") + "(\\b|$)", "gi"), " ")
    }

    function l(a, b) {
        return a.classList ? a.classList.contains(b) : new RegExp("\\b" + b + "\\b").test(a.className)
    }

    function m() {
        var a = void 0 !== window.pageXOffset,
            b = "CSS1Compat" === (document.compatMode || ""),
            c = a ? window.pageXOffset : b ? document.documentElement.scrollLeft : document.body.scrollLeft,
            d = a ? window.pageYOffset : b ? document.documentElement.scrollTop : document.body.scrollTop;
        return {
            x: c,
            y: d
        }
    }

    function n() {
        return window.navigator.pointerEnabled ? {
            start: "pointerdown",
            move: "pointermove",
            end: "pointerup"
        } : window.navigator.msPointerEnabled ? {
            start: "MSPointerDown",
            move: "MSPointerMove",
            end: "MSPointerUp"
        } : {
            start: "mousedown touchstart",
            move: "mousemove touchmove",
            end: "mouseup touchend"
        }
    }

    function o(a, b) {
        return 100 / (b - a)
    }

    function p(a, b) {
        return 100 * b / (a[1] - a[0])
    }

    function q(a, b) {
        return p(a, a[0] < 0 ? b + Math.abs(a[0]) : b - a[0])
    }

    function r(a, b) {
        return b * (a[1] - a[0]) / 100 + a[0]
    }

    function s(a, b) {
        for (var c = 1; a >= b[c];) c += 1;
        return c
    }

    function t(a, b, c) {
        if (c >= a.slice(-1)[0]) return 100;
        var d, e, f, g, h = s(c, a);
        return d = a[h - 1], e = a[h], f = b[h - 1], g = b[h], f + q([d, e], c) / o(f, g)
    }

    function u(a, b, c) {
        if (c >= 100) return a.slice(-1)[0];
        var d, e, f, g, h = s(c, b);
        return d = a[h - 1], e = a[h], f = b[h - 1], g = b[h], r([d, e], (c - f) * o(f, g))
    }

    function v(a, b, d, e) {
        if (100 === e) return e;
        var f, g, h = s(e, a);
        return d ? (f = a[h - 1], g = a[h], e - f > (g - f) / 2 ? g : f) : b[h - 1] ? a[h - 1] + c(e - a[h - 1], b[h - 1]) : e
    }

    function w(a, b, c) {
        var d;
        if ("number" == typeof b && (b = [b]), "[object Array]" !== Object.prototype.toString.call(b)) throw new Error("noUiSlider: 'range' contains invalid value.");
        if (d = "min" === a ? 0 : "max" === a ? 100 : parseFloat(a), !e(d) || !e(b[0])) throw new Error("noUiSlider: 'range' value isn't numeric.");
        c.xPct.push(d), c.xVal.push(b[0]), d ? c.xSteps.push(!isNaN(b[1]) && b[1]) : isNaN(b[1]) || (c.xSteps[0] = b[1]), c.xHighestCompleteStep.push(0)
    }

    function x(a, b, c) {
        if (!b) return !0;
        c.xSteps[a] = p([c.xVal[a], c.xVal[a + 1]], b) / o(c.xPct[a], c.xPct[a + 1]);
        var d = (c.xVal[a + 1] - c.xVal[a]) / c.xNumSteps[a],
            e = Math.ceil(Number(d.toFixed(3)) - 1),
            f = c.xVal[a] + c.xNumSteps[a] * e;
        c.xHighestCompleteStep[a] = f
    }

    function y(a, b, c, d) {
        this.xPct = [], this.xVal = [], this.xSteps = [d || !1], this.xNumSteps = [!1], this.xHighestCompleteStep = [], this.snap = b, this.direction = c;
        var e, f = [];
        for (e in a) a.hasOwnProperty(e) && f.push([a[e], e]);
        for (f.length && "object" == typeof f[0][0] ? f.sort(function(a, b) {
                return a[0][0] - b[0][0]
            }) : f.sort(function(a, b) {
                return a[0] - b[0]
            }), e = 0; e < f.length; e++) w(f[e][1], f[e][0], this);
        for (this.xNumSteps = this.xSteps.slice(0), e = 0; e < this.xNumSteps.length; e++) x(e, this.xNumSteps[e], this)
    }

    function z(a, b) {
        if (!e(b)) throw new Error("noUiSlider: 'step' is not numeric.");
        a.singleStep = b
    }

    function A(a, b) {
        if ("object" != typeof b || Array.isArray(b)) throw new Error("noUiSlider: 'range' is not an object.");
        if (void 0 === b.min || void 0 === b.max) throw new Error("noUiSlider: Missing 'min' or 'max' in 'range'.");
        if (b.min === b.max) throw new Error("noUiSlider: 'range' 'min' and 'max' cannot be equal.");
        a.spectrum = new y(b, a.snap, a.dir, a.singleStep)
    }

    function B(a, b) {
        if (b = h(b), !Array.isArray(b) || !b.length) throw new Error("noUiSlider: 'start' option is incorrect.");
        a.handles = b.length, a.start = b
    }

    function C(a, b) {
        if (a.snap = b, "boolean" != typeof b) throw new Error("noUiSlider: 'snap' option must be a boolean.")
    }

    function D(a, b) {
        if (a.animate = b, "boolean" != typeof b) throw new Error("noUiSlider: 'animate' option must be a boolean.")
    }

    function E(a, b) {
        if (a.animationDuration = b, "number" != typeof b) throw new Error("noUiSlider: 'animationDuration' option must be a number.")
    }

    function F(a, b) {
        var c, d = [!1];
        if ("lower" === b ? b = [!0, !1] : "upper" === b && (b = [!1, !0]), b === !0 || b === !1) {
            for (c = 1; c < a.handles; c++) d.push(b);
            d.push(!1)
        } else {
            if (!Array.isArray(b) || !b.length || b.length !== a.handles + 1) throw new Error("noUiSlider: 'connect' option doesn't match handle count.");
            d = b
        }
        a.connect = d
    }

    function G(a, b) {
        switch (b) {
            case "horizontal":
                a.ort = 0;
                break;
            case "vertical":
                a.ort = 1;
                break;
            default:
                throw new Error("noUiSlider: 'orientation' option is invalid.")
        }
    }

    function H(a, b) {
        if (!e(b)) throw new Error("noUiSlider: 'margin' option must be numeric.");
        if (0 !== b && (a.margin = a.spectrum.getMargin(b), !a.margin)) throw new Error("noUiSlider: 'margin' option is only supported on linear sliders.")
    }

    function I(a, b) {
        if (!e(b)) throw new Error("noUiSlider: 'limit' option must be numeric.");
        if (a.limit = a.spectrum.getMargin(b), !a.limit || a.handles < 2) throw new Error("noUiSlider: 'limit' option is only supported on linear sliders with 2 or more handles.")
    }

    function J(a, b) {
        if (!e(b)) throw new Error("noUiSlider: 'padding' option must be numeric.");
        if (0 !== b) {
            if (a.padding = a.spectrum.getMargin(b), !a.padding) throw new Error("noUiSlider: 'padding' option is only supported on linear sliders.");
            if (a.padding < 0) throw new Error("noUiSlider: 'padding' option must be a positive number.");
            if (a.padding >= 50) throw new Error("noUiSlider: 'padding' option must be less than half the range.")
        }
    }

    function K(a, b) {
        switch (b) {
            case "ltr":
                a.dir = 0;
                break;
            case "rtl":
                a.dir = 1;
                break;
            default:
                throw new Error("noUiSlider: 'direction' option was not recognized.")
        }
    }

    function L(a, b) {
        if ("string" != typeof b) throw new Error("noUiSlider: 'behaviour' must be a string containing options.");
        var c = b.indexOf("tap") >= 0,
            d = b.indexOf("drag") >= 0,
            e = b.indexOf("fixed") >= 0,
            f = b.indexOf("snap") >= 0,
            g = b.indexOf("hover") >= 0;
        if (e) {
            if (2 !== a.handles) throw new Error("noUiSlider: 'fixed' behaviour must be used with 2 handles");
            H(a, a.start[1] - a.start[0])
        }
        a.events = {
            tap: c || f,
            drag: d,
            fixed: e,
            snap: f,
            hover: g
        }
    }

    function M(a, b) {
        if (b !== !1)
            if (b === !0) {
                a.tooltips = [];
                for (var c = 0; c < a.handles; c++) a.tooltips.push(!0)
            } else {
                if (a.tooltips = h(b), a.tooltips.length !== a.handles) throw new Error("noUiSlider: must pass a formatter for all handles.");
                a.tooltips.forEach(function(a) {
                    if ("boolean" != typeof a && ("object" != typeof a || "function" != typeof a.to)) throw new Error("noUiSlider: 'tooltips' must be passed a formatter or 'false'.")
                })
            }
    }

    function N(a, b) {
        if (a.format = b, "function" == typeof b.to && "function" == typeof b.from) return !0;
        throw new Error("noUiSlider: 'format' requires 'to' and 'from' methods.")
    }

    function O(a, b) {
        if (void 0 !== b && "string" != typeof b && b !== !1) throw new Error("noUiSlider: 'cssPrefix' must be a string or `false`.");
        a.cssPrefix = b
    }

    function P(a, b) {
        if (void 0 !== b && "object" != typeof b) throw new Error("noUiSlider: 'cssClasses' must be an object.");
        if ("string" == typeof a.cssPrefix) {
            a.cssClasses = {};
            for (var c in b) b.hasOwnProperty(c) && (a.cssClasses[c] = a.cssPrefix + b[c])
        } else a.cssClasses = b
    }

    function Q(a, b) {
        if (b !== !0 && b !== !1) throw new Error("noUiSlider: 'useRequestAnimationFrame' option should be true (default) or false.");
        a.useRequestAnimationFrame = b
    }

    function R(a) {
        var b = {
                margin: 0,
                limit: 0,
                padding: 0,
                animate: !0,
                animationDuration: 300,
                format: U
            },
            c = {
                step: {
                    r: !1,
                    t: z
                },
                start: {
                    r: !0,
                    t: B
                },
                connect: {
                    r: !0,
                    t: F
                },
                direction: {
                    r: !0,
                    t: K
                },
                snap: {
                    r: !1,
                    t: C
                },
                animate: {
                    r: !1,
                    t: D
                },
                animationDuration: {
                    r: !1,
                    t: E
                },
                range: {
                    r: !0,
                    t: A
                },
                orientation: {
                    r: !1,
                    t: G
                },
                margin: {
                    r: !1,
                    t: H
                },
                limit: {
                    r: !1,
                    t: I
                },
                padding: {
                    r: !1,
                    t: J
                },
                behaviour: {
                    r: !0,
                    t: L
                },
                format: {
                    r: !1,
                    t: N
                },
                tooltips: {
                    r: !1,
                    t: M
                },
                cssPrefix: {
                    r: !1,
                    t: O
                },
                cssClasses: {
                    r: !1,
                    t: P
                },
                useRequestAnimationFrame: {
                    r: !1,
                    t: Q
                }
            },
            d = {
                connect: !1,
                direction: "ltr",
                behaviour: "tap",
                orientation: "horizontal",
                cssPrefix: "noUi-",
                cssClasses: {
                    target: "target",
                    base: "base",
                    origin: "origin",
                    handle: "handle",
                    handleLower: "handle-lower",
                    handleUpper: "handle-upper",
                    horizontal: "horizontal",
                    vertical: "vertical",
                    background: "background",
                    connect: "connect",
                    ltr: "ltr",
                    rtl: "rtl",
                    draggable: "draggable",
                    drag: "state-drag",
                    tap: "state-tap",
                    active: "active",
                    tooltip: "tooltip",
                    pips: "pips",
                    pipsHorizontal: "pips-horizontal",
                    pipsVertical: "pips-vertical",
                    marker: "marker",
                    markerHorizontal: "marker-horizontal",
                    markerVertical: "marker-vertical",
                    markerNormal: "marker-normal",
                    markerLarge: "marker-large",
                    markerSub: "marker-sub",
                    value: "value",
                    valueHorizontal: "value-horizontal",
                    valueVertical: "value-vertical",
                    valueNormal: "value-normal",
                    valueLarge: "value-large",
                    valueSub: "value-sub"
                },
                useRequestAnimationFrame: !0
            };
        Object.keys(c).forEach(function(e) {
            if (void 0 === a[e] && void 0 === d[e]) {
                if (c[e].r) throw new Error("noUiSlider: '" + e + "' is required.");
                return !0
            }
            c[e].t(b, void 0 === a[e] ? d[e] : a[e])
        }), b.pips = a.pips;
        var e = [
            ["left", "top"],
            ["right", "bottom"]
        ];
        return b.style = e[b.dir][b.ort], b.styleOposite = e[b.dir ? 0 : 1][b.ort], b
    }

    function S(c, e, i) {
        function o(b, c) {
            var d = a(b, e.cssClasses.origin),
                f = a(d, e.cssClasses.handle);
            return f.setAttribute("data-handle", c), 0 === c ? j(f, e.cssClasses.handleLower) : c === e.handles - 1 && j(f, e.cssClasses.handleUpper), d
        }

        function p(b, c) {
            return !!c && a(b, e.cssClasses.connect)
        }

        function q(a, b) {
            ba = [], ca = [], ca.push(p(b, a[0]));
            for (var c = 0; c < e.handles; c++) ba.push(o(b, c)), ha[c] = c, ca.push(p(b, a[c + 1]))
        }

        function r(b) {
            j(b, e.cssClasses.target), 0 === e.dir ? j(b, e.cssClasses.ltr) : j(b, e.cssClasses.rtl), 0 === e.ort ? j(b, e.cssClasses.horizontal) : j(b, e.cssClasses.vertical), aa = a(b, e.cssClasses.base)
        }

        function s(b, c) {
            return !!e.tooltips[c] && a(b.firstChild, e.cssClasses.tooltip)
        }

        function t() {
            var a = ba.map(s);
            Z("update", function(b, c, d) {
                if (a[c]) {
                    var f = b[c];
                    e.tooltips[c] !== !0 && (f = e.tooltips[c].to(d[c])), a[c].innerHTML = f
                }
            })
        }

        function u(a, b, c) {
            if ("range" === a || "steps" === a) return ja.xVal;
            if ("count" === a) {
                var d, e = 100 / (b - 1),
                    f = 0;
                for (b = [];
                    (d = f++ * e) <= 100;) b.push(d);
                a = "positions"
            }
            return "positions" === a ? b.map(function(a) {
                return ja.fromStepping(c ? ja.getStep(a) : a)
            }) : "values" === a ? c ? b.map(function(a) {
                return ja.fromStepping(ja.getStep(ja.toStepping(a)))
            }) : b : void 0
        }

        function v(a, c, d) {
            function e(a, b) {
                return (a + b).toFixed(7) / 1
            }
            var f = {},
                g = ja.xVal[0],
                h = ja.xVal[ja.xVal.length - 1],
                i = !1,
                j = !1,
                k = 0;
            return d = b(d.slice().sort(function(a, b) {
                return a - b
            })), d[0] !== g && (d.unshift(g), i = !0), d[d.length - 1] !== h && (d.push(h), j = !0), d.forEach(function(b, g) {
                var h, l, m, n, o, p, q, r, s, t, u = b,
                    v = d[g + 1];
                if ("steps" === c && (h = ja.xNumSteps[g]), h || (h = v - u), u !== !1 && void 0 !== v)
                    for (h = Math.max(h, 1e-7), l = u; l <= v; l = e(l, h)) {
                        for (n = ja.toStepping(l), o = n - k, r = o / a, s = Math.round(r), t = o / s, m = 1; m <= s; m += 1) p = k + m * t, f[p.toFixed(5)] = ["x", 0];
                        q = d.indexOf(l) > -1 ? 1 : "steps" === c ? 2 : 0, !g && i && (q = 0), l === v && j || (f[n.toFixed(5)] = [l, q]), k = n
                    }
            }), f
        }

        function w(a, b, c) {
            function d(a, b) {
                var c = b === e.cssClasses.value,
                    d = c ? m : n,
                    f = c ? k : l;
                return b + " " + d[e.ort] + " " + f[a]
            }

            function f(a, b, c) {
                return 'class="' + d(c[1], b) + '" style="' + e.style + ": " + a + '%"'
            }

            function g(a, d) {
                d[1] = d[1] && b ? b(d[0], d[1]) : d[1], i += "<div " + f(a, e.cssClasses.marker, d) + "></div>", d[1] && (i += "<div " + f(a, e.cssClasses.value, d) + ">" + c.to(d[0]) + "</div>")
            }
            var h = document.createElement("div"),
                i = "",
                k = [e.cssClasses.valueNormal, e.cssClasses.valueLarge, e.cssClasses.valueSub],
                l = [e.cssClasses.markerNormal, e.cssClasses.markerLarge, e.cssClasses.markerSub],
                m = [e.cssClasses.valueHorizontal, e.cssClasses.valueVertical],
                n = [e.cssClasses.markerHorizontal, e.cssClasses.markerVertical];
            return j(h, e.cssClasses.pips), j(h, 0 === e.ort ? e.cssClasses.pipsHorizontal : e.cssClasses.pipsVertical), Object.keys(a).forEach(function(b) {
                g(b, a[b])
            }), h.innerHTML = i, h
        }

        function x(a) {
            var b = a.mode,
                c = a.density || 1,
                d = a.filter || !1,
                e = a.values || !1,
                f = a.stepped || !1,
                g = u(b, e, f),
                h = v(c, b, g),
                i = a.format || {
                    to: Math.round
                };
            return fa.appendChild(w(h, d, i))
        }

        function y() {
            var a = aa.getBoundingClientRect(),
                b = "offset" + ["Width", "Height"][e.ort];
            return 0 === e.ort ? a.width || aa[b] : a.height || aa[b]
        }

        function z(a, b, c, d) {
            var f = function(b) {
                    return !fa.hasAttribute("disabled") && (!l(fa, e.cssClasses.tap) && (!!(b = A(b, d.pageOffset)) && (!(a === ea.start && void 0 !== b.buttons && b.buttons > 1) && ((!d.hover || !b.buttons) && (b.calcPoint = b.points[e.ort], void c(b, d))))))
                },
                g = [];
            return a.split(" ").forEach(function(a) {
                b.addEventListener(a, f, !1), g.push([a, f])
            }), g
        }

        function A(a, b) {
            a.preventDefault();
            var c, d, e = 0 === a.type.indexOf("touch"),
                f = 0 === a.type.indexOf("mouse"),
                g = 0 === a.type.indexOf("pointer");
            if (0 === a.type.indexOf("MSPointer") && (g = !0), e) {
                if (a.touches.length > 1) return !1;
                c = a.changedTouches[0].pageX, d = a.changedTouches[0].pageY
            }
            return b = b || m(), (f || g) && (c = a.clientX + b.x, d = a.clientY + b.y), a.pageOffset = b, a.points = [c, d], a.cursor = f || g, a
        }

        function B(a) {
            var b = a - d(aa, e.ort),
                c = 100 * b / y();
            return e.dir ? 100 - c : c
        }

        function C(a) {
            var b = 100,
                c = !1;
            return ba.forEach(function(d, e) {
                if (!d.hasAttribute("disabled")) {
                    var f = Math.abs(ga[e] - a);
                    f < b && (c = e, b = f)
                }
            }), c
        }

        function D(a, b, c, d) {
            var e = c.slice(),
                f = [!a, a],
                g = [a, !a];
            d = d.slice(), a && d.reverse(), d.length > 1 ? d.forEach(function(a, c) {
                var d = M(e, a, e[a] + b, f[c], g[c]);
                d === !1 ? b = 0 : (b = d - e[a], e[a] = d)
            }) : f = g = [!0];
            var h = !1;
            d.forEach(function(a, d) {
                h = Q(a, c[a] + b, f[d], g[d]) || h
            }), h && d.forEach(function(a) {
                E("update", a), E("slide", a)
            })
        }

        function E(a, b, c) {
            Object.keys(la).forEach(function(d) {
                var f = d.split(".")[0];
                a === f && la[d].forEach(function(a) {
                    a.call(da, ka.map(e.format.to), b, ka.slice(), c || !1, ga.slice())
                })
            })
        }

        function F(a, b) {
            "mouseout" === a.type && "HTML" === a.target.nodeName && null === a.relatedTarget && H(a, b)
        }

        function G(a, b) {
            if (navigator.appVersion.indexOf("MSIE 9") === -1 && 0 === a.buttons && 0 !== b.buttonsProperty) return H(a, b);
            var c = (e.dir ? -1 : 1) * (a.calcPoint - b.startCalcPoint),
                d = 100 * c / b.baseSize;
            D(c > 0, d, b.locations, b.handleNumbers)
        }

        function H(a, b) {
            ia && (k(ia, e.cssClasses.active), ia = !1), a.cursor && (document.body.style.cursor = "", document.body.removeEventListener("selectstart", document.body.noUiListener)), document.documentElement.noUiListeners.forEach(function(a) {
                document.documentElement.removeEventListener(a[0], a[1])
            }), k(fa, e.cssClasses.drag), P(), b.handleNumbers.forEach(function(a) {
                E("set", a), E("change", a), E("end", a)
            })
        }

        function I(a, b) {
            if (1 === b.handleNumbers.length) {
                var c = ba[b.handleNumbers[0]];
                if (c.hasAttribute("disabled")) return !1;
                ia = c.children[0], j(ia, e.cssClasses.active)
            }
            a.preventDefault(), a.stopPropagation();
            var d = z(ea.move, document.documentElement, G, {
                    startCalcPoint: a.calcPoint,
                    baseSize: y(),
                    pageOffset: a.pageOffset,
                    handleNumbers: b.handleNumbers,
                    buttonsProperty: a.buttons,
                    locations: ga.slice()
                }),
                f = z(ea.end, document.documentElement, H, {
                    handleNumbers: b.handleNumbers
                }),
                g = z("mouseout", document.documentElement, F, {
                    handleNumbers: b.handleNumbers
                });
            if (document.documentElement.noUiListeners = d.concat(f, g), a.cursor) {
                document.body.style.cursor = getComputedStyle(a.target).cursor, ba.length > 1 && j(fa, e.cssClasses.drag);
                var h = function() {
                    return !1
                };
                document.body.noUiListener = h, document.body.addEventListener("selectstart", h, !1)
            }
            b.handleNumbers.forEach(function(a) {
                E("start", a)
            })
        }

        function J(a) {
            a.stopPropagation();
            var b = B(a.calcPoint),
                c = C(b);
            return c !== !1 && (e.events.snap || f(fa, e.cssClasses.tap, e.animationDuration), Q(c, b, !0, !0), P(), E("slide", c, !0), E("set", c, !0), E("change", c, !0), E("update", c, !0), void(e.events.snap && I(a, {
                handleNumbers: [c]
            })))
        }

        function K(a) {
            var b = B(a.calcPoint),
                c = ja.getStep(b),
                d = ja.fromStepping(c);
            Object.keys(la).forEach(function(a) {
                "hover" === a.split(".")[0] && la[a].forEach(function(a) {
                    a.call(da, d)
                })
            })
        }

        function L(a) {
            a.fixed || ba.forEach(function(a, b) {
                z(ea.start, a.children[0], I, {
                    handleNumbers: [b]
                })
            }), a.tap && z(ea.start, aa, J, {}), a.hover && z(ea.move, aa, K, {
                hover: !0
            }), a.drag && ca.forEach(function(b, c) {
                if (b !== !1 && 0 !== c && c !== ca.length - 1) {
                    var d = ba[c - 1],
                        f = ba[c],
                        g = [b];
                    j(b, e.cssClasses.draggable), a.fixed && (g.push(d.children[0]), g.push(f.children[0])), g.forEach(function(a) {
                        z(ea.start, a, I, {
                            handles: [d, f],
                            handleNumbers: [c - 1, c]
                        })
                    })
                }
            })
        }

        function M(a, b, c, d, f) {
            return ba.length > 1 && (d && b > 0 && (c = Math.max(c, a[b - 1] + e.margin)), f && b < ba.length - 1 && (c = Math.min(c, a[b + 1] - e.margin))), ba.length > 1 && e.limit && (d && b > 0 && (c = Math.min(c, a[b - 1] + e.limit)), f && b < ba.length - 1 && (c = Math.max(c, a[b + 1] - e.limit))), e.padding && (0 === b && (c = Math.max(c, e.padding)), b === ba.length - 1 && (c = Math.min(c, 100 - e.padding))), c = ja.getStep(c), c = g(c), c !== a[b] && c
        }

        function N(a) {
            return a + "%"
        }

        function O(a, b) {
            ga[a] = b, ka[a] = ja.fromStepping(b);
            var c = function() {
                ba[a].style[e.style] = N(b), S(a), S(a + 1)
            };
            window.requestAnimationFrame && e.useRequestAnimationFrame ? window.requestAnimationFrame(c) : c()
        }

        function P() {
            ha.forEach(function(a) {
                var b = ga[a] > 50 ? -1 : 1,
                    c = 3 + (ba.length + b * a);
                ba[a].childNodes[0].style.zIndex = c
            })
        }

        function Q(a, b, c, d) {
            return b = M(ga, a, b, c, d), b !== !1 && (O(a, b), !0)
        }

        function S(a) {
            if (ca[a]) {
                var b = 0,
                    c = 100;
                0 !== a && (b = ga[a - 1]), a !== ca.length - 1 && (c = ga[a]), ca[a].style[e.style] = N(b), ca[a].style[e.styleOposite] = N(100 - c)
            }
        }

        function T(a, b) {
            null !== a && a !== !1 && ("number" == typeof a && (a = String(a)), a = e.format.from(a), a === !1 || isNaN(a) || Q(b, ja.toStepping(a), !1, !1))
        }

        function U(a, b) {
            var c = h(a),
                d = void 0 === ga[0];
            b = void 0 === b || !!b, c.forEach(T), e.animate && !d && f(fa, e.cssClasses.tap, e.animationDuration), ha.forEach(function(a) {
                Q(a, ga[a], !0, !1)
            }), P(), ha.forEach(function(a) {
                E("update", a), null !== c[a] && b && E("set", a)
            })
        }

        function V(a) {
            U(e.start, a)
        }

        function W() {
            var a = ka.map(e.format.to);
            return 1 === a.length ? a[0] : a
        }

        function X() {
            for (var a in e.cssClasses) e.cssClasses.hasOwnProperty(a) && k(fa, e.cssClasses[a]);
            for (; fa.firstChild;) fa.removeChild(fa.firstChild);
            delete fa.noUiSlider
        }

        function Y() {
            return ga.map(function(a, b) {
                var c = ja.getNearbySteps(a),
                    d = ka[b],
                    e = c.thisStep.step,
                    f = null;
                e !== !1 && d + e > c.stepAfter.startValue && (e = c.stepAfter.startValue - d), f = d > c.thisStep.startValue ? c.thisStep.step : c.stepBefore.step !== !1 && d - c.stepBefore.highestStep, 100 === a ? e = null : 0 === a && (f = null);
                var g = ja.countStepDecimals();
                return null !== e && e !== !1 && (e = Number(e.toFixed(g))), null !== f && f !== !1 && (f = Number(f.toFixed(g))), [f, e]
            })
        }

        function Z(a, b) {
            la[a] = la[a] || [], la[a].push(b), "update" === a.split(".")[0] && ba.forEach(function(a, b) {
                E("update", b)
            })
        }

        function $(a) {
            var b = a && a.split(".")[0],
                c = b && a.substring(b.length);
            Object.keys(la).forEach(function(a) {
                var d = a.split(".")[0],
                    e = a.substring(d.length);
                b && b !== d || c && c !== e || delete la[a]
            })
        }

        function _(a, b) {
            var c = W(),
                d = ["margin", "limit", "padding", "range", "animate", "snap", "step", "format"];
            d.forEach(function(b) {
                void 0 !== a[b] && (i[b] = a[b])
            });
            var f = R(i);
            d.forEach(function(b) {
                void 0 !== a[b] && (e[b] = f[b])
            }), f.spectrum.direction = ja.direction, ja = f.spectrum, e.margin = f.margin, e.limit = f.limit, e.padding = f.padding, ga = [], U(a.start || c, b)
        }
        var aa, ba, ca, da, ea = n(),
            fa = c,
            ga = [],
            ha = [],
            ia = !1,
            ja = e.spectrum,
            ka = [],
            la = {};
        if (fa.noUiSlider) throw new Error("Slider was already initialized.");
        return r(fa), q(e.connect, aa), da = {
            destroy: X,
            steps: Y,
            on: Z,
            off: $,
            get: W,
            set: U,
            reset: V,
            __moveHandles: function(a, b, c) {
                D(a, b, ga, c)
            },
            options: i,
            updateOptions: _,
            target: fa,
            pips: x
        }, L(e.events), U(e.start), e.pips && x(e.pips), e.tooltips && t(), da
    }

    function T(a, b) {
        if (!a.nodeName) throw new Error("noUiSlider.create requires a single element.");
        var c = R(b, a),
            d = S(a, c, b);
        return a.noUiSlider = d, d
    }
    y.prototype.getMargin = function(a) {
        var b = this.xNumSteps[0];
        if (b && a / b % 1 !== 0) throw new Error("noUiSlider: 'limit', 'margin' and 'padding' must be divisible by step.");
        return 2 === this.xPct.length && p(this.xVal, a)
    }, y.prototype.toStepping = function(a) {
        return a = t(this.xVal, this.xPct, a)
    }, y.prototype.fromStepping = function(a) {
        return u(this.xVal, this.xPct, a)
    }, y.prototype.getStep = function(a) {
        return a = v(this.xPct, this.xSteps, this.snap, a)
    }, y.prototype.getNearbySteps = function(a) {
        var b = s(a, this.xPct);
        return {
            stepBefore: {
                startValue: this.xVal[b - 2],
                step: this.xNumSteps[b - 2],
                highestStep: this.xHighestCompleteStep[b - 2]
            },
            thisStep: {
                startValue: this.xVal[b - 1],
                step: this.xNumSteps[b - 1],
                highestStep: this.xHighestCompleteStep[b - 1]
            },
            stepAfter: {
                startValue: this.xVal[b - 0],
                step: this.xNumSteps[b - 0],
                highestStep: this.xHighestCompleteStep[b - 0]
            }
        }
    }, y.prototype.countStepDecimals = function() {
        var a = this.xNumSteps.map(i);
        return Math.max.apply(null, a)
    }, y.prototype.convert = function(a) {
        return this.getStep(this.toStepping(a))
    };
    var U = {
        to: function(a) {
            return void 0 !== a && a.toFixed(2)
        },
        from: Number
    };
    return {
        create: T
    }
});