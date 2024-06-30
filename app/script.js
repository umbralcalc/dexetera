document.getElementById('startButton').addEventListener('click', () => {
    const worker = new Worker('worker.js');
    worker.postMessage({ action: 'start' });
});
