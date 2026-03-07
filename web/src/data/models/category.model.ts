import { z } from 'zod';
import { BaseEntitySchema } from './base.model';

export const CategoryTypeSchema = z.enum(['income', 'expense']);
export type CategoryType = z.infer<typeof CategoryTypeSchema>;

export const CategorySchema = BaseEntitySchema.extend({
  name: z.string(),
  type: CategoryTypeSchema,
  parent_id: z.string().uuid().nullable(),
  icon: z.string().nullable(),
});

export type Category = z.infer<typeof CategorySchema>;

export type CategoryCreateInput = Omit<Category, keyof z.infer<typeof BaseEntitySchema>>;
export type CategoryUpdateInput = Partial<Omit<Category, 'id' | 'user_id'>>;
