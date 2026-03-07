import { z } from 'zod';
import { BaseEntitySchema } from './base.model';

export const TransactionTypeSchema = z.enum(['income', 'expense', 'transfer']);
export type TransactionType = z.infer<typeof TransactionTypeSchema>;

export const TransactionSchema = BaseEntitySchema.extend({
  account_id: z.string().uuid(),
  category_id: z.string().uuid().nullable(),
  type: TransactionTypeSchema,
  amount: z.number(),
  currency: z.string().default('CNY'),
  note: z.string().nullable(),
  transaction_date: z.string().transform((str) => new Date(str)),
});

export type Transaction = z.infer<typeof TransactionSchema>;

export type TransactionCreateInput = {
  account_id: string;
  category_id?: string | null;
  type: TransactionType;
  amount: number;
  currency?: string;
  note?: string | null;
  transaction_date: string | Date;
};

export type TransactionUpdateInput = Partial<Omit<Transaction, 'id' | 'user_id'>>;
