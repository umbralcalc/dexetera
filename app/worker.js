const CACHE_NAME = 'my-app-cache-v1';
const ASSETS_TO_CACHE = [
    'index.html',
    'example_sim.wasm',
    'wasm_exec.js',
    'styles.css',
    'script.js'
];

self.addEventListener('install', event => {
    console.log('Service Worker: Install event');
    event.waitUntil(
        caches.open(CACHE_NAME)
            .then(cache => {
                console.log('Opened cache');
                return cache.addAll(ASSETS_TO_CACHE);
            })
            .then(() => {
                console.log('All assets cached');
                self.skipWaiting();
            })
            .catch(err => {
                console.error('Failed to cache assets during install', err);
            })
    );
});

self.addEventListener('activate', event => {
    console.log('Service Worker: Activate event');
    event.waitUntil(
        self.clients.claim().then(() => {
            console.log('Service Worker activated and clients claimed');
        })
    );
});

self.addEventListener('fetch', event => {
    console.log('Service Worker: Fetch event for', event.request.url);
    event.respondWith(
        caches.match(event.request)
            .then(response => {
                return response || fetch(event.request);
            })
    );
});

self.addEventListener('message', event => {
    console.log('Service Worker: Received message', event.data);
    if (event.data === 'cache-assets') {
        caches.open(CACHE_NAME).then(cache => {
            return cache.addAll(ASSETS_TO_CACHE).then(() => {
                console.log('Assets cached via message');
                self.clients.matchAll().then(clients => {
                    clients.forEach(client => {
                        console.log('Service Worker: Sending assets-cached message to client');
                        client.postMessage('assets-cached');
                    });
                });
            });
        }).catch(err => {
            console.error('Failed to cache assets on message', err);
        });
    }
});
