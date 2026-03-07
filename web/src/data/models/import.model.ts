import { z } from 'zod';

export const ImportSourceSchema = z.enum(['alipay', 'wechat', 'bank', 'generic']);
export type ImportSource = z.infer<typeof ImportSourceSchema>;

export const ImportSourceInfoSchema = z.object({
  id: ImportSourceSchema,
  name: z.string(),
  description: z.string(),
  supportedExtensions: z.array(z.string()),
});

export type ImportSourceInfo = z.infer<typeof ImportSourceInfoSchema>;

// Server-side snake_case models
export const ParsedTransactionServerSchema = z.object({
  line_number: z.number(),
  transaction_date: z.string(),
  amount: z.number(),
  type: z.enum(['income', 'expense']),
  note: z.string().nullable(),
  counterparty: z.string().nullable(),
  account_name: z.string().nullable(),
  category_hint: z.string().nullable(),
  is_duplicate: z.boolean(),
  can_be_imported: z.boolean(),
  import_warning: z.string().nullable(),
});

export type ParsedTransactionServer = z.infer<typeof ParsedTransactionServerSchema>;

export const AccountSuggestionServerSchema = z.object({
  id: z.string().uuid(),
  name: z.string(),
});

export type AccountSuggestionServer = z.infer<typeof AccountSuggestionServerSchema>;

export const ImportPreviewServerSchema = z.object({
  job_id: z.string(),
  source: ImportSourceSchema,
  file_name: z.string(),
  total_rows: z.number(),
  valid_rows: z.number(),
  duplicate_rows: z.number(),
  transactions: z.array(ParsedTransactionServerSchema),
  account_suggestions: z.record(z.array(AccountSuggestionServerSchema)).optional(),
  categories: z.array(z.any()).optional(),
});

export type ImportPreviewServer = z.infer<typeof ImportPreviewServerSchema>;

export const ImportResultServerSchema = z.object({
  job_id: z.string(),
  total_rows: z.number(),
  imported_rows: z.number(),
  skipped_rows: z.number(),
  failed_rows: z.number(),
  imported_ids: z.array(z.string().uuid()).optional(),
  errors: z.array(z.object({
    line_number: z.number(),
    error: z.string(),
  })).optional(),
});

export type ImportResultServer = z.infer<typeof ImportResultServerSchema>;

// Frontend camelCase models (for UI)
export interface ParsedTransaction {
  rawIndex: number;
  date: string;
  amount: number;
  type: 'income' | 'expense';
  description: string | null;
  counterparty: string | null;
  accountSuggestion: string | null;
  categorySuggestion: string | null;
  isDuplicate: boolean;
  canBeImported: boolean;
}

export interface ImportPreview {
  jobId: string;
  source: ImportSource;
  fileName: string;
  totalRecords: number;
  validRecords: number;
  duplicateRecords: number;
  transactions: ParsedTransaction[];
  accountSuggestions: Array<{
    name: string;
    suggestedAccountId: string | null;
  }>;
}

export interface ImportResult {
  jobId: string;
  success: boolean;
  importedCount: number;
  skippedCount: number;
  errorCount: number;
  errors: Array<{
    index: number;
    message: string;
  }>;
}

export type UploadAndParseInput = {
  source: ImportSource;
  fileName: string;
  fileBytes: ArrayBuffer;
};

export type ExecuteImportInput = {
  job_id: string;
  transactions: Array<{
    raw_index: number;
    selected_account_id: string | null;
    selected_category_id: string | null;
    note: string | null;
  }>;
};
