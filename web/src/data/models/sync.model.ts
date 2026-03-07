import { z } from 'zod';
import { AccountSchema } from './account.model';
import { CategorySchema } from './category.model';
import { TransactionSchema } from './transaction.model';

export const SyncPullRequestSchema = z.object({
  deviceId: z.string(),
  lastSyncAt: z.string().transform((str) => new Date(str)),
});

export type SyncPullRequest = z.infer<typeof SyncPullRequestSchema>;

export const SyncPullResponseSchema = z.object({
  accounts: z.array(AccountSchema),
  categories: z.array(CategorySchema),
  transactions: z.array(TransactionSchema),
  serverTimestamp: z.string().transform((str) => new Date(str)),
});

export type SyncPullResponse = z.infer<typeof SyncPullResponseSchema>;

export const SyncPushRequestSchema = z.object({
  deviceId: z.string(),
  accounts: z.array(AccountSchema),
  categories: z.array(CategorySchema),
  transactions: z.array(TransactionSchema),
  lastSyncAt: z.string().transform((str) => new Date(str)),
});

export type SyncPushRequest = z.infer<typeof SyncPushRequestSchema>;

export const SyncPushResponseSchema = z.object({
  success: z.boolean(),
  serverTimestamp: z.string().transform((str) => new Date(str)),
  conflicts: z.array(z.object({
    entityType: z.string(),
    entityId: z.string(),
    resolution: z.string(),
  })).optional(),
});

export type SyncPushResponse = z.infer<typeof SyncPushResponseSchema>;

export interface SyncState {
  deviceId: string;
  lastSyncAt: Date | null;
  status: 'idle' | 'syncing' | 'success' | 'error';
  lastError?: string;
}
