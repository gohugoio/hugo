import { readInput, writeOutput } from './common';
import katex from 'katex';

const render = function (input) {
	const data = input.data;
	const expression = data.expression;
	const options = data.options;
	const header = input.header;
	try {
		const output = katex.renderToString(expression, options);
		writeOutput({ header: header, data: { output: output } });
	} catch (e) {
		header.err = e.message;
		writeOutput({ header: header });
	}
};

readInput(render);
