import { v4 as uuidv4 } from 'uuid';

const DEVICE_ID_KEY = 'sync_device_id';
const LAST_SYNC_AT_KEY = 'sync_last_sync_at';

export class SyncStorage {
  async getOrCreateDeviceId(): Promise<string> {
    let deviceId = localStorage.getItem(DEVICE_ID_KEY);
    if (!deviceId) {
      deviceId = uuidv4();
      localStorage.setItem(DEVICE_ID_KEY, deviceId);
    }
    return deviceId;
  }

  async getDeviceId(): Promise<string | null> {
    return localStorage.getItem(DEVICE_ID_KEY);
  }

  async getLastSyncAt(): Promise<Date | null> {
    const timestamp = localStorage.getItem(LAST_SYNC_AT_KEY);
    if (!timestamp) return null;
    try {
      return new Date(timestamp);
    } catch {
      return null;
    }
  }

  async saveLastSyncAt(date: Date): Promise<void> {
    localStorage.setItem(LAST_SYNC_AT_KEY, date.toISOString());
  }

  async clear(): Promise<void> {
    localStorage.removeItem(LAST_SYNC_AT_KEY);
  }
}

export const syncStorage = new SyncStorage();
