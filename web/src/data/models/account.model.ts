import { z } from 'zod';
import { BaseEntitySchema } from './base.model';

export const AccountTypeSchema = z.enum(['cash', 'bank', 'credit_card', 'investment', 'other']);
export type AccountType = z.infer<typeof AccountTypeSchema>;

export const AccountSchema = BaseEntitySchema.extend({
  name: z.string(),
  type: AccountTypeSchema,
  currency: z.string().default('CNY'),
  balance: z.number(),
});

export type Account = z.infer<typeof AccountSchema>;

export type AccountCreateInput = Omit<Account, keyof z.infer<typeof BaseEntitySchema>>;
export type AccountUpdateInput = Partial<Omit<Account, 'id' | 'user_id'>>;
