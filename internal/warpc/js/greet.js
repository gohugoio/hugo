import { readInput, writeOutput } from './common';

const greet = function (input) {
	writeOutput({ header: input.header, data: { greeting: 'Hello ' + input.data.name + '!' } });
};

console.log('Greet module loaded');

readInput(greet);
