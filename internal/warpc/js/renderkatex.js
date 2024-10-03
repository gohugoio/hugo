import { readInput, writeOutput } from './common';
import katex from 'katex';

const render = function (input) {
	const data = input.data;
	const expression = data.expression;
	const options = data.options;
	const header = input.header;
	// Any error thrown here will be caught by the common.js readInput function.
	const output = katex.renderToString(expression, options);
	writeOutput({ header: header, data: { output: output } });
};

readInput(render);
