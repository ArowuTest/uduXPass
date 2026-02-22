/// <reference lib="webworker" />

/**
 * Service Worker for uduXPass Scanner PWA
 * Provides offline caching and background sync capabilities
 * 
 * Enterprise-grade TypeScript implementation with full type safety
 */

declare const self: ServiceWorkerGlobalScope;

const CACHE_NAME = 'uduxpass-scanner-v1';
const RUNTIME_CACHE = 'uduxpass-scanner-runtime-v1';

// Assets to cache on install
const STATIC_ASSETS: string[] = [
  '/',
  '/index.html',
  '/manifest.json',
  '/icon-192.png',
  '/icon-512.png',
];

/**
 * Install event - cache static assets
 */
self.addEventListener('install', (event: ExtendableEvent) => {
  console.log('[Service Worker] Installing...');
  
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then((cache: Cache) => {
        console.log('[Service Worker] Caching static assets');
        return cache.addAll(STATIC_ASSETS);
      })
      .then(() => {
        console.log('[Service Worker] Installed successfully');
        return self.skipWaiting(); // Activate immediately
      })
      .catch((error: Error) => {
        console.error('[Service Worker] Installation failed:', error);
      })
  );
});

/**
 * Activate event - clean up old caches
 */
self.addEventListener('activate', (event: ExtendableEvent) => {
  console.log('[Service Worker] Activating...');
  
  event.waitUntil(
    caches.keys()
      .then((cacheNames: string[]) => {
        return Promise.all(
          cacheNames
            .filter((cacheName: string) => {
              // Delete old caches
              return cacheName !== CACHE_NAME && cacheName !== RUNTIME_CACHE;
            })
            .map((cacheName: string) => {
              console.log('[Service Worker] Deleting old cache:', cacheName);
              return caches.delete(cacheName);
            })
        );
      })
      .then(() => {
        console.log('[Service Worker] Activated successfully');
        return self.clients.claim(); // Take control immediately
      })
  );
});

/**
 * Fetch event - serve from cache, fallback to network
 */
self.addEventListener('fetch', (event: FetchEvent) => {
  const { request } = event;
  const url = new URL(request.url);

  // Skip non-GET requests
  if (request.method !== 'GET') {
    return;
  }

  // Skip non-http(s) requests
  if (!url.protocol.startsWith('http')) {
    return;
  }

  // API requests - Network first, cache fallback
  if (url.pathname.startsWith('/v1/')) {
    event.respondWith(
      fetch(request)
        .then((response: Response) => {
          // Clone response before caching
          const responseClone = response.clone();
          
          // Cache successful responses
          if (response.ok) {
            caches.open(RUNTIME_CACHE).then((cache: Cache) => {
              cache.put(request, responseClone);
            });
          }
          
          return response;
        })
        .catch(() => {
          // Network failed, try cache
          return caches.match(request)
            .then((cachedResponse: Response | undefined) => {
              if (cachedResponse) {
                console.log('[Service Worker] Serving API from cache:', url.pathname);
                return cachedResponse;
              }
              
              // Return offline response
              return new Response(
                JSON.stringify({
                  success: false,
                  error: 'Offline - No cached data available',
                  offline: true
                }),
                {
                  status: 503,
                  statusText: 'Service Unavailable',
                  headers: { 'Content-Type': 'application/json' }
                }
              );
            });
        })
    );
    return;
  }

  // Static assets - Cache first, network fallback
  event.respondWith(
    caches.match(request)
      .then((cachedResponse: Response | undefined) => {
        if (cachedResponse) {
          console.log('[Service Worker] Serving from cache:', url.pathname);
          return cachedResponse;
        }

        // Not in cache, fetch from network
        return fetch(request)
          .then((response: Response) => {
            // Clone response before caching
            const responseClone = response.clone();
            
            // Cache successful responses
            if (response.ok) {
              caches.open(RUNTIME_CACHE).then((cache: Cache) => {
                cache.put(request, responseClone);
              });
            }
            
            return response;
          })
          .catch((error: Error) => {
            console.error('[Service Worker] Fetch failed:', error);
            
            // Return offline page for navigation requests
            if (request.mode === 'navigate') {
              return caches.match('/index.html').then((response: Response | undefined) => {
                return response || new Response('Offline', { status: 503 });
              });
            }
            
            throw error;
          });
      })
  );
});

/**
 * Background Sync - for offline ticket validations
 */
interface SyncEvent extends ExtendableEvent {
  tag: string;
}

self.addEventListener('sync', (event: Event) => {
  const syncEvent = event as SyncEvent;
  console.log('[Service Worker] Background sync:', syncEvent.tag);
  
  if (syncEvent.tag === 'sync-validations') {
    syncEvent.waitUntil(syncOfflineValidations());
  }
});

/**
 * Sync offline validations when back online
 */
async function syncOfflineValidations(): Promise<void> {
  try {
    console.log('[Service Worker] Syncing offline validations...');
    
    // Get offline validations from IndexedDB (will be implemented in Phase 2)
    // For now, just log
    console.log('[Service Worker] Sync complete');
  } catch (error) {
    console.error('[Service Worker] Sync failed:', error);
    throw error; // Retry sync
  }
}

/**
 * Message handler - for communication with main app
 */
interface ServiceWorkerMessage {
  type: string;
  urls?: string[];
}

self.addEventListener('message', (event: ExtendableMessageEvent) => {
  const data = event.data as ServiceWorkerMessage;
  console.log('[Service Worker] Message received:', data);
  
  if (data && data.type === 'SKIP_WAITING') {
    self.skipWaiting();
  }
  
  if (data && data.type === 'CACHE_URLS') {
    const urls = data.urls || [];
    event.waitUntil(
      caches.open(RUNTIME_CACHE)
        .then((cache: Cache) => cache.addAll(urls))
        .then(() => {
          console.log('[Service Worker] Cached additional URLs');
        })
    );
  }
});

export {};
