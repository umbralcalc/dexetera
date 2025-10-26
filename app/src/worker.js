self.importScripts('wasm_exec.js');
self.importScripts('google-protobuf.js');
self.importScripts('partition_state_pb.js');

let go;
let socket;
let wasmInstance;
let stopAtSimTime;
let reconnectInterval = 2000; // 2 seconds
let isConnected = false;
let debugMode = false;
let serverPartitionNames = [];
let timesteps = 0;
let partitionName = '';
let state = [];

self.onmessage = async function(event) {
    if (event.data.action === 'start') {
        await loadWasm(event.data.wasmBinary);
        debugMode = event.data.debugMode;
        stopAtSimTime = event.data.stopAtSimTime;
        serverPartitionNames = event.data.serverPartitionNames;
        startWebSocketClient();
    }
};

async function loadWasm(wasmBinary) {
    try {
        go = new Go();
        const result = await WebAssembly.instantiateStreaming(fetch(wasmBinary), go.importObject);
        go.run(result.instance);
        wasmInstance = result.instance;
        console.log('WebAssembly module loaded successfully');
        console.log('stepSimulation function available:', typeof stepSimulation);
    } catch (error) {
        console.error('Error loading WebAssembly module:', error);
        self.postMessage({
            type: 'error',
            data: 'Failed to load WebAssembly module: ' + error.message
        });
    }
}

// Global callback function for handling partition state from Go
function handlePartitionState(data) {
    // data is raw protobuf bytes from Go, need to deserialize it
    const message = proto.PartitionState.deserializeBinary(data);
    
    if (stopAtSimTime <= message.getCumulativeTimesteps()) return;
    timesteps = message.getCumulativeTimesteps();
    partitionName = message.getPartitionName();
    state = message.getStateList();
    
    if (debugMode) {
        console.log("-------------------------------------------------------");
        console.log("Cumulative Timesteps:", timesteps);
        console.log("Partition Name:", partitionName);
        console.log("State:", state);
    }
    // Send the data to the main display thread
    self.postMessage({
        type: 'partitionState',
        data: {
            timesteps: timesteps,
            partitionName: partitionName,
            state: {
                values: state
            },
        }
    });
    // Send the subset of the data to the server only if WebSocket is ready
    if (serverPartitionNames.includes(partitionName) && socket && socket.readyState === WebSocket.OPEN) {
        socket.send(data);
    }
}

function startWebSocketClient() {
    socket = new WebSocket('ws://localhost:2112');
    socket.binaryType = 'arraybuffer';
    socket.onopen = function() {
        console.log('WebSocket connection opened.');
        isConnected = true;
        // Register the callback function with Go
        stepSimulation(handlePartitionState, null);
        
        // Send initial message to start the simulation
        const initialMessage = new Uint8Array([1, 2, 3, 4, 5]); // Simple initial data
        socket.send(initialMessage);
    };
    socket.onmessage = async function(event) {
        if (debugMode) {
            const message = proto.State.deserializeBinary(new Uint8Array(event.data));
            console.log("*******************************************************");
            console.log("Client received values:", message.getValuesList());
        }
        // Send WebSocket data to Go for processing
        stepSimulation(null, new Uint8Array(event.data));
    };
    socket.onclose = function() {
        if (debugMode) {
            console.log('WebSocket connection closed.');
        }
        isConnected = false;
        reconnect();
    };
    socket.onerror = function(error) {
        if (debugMode) {
            console.error('WebSocket error:', error);
        }
        isConnected = false;
        reconnect();
    };
}

function reconnect() {
    if (!isConnected) {
        if (debugMode) {
            console.log(`Attempting to reconnect in ${reconnectInterval / 1000} seconds...`);
        }
        setTimeout(startWebSocketClient, reconnectInterval);
    }
}
