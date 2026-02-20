import 'package:drift/drift.dart';
import 'package:account/data/datasources/local/database/app_database.dart';

QueryExecutor createTestQueryExecutor() {
  throw UnimplementedError('Cannot create a test database without dart:io or dart:html');
}

AppDatabase createTestDatabase() {
  throw UnimplementedError('Cannot create a test database without dart:io or dart:html');
}
