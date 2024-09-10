let performanceNowShim = () => Date.now();
export { performanceNowShim as 'performance.now' };
