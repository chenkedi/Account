class ApiConstants {
  static const String baseUrl = 'http://localhost:8080';
  static const String apiPrefix = '/api/v1';

  // Auth endpoints
  static const String register = '$apiPrefix/auth/register';
  static const String login = '$apiPrefix/auth/login';

  // Account endpoints
  static const String accounts = '$apiPrefix/accounts';

  // Category endpoints
  static const String categories = '$apiPrefix/categories';

  // Transaction endpoints
  static const String transactions = '$apiPrefix/transactions';

  // Sync endpoints
  static const String syncPull = '$apiPrefix/sync/pull';
  static const String syncPush = '$apiPrefix/sync/push';

  // WebSocket
  static const String wsSync = '/ws/sync';

  // Timeout
  static const int connectTimeout = 15000;
  static const int receiveTimeout = 15000;
}
