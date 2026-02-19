import 'package:drift/drift.dart';
import 'package:drift/web.dart';
import 'package:flutter/foundation.dart';

LazyDatabase openConnection() {
  return LazyDatabase(() async {
    if (kDebugMode) {
      print('Using Web database: account.db');
    }
    return WebDatabase(
      'account.db',
      logStatements: kDebugMode,
    );
  });
}
