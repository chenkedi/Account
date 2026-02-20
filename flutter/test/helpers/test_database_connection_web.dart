import 'package:drift/drift.dart';
import 'package:drift/web.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:account/data/datasources/local/database/app_database.dart';

/// Creates a unique database name for web testing to ensure isolation
String _createUniqueDbName() {
  final timestamp = DateTime.now().millisecondsSinceEpoch;
  return 'test_db_$timestamp';
}

QueryExecutor createTestQueryExecutor() {
  return WebDatabase(_createUniqueDbName(), logStatements: false);
}

AppDatabase createTestDatabase() {
  final db = AppDatabase.forTesting(createTestQueryExecutor());
  addTearDown(() async {
    await db.close();
  });
  return db;
}
