import 'package:uuid/uuid.dart';

import 'account.dart';
import 'category.dart';
import 'transaction.dart';

enum ImportSource { alipay, wechat, bank, generic }

enum ImportStatus { pending, parsing, preview, importing, completed, failed }

class ImportJob {
  final String id;
  final String userId;
  final ImportSource source;
  final String fileName;
  final int fileSize;
  final ImportStatus status;
  final int totalRows;
  final int importedRows;
  final String? errorMsg;
  final DateTime createdAt;
  final DateTime updatedAt;

  ImportJob({
    required this.id,
    required this.userId,
    required this.source,
    required this.fileName,
    required this.fileSize,
    required this.status,
    required this.totalRows,
    required this.importedRows,
    this.errorMsg,
    required this.createdAt,
    required this.updatedAt,
  });
}

class ParsedTransaction {
  final Map<String, String> rawData;
  final DateTime transactionDate;
  final TransactionType type;
  final double amount;
  final String currency;
  final String note;
  final String? accountName;
  final String? accountNumber;
  final String? counterparty;
  final String? categoryHint;
  final ImportSource source;
  final int lineNumber;
  final bool isDuplicate;
  final bool canBeImported;
  final String? importWarning;
  String? selectedAccountId;
  String? selectedCategoryId;

  ParsedTransaction({
    required this.rawData,
    required this.transactionDate,
    required this.type,
    required this.amount,
    required this.currency,
    required this.note,
    this.accountName,
    this.accountNumber,
    this.counterparty,
    this.categoryHint,
    required this.source,
    required this.lineNumber,
    required this.isDuplicate,
    required this.canBeImported,
    this.importWarning,
    this.selectedAccountId,
    this.selectedCategoryId,
  });

  ParsedTransaction copyWith({
    Map<String, String>? rawData,
    DateTime? transactionDate,
    TransactionType? type,
    double? amount,
    String? currency,
    String? note,
    String? accountName,
    String? accountNumber,
    String? counterparty,
    String? categoryHint,
    ImportSource? source,
    int? lineNumber,
    bool? isDuplicate,
    bool? canBeImported,
    String? importWarning,
    String? selectedAccountId,
    String? selectedCategoryId,
  }) {
    return ParsedTransaction(
      rawData: rawData ?? this.rawData,
      transactionDate: transactionDate ?? this.transactionDate,
      type: type ?? this.type,
      amount: amount ?? this.amount,
      currency: currency ?? this.currency,
      note: note ?? this.note,
      accountName: accountName ?? this.accountName,
      accountNumber: accountNumber ?? this.accountNumber,
      counterparty: counterparty ?? this.counterparty,
      categoryHint: categoryHint ?? this.categoryHint,
      source: source ?? this.source,
      lineNumber: lineNumber ?? this.lineNumber,
      isDuplicate: isDuplicate ?? this.isDuplicate,
      canBeImported: canBeImported ?? this.canBeImported,
      importWarning: importWarning ?? this.importWarning,
      selectedAccountId: selectedAccountId ?? this.selectedAccountId,
      selectedCategoryId: selectedCategoryId ?? this.selectedCategoryId,
    );
  }

  factory ParsedTransaction.fromJson(Map<String, dynamic> json) {
    TransactionType parseType(String type) {
      switch (type) {
        case 'income':
          return TransactionType.income;
        case 'expense':
          return TransactionType.expense;
        case 'transfer':
          return TransactionType.transfer;
        default:
          return TransactionType.expense;
      }
    }

    ImportSource parseSource(String source) {
      switch (source) {
        case 'alipay':
          return ImportSource.alipay;
        case 'wechat':
          return ImportSource.wechat;
        case 'bank':
          return ImportSource.bank;
        case 'generic':
        default:
          return ImportSource.generic;
      }
    }

    return ParsedTransaction(
      rawData: Map<String, String>.from(json['raw_data'] ?? {}),
      transactionDate: DateTime.parse(json['transaction_date']),
      type: parseType(json['type']),
      amount: json['amount']?.toDouble() ?? 0.0,
      currency: json['currency'] ?? 'CNY',
      note: json['note'] ?? '',
      accountName: json['account_name'],
      accountNumber: json['account_number'],
      counterparty: json['counterparty'],
      categoryHint: json['category_hint'],
      source: parseSource(json['source']),
      lineNumber: json['line_number'] ?? 0,
      isDuplicate: json['is_duplicate'] ?? false,
      canBeImported: json['can_be_imported'] ?? false,
      importWarning: json['import_warning'],
      selectedAccountId: json['selected_account_id'],
      selectedCategoryId: json['selected_category_id'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'raw_data': rawData,
      'transaction_date': transactionDate.toIso8601String(),
      'type': type.name,
      'amount': amount,
      'currency': currency,
      'note': note,
      'account_name': accountName,
      'account_number': accountNumber,
      'counterparty': counterparty,
      'category_hint': categoryHint,
      'source': source.name,
      'line_number': lineNumber,
      'is_duplicate': isDuplicate,
      'can_be_imported': canBeImported,
      'import_warning': importWarning,
      'selected_account_id': selectedAccountId,
      'selected_category_id': selectedCategoryId,
    };
  }
}

class ImportPreview {
  final String jobId;
  final ImportSource source;
  final int totalRows;
  final int validRows;
  final int duplicateRows;
  final List<ParsedTransaction> transactions;
  final Map<String, List<Account>>? accountSuggestions;
  final List<Category> categories;

  ImportPreview({
    required this.jobId,
    required this.source,
    required this.totalRows,
    required this.validRows,
    required this.duplicateRows,
    required this.transactions,
    this.accountSuggestions,
    required this.categories,
  });

  factory ImportPreview.fromJson(Map<String, dynamic> json) {
    ImportSource parseSource(String source) {
      switch (source) {
        case 'alipay':
          return ImportSource.alipay;
        case 'wechat':
          return ImportSource.wechat;
        case 'bank':
          return ImportSource.bank;
        case 'generic':
        default:
          return ImportSource.generic;
      }
    }

    var transactionsJson = json['transactions'] as List? ?? [];
    var categoriesJson = json['categories'] as List? ?? [];
    var accountSuggestionsJson = json['account_suggestions'] as Map<String, dynamic>?;

    Map<String, List<Account>>? accountSuggestions;
    if (accountSuggestionsJson != null) {
      accountSuggestions = {};
      accountSuggestionsJson.forEach((key, value) {
        var list = value as List;
        accountSuggestions![key] = list.map((e) => Account.fromJson(e)).toList();
      });
    }

    return ImportPreview(
      jobId: json['job_id'],
      source: parseSource(json['source']),
      totalRows: json['total_rows'] ?? 0,
      validRows: json['valid_rows'] ?? 0,
      duplicateRows: json['duplicate_rows'] ?? 0,
      transactions: transactionsJson.map((e) => ParsedTransaction.fromJson(e)).toList(),
      accountSuggestions: accountSuggestions,
      categories: categoriesJson.map((e) => Category.fromJson(e)).toList(),
    );
  }
}

class ImportResult {
  final String jobId;
  final int totalRows;
  final int importedRows;
  final int skippedRows;
  final int failedRows;
  final List<String>? importedIds;
  final List<ImportError>? errors;

  ImportResult({
    required this.jobId,
    required this.totalRows,
    required this.importedRows,
    required this.skippedRows,
    required this.failedRows,
    this.importedIds,
    this.errors,
  });

  factory ImportResult.fromJson(Map<String, dynamic> json) {
    var importedIdsJson = json['imported_ids'] as List?;
    var errorsJson = json['errors'] as List?;

    return ImportResult(
      jobId: json['job_id'],
      totalRows: json['total_rows'] ?? 0,
      importedRows: json['imported_rows'] ?? 0,
      skippedRows: json['skipped_rows'] ?? 0,
      failedRows: json['failed_rows'] ?? 0,
      importedIds: importedIdsJson?.map((e) => e.toString()).toList(),
      errors: errorsJson?.map((e) => ImportError.fromJson(e)).toList(),
    );
  }
}

class ImportError {
  final int lineNumber;
  final String error;

  ImportError({
    required this.lineNumber,
    required this.error,
  });

  factory ImportError.fromJson(Map<String, dynamic> json) {
    return ImportError(
      lineNumber: json['line_number'] ?? 0,
      error: json['error'] ?? '',
    );
  }
}

class ImportSourceInfo {
  final String id;
  final String name;
  final String description;
  final String icon;

  ImportSourceInfo({
    required this.id,
    required this.name,
    required this.description,
    required this.icon,
  });

  factory ImportSourceInfo.fromJson(Map<String, dynamic> json) {
    return ImportSourceInfo(
      id: json['id'],
      name: json['name'],
      description: json['description'],
      icon: json['icon'],
    );
  }
}
