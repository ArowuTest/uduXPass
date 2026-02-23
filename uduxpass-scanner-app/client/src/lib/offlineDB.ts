/**
 * Offline Database Utility using IndexedDB
 * Stores tickets and validations for offline scanning
 */

const DB_NAME = 'uduxpass-scanner-offline';
const DB_VERSION = 1;

// Store names
const STORES = {
  TICKETS: 'tickets',
  VALIDATIONS: 'validations',
  SYNC_QUEUE: 'sync_queue',
} as const;

export interface OfflineTicket {
  qr_code: string;
  ticket_id: string;
  event_id: string;
  tier_name: string;
  customer_name: string;
  status: 'valid' | 'used' | 'invalid';
  cached_at: number;
}

export interface OfflineValidation {
  id: string;
  qr_code: string;
  session_id: string;
  validated_at: number;
  status: 'success' | 'error';
  message: string;
  synced: boolean;
}

export interface SyncQueueItem {
  id: string;
  type: 'validation';
  data: any;
  created_at: number;
  retry_count: number;
}

/**
 * Initialize IndexedDB
 */
export async function initOfflineDB(): Promise<IDBDatabase> {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open(DB_NAME, DB_VERSION);

    request.onerror = () => {
      console.error('[OfflineDB] Failed to open database:', request.error);
      reject(request.error);
    };

    request.onsuccess = () => {
      console.log('[OfflineDB] Database opened successfully');
      resolve(request.result);
    };

    request.onupgradeneeded = (event) => {
      const db = (event.target as IDBOpenDBRequest).result;
      console.log('[OfflineDB] Upgrading database schema...');

      // Tickets store - for offline validation
      if (!db.objectStoreNames.contains(STORES.TICKETS)) {
        const ticketStore = db.createObjectStore(STORES.TICKETS, { keyPath: 'qr_code' });
        ticketStore.createIndex('event_id', 'event_id', { unique: false });
        ticketStore.createIndex('status', 'status', { unique: false });
        ticketStore.createIndex('cached_at', 'cached_at', { unique: false });
        console.log('[OfflineDB] Created tickets store');
      }

      // Validations store - for offline scan history
      if (!db.objectStoreNames.contains(STORES.VALIDATIONS)) {
        const validationStore = db.createObjectStore(STORES.VALIDATIONS, { keyPath: 'id' });
        validationStore.createIndex('session_id', 'session_id', { unique: false });
        validationStore.createIndex('synced', 'synced', { unique: false });
        validationStore.createIndex('validated_at', 'validated_at', { unique: false });
        console.log('[OfflineDB] Created validations store');
      }

      // Sync queue store - for pending syncs
      if (!db.objectStoreNames.contains(STORES.SYNC_QUEUE)) {
        const syncStore = db.createObjectStore(STORES.SYNC_QUEUE, { keyPath: 'id' });
        syncStore.createIndex('created_at', 'created_at', { unique: false });
        syncStore.createIndex('retry_count', 'retry_count', { unique: false });
        console.log('[OfflineDB] Created sync queue store');
      }
    };
  });
}

/**
 * Cache ticket for offline validation
 */
export async function cacheTicket(ticket: OfflineTicket): Promise<void> {
  const db = await initOfflineDB();
  
  return new Promise((resolve, reject) => {
    const transaction = db.transaction([STORES.TICKETS], 'readwrite');
    const store = transaction.objectStore(STORES.TICKETS);
    
    const request = store.put(ticket);
    
    request.onsuccess = () => {
      console.log('[OfflineDB] Ticket cached:', ticket.qr_code);
      resolve();
    };
    
    request.onerror = () => {
      console.error('[OfflineDB] Failed to cache ticket:', request.error);
      reject(request.error);
    };
  });
}

/**
 * Get cached ticket by QR code
 */
export async function getCachedTicket(qrCode: string): Promise<OfflineTicket | null> {
  const db = await initOfflineDB();
  
  return new Promise((resolve, reject) => {
    const transaction = db.transaction([STORES.TICKETS], 'readonly');
    const store = transaction.objectStore(STORES.TICKETS);
    
    const request = store.get(qrCode);
    
    request.onsuccess = () => {
      resolve(request.result || null);
    };
    
    request.onerror = () => {
      console.error('[OfflineDB] Failed to get ticket:', request.error);
      reject(request.error);
    };
  });
}

/**
 * Cache tickets for an event (bulk operation)
 */
export async function cacheEventTickets(tickets: OfflineTicket[]): Promise<void> {
  const db = await initOfflineDB();
  
  return new Promise((resolve, reject) => {
    const transaction = db.transaction([STORES.TICKETS], 'readwrite');
    const store = transaction.objectStore(STORES.TICKETS);
    
    let completed = 0;
    const total = tickets.length;
    
    tickets.forEach((ticket) => {
      const request = store.put(ticket);
      
      request.onsuccess = () => {
        completed++;
        if (completed === total) {
          console.log(`[OfflineDB] Cached ${total} tickets`);
          resolve();
        }
      };
      
      request.onerror = () => {
        console.error('[OfflineDB] Failed to cache ticket:', request.error);
        reject(request.error);
      };
    });
    
    if (total === 0) {
      resolve();
    }
  });
}

/**
 * Save offline validation
 */
export async function saveOfflineValidation(validation: OfflineValidation): Promise<void> {
  const db = await initOfflineDB();
  
  return new Promise((resolve, reject) => {
    const transaction = db.transaction([STORES.VALIDATIONS], 'readwrite');
    const store = transaction.objectStore(STORES.VALIDATIONS);
    
    const request = store.put(validation);
    
    request.onsuccess = () => {
      console.log('[OfflineDB] Validation saved:', validation.id);
      resolve();
    };
    
    request.onerror = () => {
      console.error('[OfflineDB] Failed to save validation:', request.error);
      reject(request.error);
    };
  });
}

/**
 * Get unsynced validations
 */
export async function getUnsyncedValidations(): Promise<OfflineValidation[]> {
  const db = await initOfflineDB();
  
  return new Promise((resolve, reject) => {
    const transaction = db.transaction([STORES.VALIDATIONS], 'readonly');
    const store = transaction.objectStore(STORES.VALIDATIONS);
    const index = store.index('synced');
    
    const request = index.getAll(false);
    
    request.onsuccess = () => {
      resolve(request.result || []);
    };
    
    request.onerror = () => {
      console.error('[OfflineDB] Failed to get unsynced validations:', request.error);
      reject(request.error);
    };
  });
}

/**
 * Mark validation as synced
 */
export async function markValidationSynced(validationId: string): Promise<void> {
  const db = await initOfflineDB();
  
  return new Promise((resolve, reject) => {
    const transaction = db.transaction([STORES.VALIDATIONS], 'readwrite');
    const store = transaction.objectStore(STORES.VALIDATIONS);
    
    const getRequest = store.get(validationId);
    
    getRequest.onsuccess = () => {
      const validation = getRequest.result;
      if (validation) {
        validation.synced = true;
        const putRequest = store.put(validation);
        
        putRequest.onsuccess = () => {
          console.log('[OfflineDB] Validation marked as synced:', validationId);
          resolve();
        };
        
        putRequest.onerror = () => {
          reject(putRequest.error);
        };
      } else {
        resolve();
      }
    };
    
    getRequest.onerror = () => {
      reject(getRequest.error);
    };
  });
}

/**
 * Add item to sync queue
 */
export async function addToSyncQueue(item: Omit<SyncQueueItem, 'id'>): Promise<void> {
  const db = await initOfflineDB();
  
  const queueItem: SyncQueueItem = {
    ...item,
    id: `sync_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
  };
  
  return new Promise((resolve, reject) => {
    const transaction = db.transaction([STORES.SYNC_QUEUE], 'readwrite');
    const store = transaction.objectStore(STORES.SYNC_QUEUE);
    
    const request = store.put(queueItem);
    
    request.onsuccess = () => {
      console.log('[OfflineDB] Added to sync queue:', queueItem.id);
      resolve();
    };
    
    request.onerror = () => {
      console.error('[OfflineDB] Failed to add to sync queue:', request.error);
      reject(request.error);
    };
  });
}

/**
 * Get all items from sync queue
 */
export async function getSyncQueue(): Promise<SyncQueueItem[]> {
  const db = await initOfflineDB();
  
  return new Promise((resolve, reject) => {
    const transaction = db.transaction([STORES.SYNC_QUEUE], 'readonly');
    const store = transaction.objectStore(STORES.SYNC_QUEUE);
    
    const request = store.getAll();
    
    request.onsuccess = () => {
      resolve(request.result || []);
    };
    
    request.onerror = () => {
      console.error('[OfflineDB] Failed to get sync queue:', request.error);
      reject(request.error);
    };
  });
}

/**
 * Remove item from sync queue
 */
export async function removeFromSyncQueue(itemId: string): Promise<void> {
  const db = await initOfflineDB();
  
  return new Promise((resolve, reject) => {
    const transaction = db.transaction([STORES.SYNC_QUEUE], 'readwrite');
    const store = transaction.objectStore(STORES.SYNC_QUEUE);
    
    const request = store.delete(itemId);
    
    request.onsuccess = () => {
      console.log('[OfflineDB] Removed from sync queue:', itemId);
      resolve();
    };
    
    request.onerror = () => {
      console.error('[OfflineDB] Failed to remove from sync queue:', request.error);
      reject(request.error);
    };
  });
}

/**
 * Clear all offline data (for logout or reset)
 */
export async function clearOfflineData(): Promise<void> {
  const db = await initOfflineDB();
  
  return new Promise((resolve, reject) => {
    const transaction = db.transaction(
      [STORES.TICKETS, STORES.VALIDATIONS, STORES.SYNC_QUEUE],
      'readwrite'
    );
    
    const promises: Promise<void>[] = [];
    
    [STORES.TICKETS, STORES.VALIDATIONS, STORES.SYNC_QUEUE].forEach((storeName) => {
      const store = transaction.objectStore(storeName);
      promises.push(
        new Promise((res, rej) => {
          const request = store.clear();
          request.onsuccess = () => res();
          request.onerror = () => rej(request.error);
        })
      );
    });
    
    Promise.all(promises)
      .then(() => {
        console.log('[OfflineDB] All offline data cleared');
        resolve();
      })
      .catch(reject);
  });
}

/**
 * Get database statistics
 */
export async function getOfflineStats(): Promise<{
  cachedTickets: number;
  validations: number;
  unsyncedValidations: number;
  queueItems: number;
}> {
  const db = await initOfflineDB();
  
  return new Promise((resolve, reject) => {
    const transaction = db.transaction(
      [STORES.TICKETS, STORES.VALIDATIONS, STORES.SYNC_QUEUE],
      'readonly'
    );
    
    const stats = {
      cachedTickets: 0,
      validations: 0,
      unsyncedValidations: 0,
      queueItems: 0,
    };
    
    // Count tickets
    const ticketsRequest = transaction.objectStore(STORES.TICKETS).count();
    ticketsRequest.onsuccess = () => {
      stats.cachedTickets = ticketsRequest.result;
    };
    
    // Count validations
    const validationsRequest = transaction.objectStore(STORES.VALIDATIONS).count();
    validationsRequest.onsuccess = () => {
      stats.validations = validationsRequest.result;
    };
    
    // Count unsynced validations
    const unsyncedRequest = transaction
      .objectStore(STORES.VALIDATIONS)
      .index('synced')
      .count(false);
    unsyncedRequest.onsuccess = () => {
      stats.unsyncedValidations = unsyncedRequest.result;
    };
    
    // Count queue items
    const queueRequest = transaction.objectStore(STORES.SYNC_QUEUE).count();
    queueRequest.onsuccess = () => {
      stats.queueItems = queueRequest.result;
    };
    
    transaction.oncomplete = () => {
      resolve(stats);
    };
    
    transaction.onerror = () => {
      reject(transaction.error);
    };
  });
}

// Export all functions as offlineDB object
export const offlineDB = {
  initOfflineDB,
  cacheTickets,
  getTicketByQR,
  saveValidation,
  getUnsyncedValidations,
  markValidationSynced,
  addToSyncQueue,
  getSyncQueue,
  removeFromSyncQueue,
  clearCache,
  getStats,
};
