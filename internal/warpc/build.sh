# TODO1 clean up when done.
go generate ./gen
javy compile js/greet.bundle.js -d -o wasm/greet.wasm
javy compile js/renderkatex.bundle.js -d -o wasm/renderkatex.wasm
touch warpc_test.go