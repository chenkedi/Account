import { z } from 'zod';

// Backend User schema (snake_case)
export const BackendUserSchema = z.object({
  id: z.string().uuid(),
  email: z.string().email(),
  created_at: z.string().transform((str) => new Date(str)),
  updated_at: z.string().transform((str) => new Date(str)),
});

export type BackendUser = z.infer<typeof BackendUserSchema>;

// Frontend User type (camelCase)
export interface User {
  id: string;
  email: string;
  createdAt: Date;
}

// Auth Response from backend (snake_case)
export const AuthResponseSchema = z.object({
  user: BackendUserSchema,
  access_token: z.string(),
  token_type: z.string(),
  expires_in: z.number(),
});

export type AuthResponse = z.infer<typeof AuthResponseSchema>;

export type LoginInput = {
  email: string;
  password: string;
};

export type RegisterInput = {
  email: string;
  password: string;
};
