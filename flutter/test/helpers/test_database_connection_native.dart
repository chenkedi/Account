import 'package:drift/drift.dart';
import 'package:drift/native.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:account/data/datasources/local/database/app_database.dart';

QueryExecutor createTestQueryExecutor() {
  return NativeDatabase.memory();
}

AppDatabase createTestDatabase() {
  final db = AppDatabase.forTesting(createTestQueryExecutor());
  addTearDown(() async => db.close());
  return db;
}
