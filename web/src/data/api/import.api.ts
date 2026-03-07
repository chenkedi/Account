import { ApiClient } from './client';
import type {
  ImportSource,
  ImportPreview,
  ImportResult,
  ImportPreviewServer,
  ImportResultServer,
} from '../models/import.model';
import { API_CONSTANTS } from '../../core/constants/api.constants';

export interface ImportSourceResponse {
  sources: Array<{
    id: ImportSource;
    name: string;
    description: string;
    icon: string;
  }>;
}

export interface TemplateInfoResponse {
  source: ImportSource;
  description: string;
  required_columns: string[];
  optional_columns: string[];
  file_extensions: string[];
}

export class ImportApi {
  constructor(private apiClient: ApiClient) {}

  async getSupportedSources(): Promise<ImportSourceResponse> {
    return this.apiClient.get<ImportSourceResponse>(API_CONSTANTS.importSources);
  }

  async getTemplateInfo(source: ImportSource): Promise<TemplateInfoResponse> {
    return this.apiClient.get<TemplateInfoResponse>(API_CONSTANTS.importTemplate, {
      params: { source },
    });
  }

  async uploadAndParse(source: ImportSource, file: File): Promise<ImportPreview> {
    const serverData = await this.apiClient.upload<ImportPreviewServer>(API_CONSTANTS.importUpload, file, {
      source,
    });
    // Transform snake_case to camelCase for UI
    return {
      jobId: serverData.job_id,
      source: serverData.source,
      fileName: serverData.file_name,
      totalRecords: serverData.total_rows,
      validRecords: serverData.valid_rows,
      duplicateRecords: serverData.duplicate_rows,
      transactions: serverData.transactions.map((t: any) => ({
        rawIndex: t.line_number,
        date: t.transaction_date,
        amount: t.amount,
        type: t.type,
        description: t.note,
        counterparty: t.counterparty,
        accountSuggestion: t.account_name,
        categorySuggestion: t.category_hint,
        isDuplicate: t.is_duplicate,
        canBeImported: t.can_be_imported,
      })),
      accountSuggestions: [],
    };
  }

  async executeImport(
    jobId: string,
    transactions: Array<{
      raw_index: number;
      selected_account_id: string | null;
      selected_category_id: string | null;
      note: string | null;
    }>
  ): Promise<ImportResult> {
    const serverData = await this.apiClient.post<ImportResultServer>(API_CONSTANTS.importExecute, {
      job_id: jobId,
      transactions,
    });
    // Transform snake_case to camelCase for UI
    return {
      jobId: serverData.job_id,
      success: serverData.failed_rows === 0,
      importedCount: serverData.imported_rows,
      skippedCount: serverData.skipped_rows,
      errorCount: serverData.failed_rows,
      errors: (serverData.errors || []).map((e: any) => ({
        index: e.line_number,
        message: e.error,
      })),
    };
  }
}
