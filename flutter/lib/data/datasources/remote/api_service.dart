import '../../../../core/constants/api_constants.dart';
import '../../../../core/errors/exceptions.dart';
import '../../../../core/network/api_client.dart';
import '../../models/account.dart';
import '../../models/category.dart';
import '../../models/transaction.dart';
import '../../models/user.dart';
import '../../models/import.dart';
import 'models/requests/auth_requests.dart';
import 'models/responses/auth_responses.dart';

class ApiService {
  final ApiClient _apiClient;

  ApiService(this._apiClient);

  // Auth
  Future<AuthResponse> login(String email, String password) async {
    final response = await _apiClient.post(
      ApiConstants.login,
      data: LoginRequest(email: email, password: password).toJson(),
    );
    return AuthResponse.fromJson(response);
  }

  Future<AuthResponse> register(String email, String password) async {
    final response = await _apiClient.post(
      ApiConstants.register,
      data: RegisterRequest(email: email, password: password).toJson(),
    );
    return AuthResponse.fromJson(response);
  }

  // Accounts
  Future<List<Account>> getAccounts() async {
    final response = await _apiClient.get(ApiConstants.accounts);
    return (response as List)
        .map((json) => Account.fromJson(json as Map<String, dynamic>))
        .toList();
  }

  Future<Account> createAccount(Map<String, dynamic> data) async {
    final response = await _apiClient.post(
      ApiConstants.accounts,
      data: data,
    );
    return Account.fromJson(response);
  }

  Future<Account> updateAccount(String id, Map<String, dynamic> data) async {
    final response = await _apiClient.put(
      '${ApiConstants.accounts}/$id',
      data: data,
    );
    return Account.fromJson(response);
  }

  Future<void> deleteAccount(String id) async {
    await _apiClient.delete('${ApiConstants.accounts}/$id');
  }

  // Categories
  Future<List<Category>> getCategories() async {
    final response = await _apiClient.get(ApiConstants.categories);
    return (response as List)
        .map((json) => Category.fromJson(json as Map<String, dynamic>))
        .toList();
  }

  Future<List<Category>> getCategoriesByType(String type) async {
    final response = await _apiClient.get('${ApiConstants.categories}/type/$type');
    return (response as List)
        .map((json) => Category.fromJson(json as Map<String, dynamic>))
        .toList();
  }

  Future<Category> createCategory(Map<String, dynamic> data) async {
    final response = await _apiClient.post(
      ApiConstants.categories,
      data: data,
    );
    return Category.fromJson(response);
  }

  Future<Category> updateCategory(String id, Map<String, dynamic> data) async {
    final response = await _apiClient.put(
      '${ApiConstants.categories}/$id',
      data: data,
    );
    return Category.fromJson(response);
  }

  Future<void> deleteCategory(String id) async {
    await _apiClient.delete('${ApiConstants.categories}/$id');
  }

  // Transactions
  Future<List<Transaction>> getTransactions({int limit = 50, int offset = 0}) async {
    final response = await _apiClient.get(
      ApiConstants.transactions,
      queryParameters: {'limit': limit, 'offset': offset},
    );
    return (response as List)
        .map((json) => Transaction.fromJson(json as Map<String, dynamic>))
        .toList();
  }

  Future<List<Transaction>> getTransactionsByDateRange(DateTime start, DateTime end) async {
    final response = await _apiClient.get(
      '${ApiConstants.transactions}/range',
      queryParameters: {
        'start_date': start.toIso8601String(),
        'end_date': end.toIso8601String(),
      },
    );
    return (response as List)
        .map((json) => Transaction.fromJson(json as Map<String, dynamic>))
        .toList();
  }

  Future<Map<String, dynamic>> getStats(DateTime start, DateTime end) async {
    return await _apiClient.get(
      '${ApiConstants.transactions}/stats',
      queryParameters: {
        'start_date': start.toIso8601String(),
        'end_date': end.toIso8601String(),
      },
    );
  }

  Future<Transaction> createTransaction(Map<String, dynamic> data) async {
    final response = await _apiClient.post(
      ApiConstants.transactions,
      data: data,
    );
    return Transaction.fromJson(response);
  }

  Future<Transaction> updateTransaction(String id, Map<String, dynamic> data) async {
    final response = await _apiClient.put(
      '${ApiConstants.transactions}/$id',
      data: data,
    );
    return Transaction.fromJson(response);
  }

  Future<void> deleteTransaction(String id) async {
    await _apiClient.delete('${ApiConstants.transactions}/$id');
  }

  Future<Map<String, dynamic>> getDetailedStats(DateTime start, DateTime end) async {
    return await _apiClient.get(
      '${ApiConstants.transactions}/stats/detailed',
      queryParameters: {
        'start_date': start.toIso8601String(),
        'end_date': end.toIso8601String(),
      },
    );
  }

  // Import
  Future<List<ImportSourceInfo>> getImportSources() async {
    final response = await _apiClient.get('/api/v1/import/sources');
    final sources = response['sources'] as List;
    return sources.map((json) => ImportSourceInfo.fromJson(json)).toList();
  }

  Future<ImportPreview> uploadAndParseFile(
    ImportSource source,
    String fileName,
    List<int> fileBytes,
  ) async {
    final sourceString = source.name;

    final formData = FormData.fromMap({
      'source': sourceString,
      'file': MultipartFile.fromBytes(
        fileBytes,
        filename: fileName,
      ),
    });

    final response = await _apiClient.post(
      '/api/v1/import/upload',
      data: formData,
    );

    return ImportPreview.fromJson(response);
  }

  Future<ImportResult> executeImport(
    String jobId,
    List<ParsedTransaction> transactions,
  ) async {
    final response = await _apiClient.post(
      '/api/v1/import/execute',
      data: {
        'job_id': jobId,
        'transactions': transactions.map((t) => t.toJson()).toList(),
      },
    );

    return ImportResult.fromJson(response);
  }
}
