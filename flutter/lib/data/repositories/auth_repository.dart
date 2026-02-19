import '../datasources/local/shared_prefs/app_shared_prefs.dart';
import '../datasources/remote/api_service.dart';
import '../datasources/remote/models/responses/auth_responses.dart';
import '../models/user.dart';

class AuthRepository {
  final ApiService _apiService;
  final AppSharedPrefs _sharedPrefs;

  AuthRepository(this._apiService, this._sharedPrefs);

  Future<AuthResponse> register(String email, String password) async {
    final response = await _apiService.register(email, password);
    await _saveAuthData(response);
    return response;
  }

  Future<AuthResponse> login(String email, String password) async {
    final response = await _apiService.login(email, password);
    await _saveAuthData(response);
    return response;
  }

  Future<void> logout() async {
    await _sharedPrefs.clearAuthToken();
    await _sharedPrefs.clearUser();
  }

  Future<void> _saveAuthData(AuthResponse response) async {
    await _sharedPrefs.saveAuthToken(response.accessToken);
    await _sharedPrefs.saveUserId(response.user.id);
    await _sharedPrefs.saveUserEmail(response.user.email);
  }

  Future<bool> isAuthenticated() async {
    final token = _sharedPrefs.authToken;
    return token != null && token.isNotEmpty;
  }

  User? get currentUser {
    final userId = _sharedPrefs.userId;
    final userEmail = _sharedPrefs.userEmail;
    if (userId != null && userEmail != null) {
      return User(
        id: userId,
        email: userEmail,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );
    }
    return null;
  }

  String? get authToken => _sharedPrefs.authToken;

  bool get hasToken => _sharedPrefs.authToken != null;
}
