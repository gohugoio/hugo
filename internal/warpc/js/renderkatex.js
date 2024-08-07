import { readInput, writeOutput } from './common';
import katex from 'katex';

const render = function (input) {
	const data = input.data;
	const expression = data.expression;
	const options = data.options;
	writeOutput({ header: input.header, data: { output: katex.renderToString(expression, options) } });
};

readInput(render);
