import { readInput, writeOutput } from './common';
import katex from 'katex';
import 'katex/contrib/mhchem/mhchem.js';

const render = function (input) {
	const data = input.data;
	const expression = data.expression;
	const options = data.options;
	const header = input.header;
	header.warnings = [];

	if (options.strict == 'warn') {
		// By default, KaTeX will write to console.warn, that's a little hard to handle.
		options.strict = (errorCode, errorMsg) => {
			header.warnings.push(
				`katex: LaTeX-incompatible input and strict mode is set to 'warn': ${errorMsg} [${errorCode}]`,
			);
		};
	}
	// Any error thrown here will be caught by the common.js readInput function.
	const output = katex.renderToString(expression, options);
	writeOutput({ header: header, data: { output: output } });
};

readInput(render);
