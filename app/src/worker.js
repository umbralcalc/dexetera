self.importScripts('wasm_exec.js');
self.importScripts('google-protobuf.js');
self.importScripts('partition_state_pb.js');

let go;
let socket;
let wasmInstance;
let reconnectInterval = 2000; // 2 seconds
let isConnected = false;
let debugMode = false;
let serverPartitionIndices = [0];

self.onmessage = async function(event) {
    if (event.data.action === 'start') {
        await loadWasm(event.data.wasmBinary);
        debugMode = event.data.debugMode;
        serverPartitionIndices = event.data.serverPartitionIndices;
        startWebSocketClient();
    }
};

async function loadWasm(wasmBinary) {
    go = new Go();
    const result = await WebAssembly.instantiateStreaming(fetch(wasmBinary), go.importObject);
    go.run(result.instance);
    wasmInstance = result.instance;
}

function startWebSocketClient() {
    const socket = new WebSocket('ws://localhost:2112');
    socket.binaryType = 'arraybuffer';
    socket.onopen = function() {
        console.log('WebSocket connection opened.');
        isConnected = true;
        stepSimulation(handlePartitionState, null);
    };
    socket.onmessage = async function(event) {
        if (debugMode) {
            const message = proto.State.deserializeBinary(new Uint8Array(event.data));
            console.log("*******************************************************");
            console.log("Client received values:", message.getValuesList());
        }
        stepSimulation(handlePartitionState, new Uint8Array(event.data));
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
    // Callback function
    function handlePartitionState(data) {
        const message = proto.PartitionState.deserializeBinary(new Uint8Array(data));
        timesteps = message.getCumulativeTimesteps();
        partitionIndex = message.getPartitionIndex();
        state = message.getState().getValuesList();
        if (debugMode) {
            console.log("-------------------------------------------------------");
            console.log("Cumulative Timesteps:", timesteps);
            console.log("Partition Index:", partitionIndex);
            console.log("State:", state);
        }
        // Send the data to the main display thread
        self.postMessage({
            type: 'partitionState',
            data: {
                timesteps: timesteps,
                partitionIndex: partitionIndex,
                state: state,
            }
        });
        // Send the subset of the data to the server
        if (serverPartitionIndices.includes(partitionIndex)) {
            socket.send(data);
        }
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
