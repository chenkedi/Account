import 'package:flutter_test/flutter_test.dart';
import 'package:account/data/datasources/local/database/app_database.dart';
import 'test_database_connection.dart' as conn;

/// Creates a platform-appropriate test database.
///
/// - On native platforms: Uses in-memory SQLite database
/// - On web: Uses WebDatabase with unique name for test isolation
AppDatabase createTestDatabase() {
  return conn.createTestDatabase();
}
