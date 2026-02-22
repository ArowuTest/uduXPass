/**
 * Offline Validation Service
 * Handles ticket validation when offline using cached data
 */

import {
  getCachedTicket,
  saveOfflineValidation,
  addToSyncQueue,
  type OfflineTicket,
  type OfflineValidation,
} from './offlineDB';
import { isOnline } from './registerServiceWorker';

export interface ValidationResult {
  success: boolean;
  message: string;
  offline: boolean;
  ticket?: {
    id: string;
    qr_code: string;
    status: string;
    tier_name?: string;
    customer_name?: string;
  };
  error?: string;
}

/**
 * Validate ticket offline using cached data
 */
export async function validateTicketOffline(
  qrCode: string,
  sessionId: string
): Promise<ValidationResult> {
  try {
    console.log('[Offline] Validating ticket offline:', qrCode);

    // Get cached ticket
    const cachedTicket = await getCachedTicket(qrCode);

    if (!cachedTicket) {
      // Ticket not in cache
      const validation: OfflineValidation = {
        id: `val_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
        qr_code: qrCode,
        session_id: sessionId,
        validated_at: Date.now(),
        status: 'error',
        message: 'Ticket not found in offline cache',
        synced: false,
      };

      await saveOfflineValidation(validation);

      return {
        success: false,
        offline: true,
        message: 'Ticket not found in offline cache',
        error: 'NOT_CACHED',
      };
    }

    // Check ticket status
    if (cachedTicket.status === 'used') {
      // Ticket already used
      const validation: OfflineValidation = {
        id: `val_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
        qr_code: qrCode,
        session_id: sessionId,
        validated_at: Date.now(),
        status: 'error',
        message: 'Ticket already used (offline validation)',
        synced: false,
      };

      await saveOfflineValidation(validation);

      return {
        success: false,
        offline: true,
        message: 'Ticket already used',
        error: 'ALREADY_USED',
        ticket: {
          id: cachedTicket.ticket_id,
          qr_code: cachedTicket.qr_code,
          status: cachedTicket.status,
          tier_name: cachedTicket.tier_name,
          customer_name: cachedTicket.customer_name,
        },
      };
    }

    if (cachedTicket.status === 'invalid') {
      // Invalid ticket
      const validation: OfflineValidation = {
        id: `val_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
        qr_code: qrCode,
        session_id: sessionId,
        validated_at: Date.now(),
        status: 'error',
        message: 'Invalid ticket (offline validation)',
        synced: false,
      };

      await saveOfflineValidation(validation);

      return {
        success: false,
        offline: true,
        message: 'Invalid ticket',
        error: 'INVALID',
        ticket: {
          id: cachedTicket.ticket_id,
          qr_code: cachedTicket.qr_code,
          status: cachedTicket.status,
          tier_name: cachedTicket.tier_name,
          customer_name: cachedTicket.customer_name,
        },
      };
    }

    // Ticket is valid - mark as used locally
    cachedTicket.status = 'used';
    
    // Save validation
    const validation: OfflineValidation = {
      id: `val_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
      qr_code: qrCode,
      session_id: sessionId,
      validated_at: Date.now(),
      status: 'success',
      message: 'Ticket validated successfully (offline)',
      synced: false,
    };

    await saveOfflineValidation(validation);

    // Add to sync queue for later sync
    await addToSyncQueue({
      type: 'validation',
      data: {
        qr_code: qrCode,
        session_id: sessionId,
        validated_at: validation.validated_at,
      },
      created_at: Date.now(),
      retry_count: 0,
    });

    console.log('[Offline] Ticket validated successfully offline');

    return {
      success: true,
      offline: true,
      message: 'Ticket validated successfully (offline mode)',
      ticket: {
        id: cachedTicket.ticket_id,
        qr_code: cachedTicket.qr_code,
        status: 'valid',
        tier_name: cachedTicket.tier_name,
        customer_name: cachedTicket.customer_name,
      },
    };
  } catch (error) {
    console.error('[Offline] Validation error:', error);
    
    return {
      success: false,
      offline: true,
      message: 'Offline validation failed',
      error: error instanceof Error ? error.message : 'Unknown error',
    };
  }
}

/**
 * Check if ticket should be validated offline
 */
export function shouldUseOfflineValidation(): boolean {
  return !isOnline();
}

/**
 * Sync offline validations to server
 */
export async function syncOfflineValidations(
  apiValidateFunction: (qrCode: string, sessionId: string) => Promise<any>
): Promise<{
  synced: number;
  failed: number;
  errors: string[];
}> {
  const { getUnsyncedValidations, markValidationSynced } = await import('./offlineDB');
  
  const unsyncedValidations = await getUnsyncedValidations();
  
  console.log(`[Offline] Syncing ${unsyncedValidations.length} offline validations...`);
  
  let synced = 0;
  let failed = 0;
  const errors: string[] = [];
  
  for (const validation of unsyncedValidations) {
    try {
      // Try to sync with server
      await apiValidateFunction(validation.qr_code, validation.session_id);
      
      // Mark as synced
      await markValidationSynced(validation.id);
      synced++;
      
      console.log('[Offline] Synced validation:', validation.id);
    } catch (error) {
      failed++;
      const errorMsg = error instanceof Error ? error.message : 'Unknown error';
      errors.push(`${validation.qr_code}: ${errorMsg}`);
      console.error('[Offline] Failed to sync validation:', validation.id, error);
    }
  }
  
  console.log(`[Offline] Sync complete: ${synced} synced, ${failed} failed`);
  
  return { synced, failed, errors };
}
