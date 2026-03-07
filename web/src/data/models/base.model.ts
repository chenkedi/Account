import { z } from 'zod';

export const BaseEntitySchema = z.object({
  id: z.string().uuid(),
  user_id: z.string().uuid(),
  created_at: z.string().transform((str) => new Date(str)),
  updated_at: z.string().transform((str) => new Date(str)),
  last_modified_at: z.string().transform((str) => new Date(str)),
  version: z.number().int(),
  is_deleted: z.boolean(),
});

export type BaseEntity = z.infer<typeof BaseEntitySchema>;

export interface LwwEntity {
  id: string;
  last_modified_at: Date;
  is_deleted: boolean;
  version: number;
}
