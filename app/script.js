if ('serviceWorker' in navigator) {
    navigator.serviceWorker.register('/worker.js')
        .then(registration => {
            console.log('Service Worker registered with scope:', registration.scope);
            return navigator.serviceWorker.ready;
        })
        .then(registration => {
            console.log('Service Worker is ready and controlling the page');
        })
        .catch(error => {
            console.error('Service Worker registration failed:', error);
        });

    navigator.serviceWorker.addEventListener('message', event => {
        console.log('Main thread: Received message from Service Worker:', event.data);
        if (event.data === 'assets-cached') {
            console.log('Assets have been cached. Running WebAssembly...');
            runWasm();
        } else {
            console.log('Received unknown message from Service Worker:', event.data);
        }
    });
} else {
    console.error('Service Worker not supported in this browser');
}

// Function to trigger downloading assets
function downloadAssets() {
    console.log('Main thread: Sending message to Service Worker to cache assets...');
    if (navigator.serviceWorker.controller) {
        navigator.serviceWorker.controller.postMessage('cache-assets');
    } else {
        console.error('No active Service Worker controller found, retrying...');
        navigator.serviceWorker.ready.then(() => {
            navigator.serviceWorker.controller.postMessage('cache-assets');
        }).catch(err => {
            console.error('Failed to find an active Service Worker controller', err);
        });
    }
}

// Function to run the WebAssembly code
function runWasm() {
    const go = new Go();
    WebAssembly.instantiateStreaming(fetch("example_sim.wasm"), go.importObject).then(result => {
        go.run(result.instance);

        // JavaScript callback function
        function handleData(data) {
            console.log("Data from Go:", data);
            // Process the data as needed
        }

        // Call the exported Go function
        loop(handleData);
    }).catch(err => {
        console.error("Error loading WebAssembly:", err);
    });
}

// Add event listener to button
document.getElementById('downloadButton').addEventListener('click', downloadAssets);
