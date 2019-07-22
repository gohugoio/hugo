// Copyright 2015 The Hugo Authors. All rights reserved.
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
//
// Contains an embedded version of livereload.js
//
// Copyright (c) 2010-2015 Andrey Tarantsov
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package livereload

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/gorilla/websocket"
)

// Prefix to signal to LiveReload that we need to navigate to another path.
const hugoNavigatePrefix = "__hugo_navigate"

var upgrader = &websocket.Upgrader{
	// Hugo may potentially spin up multiple HTTP servers, so we need to exclude the
	// port when checking the origin.
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header["Origin"]
		if len(origin) == 0 {
			return true
		}
		u, err := url.Parse(origin[0])
		if err != nil {
			return false
		}

		if u.Host == r.Host {
			return true
		}

		h1, _, err := net.SplitHostPort(u.Host)
		if err != nil {
			return false
		}
		h2, _, err := net.SplitHostPort(r.Host)
		if err != nil {
			return false
		}

		return h1 == h2
	},
	ReadBufferSize: 1024, WriteBufferSize: 1024}

// Handler is a HandlerFunc handling the livereload
// Websocket interaction.
func Handler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	wsHub.register <- c
	defer func() { wsHub.unregister <- c }()
	go c.writer()
	c.reader()
}

// Initialize starts the Websocket Hub handling live reloads.
func Initialize() {
	go wsHub.run()
}

// ForceRefresh tells livereload to force a hard refresh.
func ForceRefresh() {
	RefreshPath("/x.js")
}

// NavigateToPath tells livereload to navigate to the given path.
// This translates to `window.location.href = path` in the client.
func NavigateToPath(path string) {
	RefreshPath(hugoNavigatePrefix + path)
}

// NavigateToPathForPort is similar to NavigateToPath but will also
// set window.location.port to the given port value.
func NavigateToPathForPort(path string, port int) {
	refreshPathForPort(hugoNavigatePrefix+path, port)
}

// RefreshPath tells livereload to refresh only the given path.
// If that path points to a CSS stylesheet or an image, only the changes
// will be updated in the browser, not the entire page.
func RefreshPath(s string) {
	refreshPathForPort(s, -1)
}

func refreshPathForPort(s string, port int) {
	// Tell livereload a file has changed - will force a hard refresh if not CSS or an image
	urlPath := filepath.ToSlash(s)
	portStr := ""
	if port > 0 {
		portStr = fmt.Sprintf(`, "overrideURL": %d`, port)
	}
	msg := fmt.Sprintf(`{"command":"reload","path":%q,"originalPath":"","liveCSS":true,"liveImg":true%s}`, urlPath, portStr)
	wsHub.broadcast <- []byte(msg)
}

// ServeJS serves the liverreload.js who's reference is injected into the page.
func ServeJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	w.Write(liveReloadJS())
}

func liveReloadJS() []byte {
	return []byte(livereloadJS + hugoLiveReloadPlugin)
}

var (
	// This is a patched version, see https://github.com/livereload/livereload-js/pull/84
	livereloadJS         = `!function(){return function e(t,o,n){function r(s,c){if(!o[s]){if(!t[s]){var a="function"==typeof require&&require;if(!c&&a)return a(s,!0);if(i)return i(s,!0);var l=new Error("Cannot find module '"+s+"'");throw l.code="MODULE_NOT_FOUND",l}var h=o[s]={exports:{}};t[s][0].call(h.exports,function(e){return r(t[s][1][e]||e)},h,h.exports,e,t,o,n)}return o[s].exports}for(var i="function"==typeof require&&require,s=0;s<n.length;s++)r(n[s]);return r}}()({1:[function(e,t,o){t.exports=function(e){if("function"!=typeof e)throw TypeError(e+" is not a function!");return e}},{}],2:[function(e,t,o){var n=e("./_wks")("unscopables"),r=Array.prototype;null==r[n]&&e("./_hide")(r,n,{}),t.exports=function(e){r[n][e]=!0}},{"./_hide":17,"./_wks":45}],3:[function(e,t,o){var n=e("./_is-object");t.exports=function(e){if(!n(e))throw TypeError(e+" is not an object!");return e}},{"./_is-object":21}],4:[function(e,t,o){var n=e("./_to-iobject"),r=e("./_to-length"),i=e("./_to-absolute-index");t.exports=function(e){return function(t,o,s){var c,a=n(t),l=r(a.length),h=i(s,l);if(e&&o!=o){for(;l>h;)if((c=a[h++])!=c)return!0}else for(;l>h;h++)if((e||h in a)&&a[h]===o)return e||h||0;return!e&&-1}}},{"./_to-absolute-index":38,"./_to-iobject":40,"./_to-length":41}],5:[function(e,t,o){var n={}.toString;t.exports=function(e){return n.call(e).slice(8,-1)}},{}],6:[function(e,t,o){var n=t.exports={version:"2.6.5"};"number"==typeof __e&&(__e=n)},{}],7:[function(e,t,o){var n=e("./_a-function");t.exports=function(e,t,o){if(n(e),void 0===t)return e;switch(o){case 1:return function(o){return e.call(t,o)};case 2:return function(o,n){return e.call(t,o,n)};case 3:return function(o,n,r){return e.call(t,o,n,r)}}return function(){return e.apply(t,arguments)}}},{"./_a-function":1}],8:[function(e,t,o){t.exports=function(e){if(null==e)throw TypeError("Can't call method on  "+e);return e}},{}],9:[function(e,t,o){t.exports=!e("./_fails")(function(){return 7!=Object.defineProperty({},"a",{get:function(){return 7}}).a})},{"./_fails":13}],10:[function(e,t,o){var n=e("./_is-object"),r=e("./_global").document,i=n(r)&&n(r.createElement);t.exports=function(e){return i?r.createElement(e):{}}},{"./_global":15,"./_is-object":21}],11:[function(e,t,o){t.exports="constructor,hasOwnProperty,isPrototypeOf,propertyIsEnumerable,toLocaleString,toString,valueOf".split(",")},{}],12:[function(e,t,o){var n=e("./_global"),r=e("./_core"),i=e("./_hide"),s=e("./_redefine"),c=e("./_ctx"),a=function(e,t,o){var l,h,u,d,f=e&a.F,p=e&a.G,_=e&a.S,m=e&a.P,g=e&a.B,y=p?n:_?n[t]||(n[t]={}):(n[t]||{}).prototype,v=p?r:r[t]||(r[t]={}),w=v.prototype||(v.prototype={});for(l in p&&(o=t),o)u=((h=!f&&y&&void 0!==y[l])?y:o)[l],d=g&&h?c(u,n):m&&"function"==typeof u?c(Function.call,u):u,y&&s(y,l,u,e&a.U),v[l]!=u&&i(v,l,d),m&&w[l]!=u&&(w[l]=u)};n.core=r,a.F=1,a.G=2,a.S=4,a.P=8,a.B=16,a.W=32,a.U=64,a.R=128,t.exports=a},{"./_core":6,"./_ctx":7,"./_global":15,"./_hide":17,"./_redefine":34}],13:[function(e,t,o){t.exports=function(e){try{return!!e()}catch(e){return!0}}},{}],14:[function(e,t,o){t.exports=e("./_shared")("native-function-to-string",Function.toString)},{"./_shared":37}],15:[function(e,t,o){var n=t.exports="undefined"!=typeof window&&window.Math==Math?window:"undefined"!=typeof self&&self.Math==Math?self:Function("return this")();"number"==typeof __g&&(__g=n)},{}],16:[function(e,t,o){var n={}.hasOwnProperty;t.exports=function(e,t){return n.call(e,t)}},{}],17:[function(e,t,o){var n=e("./_object-dp"),r=e("./_property-desc");t.exports=e("./_descriptors")?function(e,t,o){return n.f(e,t,r(1,o))}:function(e,t,o){return e[t]=o,e}},{"./_descriptors":9,"./_object-dp":28,"./_property-desc":33}],18:[function(e,t,o){var n=e("./_global").document;t.exports=n&&n.documentElement},{"./_global":15}],19:[function(e,t,o){t.exports=!e("./_descriptors")&&!e("./_fails")(function(){return 7!=Object.defineProperty(e("./_dom-create")("div"),"a",{get:function(){return 7}}).a})},{"./_descriptors":9,"./_dom-create":10,"./_fails":13}],20:[function(e,t,o){var n=e("./_cof");t.exports=Object("z").propertyIsEnumerable(0)?Object:function(e){return"String"==n(e)?e.split(""):Object(e)}},{"./_cof":5}],21:[function(e,t,o){t.exports=function(e){return"object"==typeof e?null!==e:"function"==typeof e}},{}],22:[function(e,t,o){"use strict";var n=e("./_object-create"),r=e("./_property-desc"),i=e("./_set-to-string-tag"),s={};e("./_hide")(s,e("./_wks")("iterator"),function(){return this}),t.exports=function(e,t,o){e.prototype=n(s,{next:r(1,o)}),i(e,t+" Iterator")}},{"./_hide":17,"./_object-create":27,"./_property-desc":33,"./_set-to-string-tag":35,"./_wks":45}],23:[function(e,t,o){"use strict";var n=e("./_library"),r=e("./_export"),i=e("./_redefine"),s=e("./_hide"),c=e("./_iterators"),a=e("./_iter-create"),l=e("./_set-to-string-tag"),h=e("./_object-gpo"),u=e("./_wks")("iterator"),d=!([].keys&&"next"in[].keys()),f=function(){return this};t.exports=function(e,t,o,p,_,m,g){a(o,t,p);var y,v,w,b=function(e){if(!d&&e in L)return L[e];switch(e){case"keys":case"values":return function(){return new o(this,e)}}return function(){return new o(this,e)}},S=t+" Iterator",R="values"==_,k=!1,L=e.prototype,x=L[u]||L["@@iterator"]||_&&L[_],j=x||b(_),C=_?R?b("entries"):j:void 0,O="Array"==t&&L.entries||x;if(O&&(w=h(O.call(new e)))!==Object.prototype&&w.next&&(l(w,S,!0),n||"function"==typeof w[u]||s(w,u,f)),R&&x&&"values"!==x.name&&(k=!0,j=function(){return x.call(this)}),n&&!g||!d&&!k&&L[u]||s(L,u,j),c[t]=j,c[S]=f,_)if(y={values:R?j:b("values"),keys:m?j:b("keys"),entries:C},g)for(v in y)v in L||i(L,v,y[v]);else r(r.P+r.F*(d||k),t,y);return y}},{"./_export":12,"./_hide":17,"./_iter-create":22,"./_iterators":25,"./_library":26,"./_object-gpo":30,"./_redefine":34,"./_set-to-string-tag":35,"./_wks":45}],24:[function(e,t,o){t.exports=function(e,t){return{value:t,done:!!e}}},{}],25:[function(e,t,o){t.exports={}},{}],26:[function(e,t,o){t.exports=!1},{}],27:[function(e,t,o){var n=e("./_an-object"),r=e("./_object-dps"),i=e("./_enum-bug-keys"),s=e("./_shared-key")("IE_PROTO"),c=function(){},a=function(){var t,o=e("./_dom-create")("iframe"),n=i.length;for(o.style.display="none",e("./_html").appendChild(o),o.src="javascript:",(t=o.contentWindow.document).open(),t.write("<script>document.F=Object<\/script>"),t.close(),a=t.F;n--;)delete a.prototype[i[n]];return a()};t.exports=Object.create||function(e,t){var o;return null!==e?(c.prototype=n(e),o=new c,c.prototype=null,o[s]=e):o=a(),void 0===t?o:r(o,t)}},{"./_an-object":3,"./_dom-create":10,"./_enum-bug-keys":11,"./_html":18,"./_object-dps":29,"./_shared-key":36}],28:[function(e,t,o){var n=e("./_an-object"),r=e("./_ie8-dom-define"),i=e("./_to-primitive"),s=Object.defineProperty;o.f=e("./_descriptors")?Object.defineProperty:function(e,t,o){if(n(e),t=i(t,!0),n(o),r)try{return s(e,t,o)}catch(e){}if("get"in o||"set"in o)throw TypeError("Accessors not supported!");return"value"in o&&(e[t]=o.value),e}},{"./_an-object":3,"./_descriptors":9,"./_ie8-dom-define":19,"./_to-primitive":43}],29:[function(e,t,o){var n=e("./_object-dp"),r=e("./_an-object"),i=e("./_object-keys");t.exports=e("./_descriptors")?Object.defineProperties:function(e,t){r(e);for(var o,s=i(t),c=s.length,a=0;c>a;)n.f(e,o=s[a++],t[o]);return e}},{"./_an-object":3,"./_descriptors":9,"./_object-dp":28,"./_object-keys":32}],30:[function(e,t,o){var n=e("./_has"),r=e("./_to-object"),i=e("./_shared-key")("IE_PROTO"),s=Object.prototype;t.exports=Object.getPrototypeOf||function(e){return e=r(e),n(e,i)?e[i]:"function"==typeof e.constructor&&e instanceof e.constructor?e.constructor.prototype:e instanceof Object?s:null}},{"./_has":16,"./_shared-key":36,"./_to-object":42}],31:[function(e,t,o){var n=e("./_has"),r=e("./_to-iobject"),i=e("./_array-includes")(!1),s=e("./_shared-key")("IE_PROTO");t.exports=function(e,t){var o,c=r(e),a=0,l=[];for(o in c)o!=s&&n(c,o)&&l.push(o);for(;t.length>a;)n(c,o=t[a++])&&(~i(l,o)||l.push(o));return l}},{"./_array-includes":4,"./_has":16,"./_shared-key":36,"./_to-iobject":40}],32:[function(e,t,o){var n=e("./_object-keys-internal"),r=e("./_enum-bug-keys");t.exports=Object.keys||function(e){return n(e,r)}},{"./_enum-bug-keys":11,"./_object-keys-internal":31}],33:[function(e,t,o){t.exports=function(e,t){return{enumerable:!(1&e),configurable:!(2&e),writable:!(4&e),value:t}}},{}],34:[function(e,t,o){var n=e("./_global"),r=e("./_hide"),i=e("./_has"),s=e("./_uid")("src"),c=e("./_function-to-string"),a=(""+c).split("toString");e("./_core").inspectSource=function(e){return c.call(e)},(t.exports=function(e,t,o,c){var l="function"==typeof o;l&&(i(o,"name")||r(o,"name",t)),e[t]!==o&&(l&&(i(o,s)||r(o,s,e[t]?""+e[t]:a.join(String(t)))),e===n?e[t]=o:c?e[t]?e[t]=o:r(e,t,o):(delete e[t],r(e,t,o)))})(Function.prototype,"toString",function(){return"function"==typeof this&&this[s]||c.call(this)})},{"./_core":6,"./_function-to-string":14,"./_global":15,"./_has":16,"./_hide":17,"./_uid":44}],35:[function(e,t,o){var n=e("./_object-dp").f,r=e("./_has"),i=e("./_wks")("toStringTag");t.exports=function(e,t,o){e&&!r(e=o?e:e.prototype,i)&&n(e,i,{configurable:!0,value:t})}},{"./_has":16,"./_object-dp":28,"./_wks":45}],36:[function(e,t,o){var n=e("./_shared")("keys"),r=e("./_uid");t.exports=function(e){return n[e]||(n[e]=r(e))}},{"./_shared":37,"./_uid":44}],37:[function(e,t,o){var n=e("./_core"),r=e("./_global"),i=r["__core-js_shared__"]||(r["__core-js_shared__"]={});(t.exports=function(e,t){return i[e]||(i[e]=void 0!==t?t:{})})("versions",[]).push({version:n.version,mode:e("./_library")?"pure":"global",copyright:"© 2019 Denis Pushkarev (zloirock.ru)"})},{"./_core":6,"./_global":15,"./_library":26}],38:[function(e,t,o){var n=e("./_to-integer"),r=Math.max,i=Math.min;t.exports=function(e,t){return(e=n(e))<0?r(e+t,0):i(e,t)}},{"./_to-integer":39}],39:[function(e,t,o){var n=Math.ceil,r=Math.floor;t.exports=function(e){return isNaN(e=+e)?0:(e>0?r:n)(e)}},{}],40:[function(e,t,o){var n=e("./_iobject"),r=e("./_defined");t.exports=function(e){return n(r(e))}},{"./_defined":8,"./_iobject":20}],41:[function(e,t,o){var n=e("./_to-integer"),r=Math.min;t.exports=function(e){return e>0?r(n(e),9007199254740991):0}},{"./_to-integer":39}],42:[function(e,t,o){var n=e("./_defined");t.exports=function(e){return Object(n(e))}},{"./_defined":8}],43:[function(e,t,o){var n=e("./_is-object");t.exports=function(e,t){if(!n(e))return e;var o,r;if(t&&"function"==typeof(o=e.toString)&&!n(r=o.call(e)))return r;if("function"==typeof(o=e.valueOf)&&!n(r=o.call(e)))return r;if(!t&&"function"==typeof(o=e.toString)&&!n(r=o.call(e)))return r;throw TypeError("Can't convert object to primitive value")}},{"./_is-object":21}],44:[function(e,t,o){var n=0,r=Math.random();t.exports=function(e){return"Symbol(".concat(void 0===e?"":e,")_",(++n+r).toString(36))}},{}],45:[function(e,t,o){var n=e("./_shared")("wks"),r=e("./_uid"),i=e("./_global").Symbol,s="function"==typeof i;(t.exports=function(e){return n[e]||(n[e]=s&&i[e]||(s?i:r)("Symbol."+e))}).store=n},{"./_global":15,"./_shared":37,"./_uid":44}],46:[function(e,t,o){"use strict";var n=e("./_add-to-unscopables"),r=e("./_iter-step"),i=e("./_iterators"),s=e("./_to-iobject");t.exports=e("./_iter-define")(Array,"Array",function(e,t){this._t=s(e),this._i=0,this._k=t},function(){var e=this._t,t=this._k,o=this._i++;return!e||o>=e.length?(this._t=void 0,r(1)):r(0,"keys"==t?o:"values"==t?e[o]:[o,e[o]])},"values"),i.Arguments=i.Array,n("keys"),n("values"),n("entries")},{"./_add-to-unscopables":2,"./_iter-define":23,"./_iter-step":24,"./_iterators":25,"./_to-iobject":40}],47:[function(e,t,o){for(var n=e("./es6.array.iterator"),r=e("./_object-keys"),i=e("./_redefine"),s=e("./_global"),c=e("./_hide"),a=e("./_iterators"),l=e("./_wks"),h=l("iterator"),u=l("toStringTag"),d=a.Array,f={CSSRuleList:!0,CSSStyleDeclaration:!1,CSSValueList:!1,ClientRectList:!1,DOMRectList:!1,DOMStringList:!1,DOMTokenList:!0,DataTransferItemList:!1,FileList:!1,HTMLAllCollection:!1,HTMLCollection:!1,HTMLFormElement:!1,HTMLSelectElement:!1,MediaList:!0,MimeTypeArray:!1,NamedNodeMap:!1,NodeList:!0,PaintRequestList:!1,Plugin:!1,PluginArray:!1,SVGLengthList:!1,SVGNumberList:!1,SVGPathSegList:!1,SVGPointList:!1,SVGStringList:!1,SVGTransformList:!1,SourceBufferList:!1,StyleSheetList:!0,TextTrackCueList:!1,TextTrackList:!1,TouchList:!1},p=r(f),_=0;_<p.length;_++){var m,g=p[_],y=f[g],v=s[g],w=v&&v.prototype;if(w&&(w[h]||c(w,h,d),w[u]||c(w,u,g),a[g]=d,y))for(m in n)w[m]||i(w,m,n[m],!0)}},{"./_global":15,"./_hide":17,"./_iterators":25,"./_object-keys":32,"./_redefine":34,"./_wks":45,"./es6.array.iterator":46}],48:[function(e,t,o){"use strict";const{Parser:n,PROTOCOL_6:r,PROTOCOL_7:i}=e("./protocol"),s="3.0.0";o.Connector=class{constructor(e,t,o,r){this.options=e,this.WebSocket=t,this.Timer=o,this.handlers=r;const i=this.options.path?"".concat(this.options.path):"livereload";this._uri="ws".concat(this.options.https?"s":"","://").concat(this.options.host,":").concat(this.options.port,"/").concat(i),this._nextDelay=this.options.mindelay,this._connectionDesired=!1,this.protocol=0,this.protocolParser=new n({connected:e=>(this.protocol=e,this._handshakeTimeout.stop(),this._nextDelay=this.options.mindelay,this._disconnectionReason="broken",this.handlers.connected(this.protocol)),error:e=>(this.handlers.error(e),this._closeOnError()),message:e=>this.handlers.message(e)}),this._handshakeTimeout=new this.Timer(()=>{if(this._isSocketConnected())return this._disconnectionReason="handshake-timeout",this.socket.close()}),this._reconnectTimer=new this.Timer(()=>{if(this._connectionDesired)return this.connect()}),this.connect()}_isSocketConnected(){return this.socket&&this.socket.readyState===this.WebSocket.OPEN}connect(){this._connectionDesired=!0,this._isSocketConnected()||(this._reconnectTimer.stop(),this._disconnectionReason="cannot-connect",this.protocolParser.reset(),this.handlers.connecting(),this.socket=new this.WebSocket(this._uri),this.socket.onopen=(e=>this._onopen(e)),this.socket.onclose=(e=>this._onclose(e)),this.socket.onmessage=(e=>this._onmessage(e)),this.socket.onerror=(e=>this._onerror(e)))}disconnect(){if(this._connectionDesired=!1,this._reconnectTimer.stop(),this._isSocketConnected())return this._disconnectionReason="manual",this.socket.close()}_scheduleReconnection(){this._connectionDesired&&(this._reconnectTimer.running||(this._reconnectTimer.start(this._nextDelay),this._nextDelay=Math.min(this.options.maxdelay,2*this._nextDelay)))}sendCommand(e){if(this.protocol)return this._sendCommand(e)}_sendCommand(e){return this.socket.send(JSON.stringify(e))}_closeOnError(){return this._handshakeTimeout.stop(),this._disconnectionReason="error",this.socket.close()}_onopen(e){this.handlers.socketConnected(),this._disconnectionReason="handshake-failed";const t={command:"hello",protocols:[r,i]};return t.ver=s,this.options.ext&&(t.ext=this.options.ext),this.options.extver&&(t.extver=this.options.extver),this.options.snipver&&(t.snipver=this.options.snipver),this._sendCommand(t),this._handshakeTimeout.start(this.options.handshake_timeout)}_onclose(e){return this.protocol=0,this.handlers.disconnected(this._disconnectionReason,this._nextDelay),this._scheduleReconnection()}_onerror(e){}_onmessage(e){return this.protocolParser.process(e.data)}}},{"./protocol":53}],49:[function(e,t,o){"use strict";const n={bind(e,t,o){if(e.addEventListener)return e.addEventListener(t,o,!1);if(e.attachEvent)return e[t]=1,e.attachEvent("onpropertychange",function(e){if(e.propertyName===t)return o()});throw new Error("Attempt to attach custom event ".concat(t," to something which isn't a DOMElement"))},fire(e,t){if(e.addEventListener){const e=document.createEvent("HTMLEvents");return e.initEvent(t,!0,!0),document.dispatchEvent(e)}if(!e.attachEvent)throw new Error("Attempt to fire custom event ".concat(t," on something which isn't a DOMElement"));if(e[t])return e[t]++}};o.bind=n.bind,o.fire=n.fire},{}],50:[function(e,t,o){"use strict";class n{constructor(e,t){this.window=e,this.host=t}reload(e,t){if(this.window.less&&this.window.less.refresh){if(e.match(/\.less$/i))return this.reloadLess(e);if(t.originalPath.match(/\.less$/i))return this.reloadLess(t.originalPath)}return!1}reloadLess(e){let t;const o=(()=>{const e=[];for(t of Array.from(document.getElementsByTagName("link")))(t.href&&t.rel.match(/^stylesheet\/less$/i)||t.rel.match(/stylesheet/i)&&t.type.match(/^text\/(x-)?less$/i))&&e.push(t);return e})();if(0===o.length)return!1;for(t of Array.from(o))t.href=this.host.generateCacheBustUrl(t.href);return this.host.console.log("LiveReload is asking LESS to recompile all stylesheets"),this.window.less.refresh(!0),!0}analyze(){return{disable:!(!this.window.less||!this.window.less.refresh)}}}n.identifier="less",n.version="1.0",t.exports=n},{}],51:[function(e,t,o){"use strict";e("core-js/modules/web.dom.iterable");const{Connector:n}=e("./connector"),{Timer:r}=e("./timer"),{Options:i}=e("./options"),{Reloader:s}=e("./reloader"),{ProtocolError:c}=e("./protocol");o.LiveReload=class{constructor(e){if(this.window=e,this.listeners={},this.plugins=[],this.pluginIdentifiers={},this.console=this.window.console&&this.window.console.log&&this.window.console.error?this.window.location.href.match(/LR-verbose/)?this.window.console:{log(){},error:this.window.console.error.bind(this.window.console)}:{log(){},error(){}},this.WebSocket=this.window.WebSocket||this.window.MozWebSocket){if("LiveReloadOptions"in e){this.options=new i;for(let t of Object.keys(e.LiveReloadOptions||{})){const o=e.LiveReloadOptions[t];this.options.set(t,o)}}else if(this.options=i.extract(this.window.document),!this.options)return void this.console.error("LiveReload disabled because it could not find its own <SCRIPT> tag");this.reloader=new s(this.window,this.console,r),this.connector=new n(this.options,this.WebSocket,r,{connecting:()=>{},socketConnected:()=>{},connected:e=>("function"==typeof this.listeners.connect&&this.listeners.connect(),this.log("LiveReload is connected to ".concat(this.options.host,":").concat(this.options.port," (protocol v").concat(e,").")),this.analyze()),error:e=>{if(e instanceof c){if("undefined"!=typeof console&&null!==console)return console.log("".concat(e.message,"."))}else if("undefined"!=typeof console&&null!==console)return console.log("LiveReload internal error: ".concat(e.message))},disconnected:(e,t)=>{switch("function"==typeof this.listeners.disconnect&&this.listeners.disconnect(),e){case"cannot-connect":return this.log("LiveReload cannot connect to ".concat(this.options.host,":").concat(this.options.port,", will retry in ").concat(t," sec."));case"broken":return this.log("LiveReload disconnected from ".concat(this.options.host,":").concat(this.options.port,", reconnecting in ").concat(t," sec."));case"handshake-timeout":return this.log("LiveReload cannot connect to ".concat(this.options.host,":").concat(this.options.port," (handshake timeout), will retry in ").concat(t," sec."));case"handshake-failed":return this.log("LiveReload cannot connect to ".concat(this.options.host,":").concat(this.options.port," (handshake failed), will retry in ").concat(t," sec."));case"manual":case"error":default:return this.log("LiveReload disconnected from ".concat(this.options.host,":").concat(this.options.port," (").concat(e,"), reconnecting in ").concat(t," sec."))}},message:e=>{switch(e.command){case"reload":return this.performReload(e);case"alert":return this.performAlert(e)}}}),this.initialized=!0}else this.console.error("LiveReload disabled because the browser does not seem to support web sockets")}on(e,t){this.listeners[e]=t}log(e){return this.console.log("".concat(e))}performReload(e){return this.log("LiveReload received reload request: ".concat(JSON.stringify(e,null,2))),this.reloader.reload(e.path,{liveCSS:null==e.liveCSS||e.liveCSS,liveImg:null==e.liveImg||e.liveImg,reloadMissingCSS:null==e.reloadMissingCSS||e.reloadMissingCSS,originalPath:e.originalPath||"",overrideURL:e.overrideURL||"",serverURL:"http://".concat(this.options.host,":").concat(this.options.port)})}performAlert(e){return alert(e.message)}shutDown(){if(this.initialized)return this.connector.disconnect(),this.log("LiveReload disconnected."),"function"==typeof this.listeners.shutdown?this.listeners.shutdown():void 0}hasPlugin(e){return!!this.pluginIdentifiers[e]}addPlugin(e){if(!this.initialized)return;if(this.hasPlugin(e.identifier))return;this.pluginIdentifiers[e.identifier]=!0;const t=new e(this.window,{_livereload:this,_reloader:this.reloader,_connector:this.connector,console:this.console,Timer:r,generateCacheBustUrl:e=>this.reloader.generateCacheBustUrl(e)});this.plugins.push(t),this.reloader.addPlugin(t)}analyze(){if(!this.initialized)return;if(!(this.connector.protocol>=7))return;const e={};for(let o of this.plugins){var t=("function"==typeof o.analyze?o.analyze():void 0)||{};e[o.constructor.identifier]=t,t.version=o.constructor.version}this.connector.sendCommand({command:"info",plugins:e,url:this.window.location.href})}}},{"./connector":48,"./options":52,"./protocol":53,"./reloader":54,"./timer":56,"core-js/modules/web.dom.iterable":47}],52:[function(e,t,o){"use strict";class n{constructor(){this.https=!1,this.host=null,this.port=35729,this.snipver=null,this.ext=null,this.extver=null,this.mindelay=1e3,this.maxdelay=6e4,this.handshake_timeout=5e3}set(e,t){void 0!==t&&(isNaN(+t)||(t=+t),this[e]=t)}}n.extract=function(e){for(let s of Array.from(e.getElementsByTagName("script"))){var t,o;if((o=s.src)&&(t=o.match(new RegExp("^[^:]+://(.*)/z?livereload\\.js(?:\\?(.*))?$")))){var r;const e=new n;if(e.https=0===o.indexOf("https"),(r=t[1].match(new RegExp("^([^/:]+)(?::(\\d+))?(\\/+.*)?$")))&&(e.host=r[1],r[2]&&(e.port=parseInt(r[2],10))),t[2])for(let o of t[2].split("&")){var i;(i=o.split("=")).length>1&&e.set(i[0].replace(/-/g,"_"),i.slice(1).join("="))}return e}}return null},o.Options=n},{}],53:[function(e,t,o){"use strict";let n,r;o.PROTOCOL_6=n="http://livereload.com/protocols/official-6",o.PROTOCOL_7=r="http://livereload.com/protocols/official-7";class i{constructor(e,t){this.message="LiveReload protocol error (".concat(e,') after receiving data: "').concat(t,'".')}}o.ProtocolError=i,o.Parser=class{constructor(e){this.handlers=e,this.reset()}reset(){this.protocol=null}process(e){try{let t;if(this.protocol){if(6===this.protocol){if(!(t=JSON.parse(e)).length)throw new i("protocol 6 messages must be arrays");const[o,n]=Array.from(t);if("refresh"!==o)throw new i("unknown protocol 6 command");return this.handlers.message({command:"reload",path:n.path,liveCSS:null==n.apply_css_live||n.apply_css_live})}return t=this._parseMessage(e,["reload","alert"]),this.handlers.message(t)}if(e.match(new RegExp("^!!ver:([\\d.]+)$")))this.protocol=6;else if(t=this._parseMessage(e,["hello"])){if(!t.protocols.length)throw new i("no protocols specified in handshake message");if(Array.from(t.protocols).includes(r))this.protocol=7;else{if(!Array.from(t.protocols).includes(n))throw new i("no supported protocols found");this.protocol=6}}return this.handlers.connected(this.protocol)}catch(e){if(e instanceof i)return this.handlers.error(e);throw e}}_parseMessage(e,t){let o;try{o=JSON.parse(e)}catch(t){throw new i("unparsable JSON",e)}if(!o.command)throw new i('missing "command" key',e);if(!t.includes(o.command))throw new i("invalid command '".concat(o.command,"', only valid commands are: ").concat(t.join(", "),")"),e);return o}}},{}],54:[function(e,t,o){"use strict";const n=function(e){let t,o,n;(o=e.indexOf("#"))>=0?(t=e.slice(o),e=e.slice(0,o)):t="";const r=e.indexOf("??");return r>=0?r+1!==e.lastIndexOf("?")&&(o=e.lastIndexOf("?")):o=e.indexOf("?"),o>=0?(n=e.slice(o),e=e.slice(0,o)):n="",{url:e,params:n,hash:t}},r=function(e){if(!e)return"";let t;return({url:e}=n(e)),t=0===e.indexOf("file://")?e.replace(new RegExp("^file://(localhost)?"),""):e.replace(new RegExp("^([^:]+:)?//([^:/]+)(:\\d*)?/"),"/"),decodeURIComponent(t)},i=function(e,t,o){let n,r={score:0};for(let i of t)(n=s(e,o(i)))>r.score&&(r={object:i,score:n});return 0===r.score?null:r};var s=function(e,t){if((e=e.replace(/^\/+/,"").toLowerCase())===(t=t.replace(/^\/+/,"").toLowerCase()))return 1e4;const o=e.split("/").reverse(),n=t.split("/").reverse(),r=Math.min(o.length,n.length);let i=0;for(;i<r&&o[i]===n[i];)++i;return i};const c=(e,t)=>s(e,t)>0,a=[{selector:"background",styleNames:["backgroundImage"]},{selector:"border",styleNames:["borderImage","webkitBorderImage","MozBorderImage"]}];o.Reloader=class{constructor(e,t,o){this.window=e,this.console=t,this.Timer=o,this.document=this.window.document,this.importCacheWaitPeriod=200,this.plugins=[]}addPlugin(e){return this.plugins.push(e)}analyze(e){}reload(e,t){this.options=t,this.options.stylesheetReloadTimeout||(this.options.stylesheetReloadTimeout=15e3);for(let o of Array.from(this.plugins))if(o.reload&&o.reload(e,t))return;if(!(t.liveCSS&&e.match(/\.css(?:\.map)?$/i)&&this.reloadStylesheet(e)))if(t.liveImg&&e.match(/\.(jpe?g|png|gif)$/i))this.reloadImages(e);else{if(!t.isChromeExtension)return this.reloadPage();this.reloadChromeExtension()}}reloadPage(){return this.window.document.location.reload()}reloadChromeExtension(){return this.window.chrome.runtime.reload()}reloadImages(e){let t;const o=this.generateUniqueString();for(t of Array.from(this.document.images))c(e,r(t.src))&&(t.src=this.generateCacheBustUrl(t.src,o));if(this.document.querySelectorAll)for(let{selector:n,styleNames:r}of a)for(t of Array.from(this.document.querySelectorAll("[style*=".concat(n,"]"))))this.reloadStyleImages(t.style,r,e,o);if(this.document.styleSheets)return Array.from(this.document.styleSheets).map(t=>this.reloadStylesheetImages(t,e,o))}reloadStylesheetImages(e,t,o){let n;try{n=(e||{}).cssRules}catch(e){}if(n)for(let e of Array.from(n))switch(e.type){case CSSRule.IMPORT_RULE:this.reloadStylesheetImages(e.styleSheet,t,o);break;case CSSRule.STYLE_RULE:for(let{styleNames:n}of a)this.reloadStyleImages(e.style,n,t,o);break;case CSSRule.MEDIA_RULE:this.reloadStylesheetImages(e,t,o)}}reloadStyleImages(e,t,o,n){for(let i of t){const t=e[i];if("string"==typeof t){const s=t.replace(new RegExp("\\burl\\s*\\(([^)]*)\\)"),(e,t)=>c(o,r(t))?"url(".concat(this.generateCacheBustUrl(t,n),")"):e);s!==t&&(e[i]=s)}}}reloadStylesheet(e){let t,o;const n=(()=>{const e=[];for(o of Array.from(this.document.getElementsByTagName("link")))o.rel.match(/^stylesheet$/i)&&!o.__LiveReload_pendingRemoval&&e.push(o);return e})(),s=[];for(t of Array.from(this.document.getElementsByTagName("style")))t.sheet&&this.collectImportedStylesheets(t,t.sheet,s);for(o of Array.from(n))this.collectImportedStylesheets(o,o.sheet,s);if(this.window.StyleFix&&this.document.querySelectorAll)for(t of Array.from(this.document.querySelectorAll("style[data-href]")))n.push(t);this.console.log("LiveReload found ".concat(n.length," LINKed stylesheets, ").concat(s.length," @imported stylesheets"));const c=i(e,n.concat(s),e=>r(this.linkHref(e)));if(c)c.object.rule?(this.console.log("LiveReload is reloading imported stylesheet: ".concat(c.object.href)),this.reattachImportedRule(c.object)):(this.console.log("LiveReload is reloading stylesheet: ".concat(this.linkHref(c.object))),this.reattachStylesheetLink(c.object));else if(this.options.reloadMissingCSS)for(o of(this.console.log("LiveReload will reload all stylesheets because path '".concat(e,"' did not match any specific one. To disable this behavior, set 'options.reloadMissingCSS' to 'false'.")),Array.from(n)))this.reattachStylesheetLink(o);else this.console.log("LiveReload will not reload path '".concat(e,"' because the stylesheet was not found on the page and 'options.reloadMissingCSS' was set to 'false'."));return!0}collectImportedStylesheets(e,t,o){let n;try{n=(t||{}).cssRules}catch(e){}if(n&&n.length)for(let t=0;t<n.length;t++){const r=n[t];switch(r.type){case CSSRule.CHARSET_RULE:continue;case CSSRule.IMPORT_RULE:o.push({link:e,rule:r,index:t,href:r.href}),this.collectImportedStylesheets(e,r.styleSheet,o)}}}waitUntilCssLoads(e,t){let o=!1;const n=()=>{if(!o)return o=!0,t()};if(e.onload=(()=>(this.console.log("LiveReload: the new stylesheet has finished loading"),this.knownToSupportCssOnLoad=!0,n())),!this.knownToSupportCssOnLoad){let t;(t=(()=>e.sheet?(this.console.log("LiveReload is polling until the new CSS finishes loading..."),n()):this.Timer.start(50,t)))()}return this.Timer.start(this.options.stylesheetReloadTimeout,n)}linkHref(e){return e.href||e.getAttribute&&e.getAttribute("data-href")}reattachStylesheetLink(e){let t;if(e.__LiveReload_pendingRemoval)return;e.__LiveReload_pendingRemoval=!0,"STYLE"===e.tagName?((t=this.document.createElement("link")).rel="stylesheet",t.media=e.media,t.disabled=e.disabled):t=e.cloneNode(!1),t.href=this.generateCacheBustUrl(this.linkHref(e));const o=e.parentNode;return o.lastChild===e?o.appendChild(t):o.insertBefore(t,e.nextSibling),this.waitUntilCssLoads(t,()=>{let o;return o=/AppleWebKit/.test(navigator.userAgent)?5:200,this.Timer.start(o,()=>{if(e.parentNode)return e.parentNode.removeChild(e),t.onreadystatechange=null,this.window.StyleFix?this.window.StyleFix.link(t):void 0})})}reattachImportedRule({rule:e,index:t,link:o}){const n=e.parentStyleSheet,r=this.generateCacheBustUrl(e.href),i=e.media.length?[].join.call(e.media,", "):"",s='@import url("'.concat(r,'") ').concat(i,";");e.__LiveReload_newHref=r;const c=this.document.createElement("link");return c.rel="stylesheet",c.href=r,c.__LiveReload_pendingRemoval=!0,o.parentNode&&o.parentNode.insertBefore(c,o),this.Timer.start(this.importCacheWaitPeriod,()=>{if(c.parentNode&&c.parentNode.removeChild(c),e.__LiveReload_newHref===r)return n.insertRule(s,t),n.deleteRule(t+1),(e=n.cssRules[t]).__LiveReload_newHref=r,this.Timer.start(this.importCacheWaitPeriod,()=>{if(e.__LiveReload_newHref===r)return n.insertRule(s,t),n.deleteRule(t+1)})})}generateUniqueString(){return"livereload=".concat(Date.now())}generateCacheBustUrl(e,t){let o,r;if(t||(t=this.generateUniqueString()),({url:e,hash:o,params:r}=n(e)),this.options.overrideURL&&e.indexOf(this.options.serverURL)<0){const t=e;e=this.options.serverURL+this.options.overrideURL+"?url="+encodeURIComponent(e),this.console.log("LiveReload is overriding source URL ".concat(t," with ").concat(e))}let i=r.replace(/(\?|&)livereload=(\d+)/,(e,o)=>"".concat(o).concat(t));return i===r&&(i=0===r.length?"?".concat(t):"".concat(r,"&").concat(t)),e+i+o}}},{}],55:[function(e,t,o){"use strict";const n=e("./customevents"),r=window.LiveReload=new(e("./livereload").LiveReload)(window);for(let e in window)e.match(/^LiveReloadPlugin/)&&r.addPlugin(window[e]);r.addPlugin(e("./less")),r.on("shutdown",()=>delete window.LiveReload),r.on("connect",()=>n.fire(document,"LiveReloadConnect")),r.on("disconnect",()=>n.fire(document,"LiveReloadDisconnect")),n.bind(document,"LiveReloadShutDown",()=>r.shutDown())},{"./customevents":49,"./less":50,"./livereload":51}],56:[function(e,t,o){"use strict";class n{constructor(e){this.func=e,this.running=!1,this.id=null,this._handler=(()=>(this.running=!1,this.id=null,this.func()))}start(e){this.running&&clearTimeout(this.id),this.id=setTimeout(this._handler,e),this.running=!0}stop(){this.running&&(clearTimeout(this.id),this.running=!1,this.id=null)}}n.start=((e,t)=>setTimeout(t,e)),o.Timer=n},{}]},{},[55]);`
	hugoLiveReloadPlugin = fmt.Sprintf(`
/*
Hugo adds a specific prefix, "__hugo_navigate", to the path in certain situations to signal
navigation to another content page.
*/

function HugoReload() {}

HugoReload.identifier = 'hugoReloader';
HugoReload.version = '0.9';

HugoReload.prototype.reload = function(path, options) {
	var prefix = %q;

	if (path.lastIndexOf(prefix, 0) !== 0) {
		return false
	}
	
	path = path.substring(prefix.length);

	var portChanged = options.overrideURL && options.overrideURL != window.location.port
	
	if (!portChanged && window.location.pathname === path) {
		window.location.reload();
	} else {
		if (portChanged) {
			window.location = location.protocol + "//" + location.hostname + ":" + options.overrideURL + path;
		} else {
			window.location.pathname = path;
		}
	}

	return true;
};

LiveReload.addPlugin(HugoReload)
`, hugoNavigatePrefix)
)
