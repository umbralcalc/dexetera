self.importScripts('wasm_exec.js');
self.importScripts('google-protobuf.js');
self.importScripts('partition_state_pb.js');

let wasmInstance;
let go;

self.onmessage = async function(event) {
    if (event.data.action === 'start') {
        await loadWasm();
        startWebSocketClient();
    }
};

async function loadWasm() {
    go = new Go();
    const result = await WebAssembly.instantiateStreaming(fetch("example_sim.wasm"), go.importObject);
    go.run(result.instance);
    wasmInstance = result.instance;
}

function startWebSocketClient() {
    const socket = new WebSocket('ws://localhost:2112');
    socket.binaryType = 'arraybuffer';
    socket.onopen = function() {
        console.log('WebSocket connection opened.');
        stepSimulation(handlePartitionState, null);
    };
    socket.onmessage = async function(event) {
        // const message = proto.State.deserializeBinary(new Uint8Array(event.data));
        // console.log(message.getValuesList());
        stepSimulation(handlePartitionState, new Uint8Array(event.data));
    };
    socket.onclose = function() {
        console.log('WebSocket connection closed.');
    };
    socket.onerror = function(error) {
        console.error('WebSocket error:', error);
    };
    // Callback function
    function handlePartitionState(data) {
        const message = proto.PartitionState.deserializeBinary(new Uint8Array(data));
        console.log("-------------------------------------------------------");
        console.log("Cumulative Timesteps:", message.getCumulativeTimesteps());
        console.log("Partition Index:", message.getPartitionIndex());
        console.log("State:", message.getState().getValuesList());
        socket.send(data);
    };
}
