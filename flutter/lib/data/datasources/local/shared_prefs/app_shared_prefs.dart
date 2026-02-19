import 'package:shared_preferences/shared_preferences.dart';
import 'package:uuid/uuid.dart';

class AppSharedPrefs {
  static const String _keyAuthToken = 'auth_token';
  static const String _keyUserId = 'user_id';
  static const String _keyUserEmail = 'user_email';
  static const String _keyDeviceId = 'device_id';
  static const String _keyLastSyncAt = 'last_sync_at';
  static const String _keyOnboardingComplete = 'onboarding_complete';

  static const Uuid _uuid = Uuid();

  final SharedPreferences _prefs;

  AppSharedPrefs(this._prefs);

  // Auth
  Future<void> saveAuthToken(String token) => _prefs.setString(_keyAuthToken, token);
  String? get authToken => _prefs.getString(_keyAuthToken);
  Future<void> clearAuthToken() => _prefs.remove(_keyAuthToken);

  // User
  Future<void> saveUserId(String id) => _prefs.setString(_keyUserId, id);
  String? get userId => _prefs.getString(_keyUserId);
  Future<void> saveUserEmail(String email) => _prefs.setString(_keyUserEmail, email);
  String? get userEmail => _prefs.getString(_keyUserEmail);
  Future<void> clearUser() async {
    await _prefs.remove(_keyUserId);
    await _prefs.remove(_keyUserEmail);
  }

  // Sync
  Future<void> saveDeviceId(String id) => _prefs.setString(_keyDeviceId, id);
  String? get deviceId => _prefs.getString(_keyDeviceId);
  Future<String> getOrCreateDeviceId() async {
    var id = _prefs.getString(_keyDeviceId);
    if (id == null) {
      id = _uuid.v4();
      await _prefs.setString(_keyDeviceId, id);
    }
    return id;
  }

  Future<void> saveLastSyncAt(DateTime date) =>
      _prefs.setString(_keyLastSyncAt, date.toIso8601String());
  DateTime? get lastSyncAt {
    final dateStr = _prefs.getString(_keyLastSyncAt);
    if (dateStr != null) {
      return DateTime.tryParse(dateStr);
    }
    return null;
  }

  // Onboarding
  Future<void> setOnboardingComplete() => _prefs.setBool(_keyOnboardingComplete, true);
  bool get isOnboardingComplete => _prefs.getBool(_keyOnboardingComplete) ?? false;

  // Clear all
  Future<void> clearAll() async {
    await _prefs.clear();
  }
}
