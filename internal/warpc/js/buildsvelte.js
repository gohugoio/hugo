import { readInput, writeOutput } from './common';
import { compile } from 'svelte/compiler';

const build = function (input) {
	const data = input.data; //
	const source = data.source;
	const opts = data.options || {}; // //
	const header = input.header;

	try {
		let { js, css, warnings } = compile(source);
		writeOutput({ header: header, data: { result: js.code } });
	} catch (e) {
		header.err = e.message;
		writeOutput({ header: header });
	}
};

readInput(build);
