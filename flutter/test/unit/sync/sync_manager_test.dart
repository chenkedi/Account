import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:account/sync/sync_manager.dart';
import 'package:account/data/datasources/remote/api_service.dart';
import 'package:account/data/datasources/local/shared_prefs/app_shared_prefs.dart';
import 'package:account/data/datasources/local/database/daos/account_dao.dart';
import 'package:account/data/datasources/local/database/daos/category_dao.dart';
import 'package:account/data/datasources/local/database/daos/transaction_dao.dart';
import 'package:account/core/network/websocket_client.dart';
import '../../helpers/mocks.dart';

void main() {
  group('SyncManager', () {
    late MockApiService mockApiService;
    late MockAppSharedPrefs mockSharedPrefs;
    late MockAccountDao mockAccountDao;
    late MockCategoryDao mockCategoryDao;
    late MockTransactionDao mockTransactionDao;
    late MockWebSocketClient mockWebSocketClient;

    setUp(() {
      mockApiService = MockApiService();
      mockSharedPrefs = MockAppSharedPrefs();
      mockAccountDao = MockAccountDao();
      mockCategoryDao = MockCategoryDao();
      mockTransactionDao = MockTransactionDao();
      mockWebSocketClient = MockWebSocketClient();
    });

    test('initialize sets up device id and last sync time', () async {
      const deviceId = 'test-device-123';
      final lastSyncAt = DateTime(2024, 1, 1).toUtc();

      when(() => mockSharedPrefs.getOrCreateDeviceId())
          .thenAnswer((_) async => deviceId);
      when(() => mockSharedPrefs.lastSyncAt).thenReturn(lastSyncAt);
      when(() => mockSharedPrefs.authToken).thenReturn(null);

      final syncManager = SyncManager(
        mockApiService,
        mockSharedPrefs,
        mockAccountDao,
        mockCategoryDao,
        mockTransactionDao,
        mockWebSocketClient,
      );

      await syncManager.initialize();

      verify(() => mockSharedPrefs.getOrCreateDeviceId()).called(1);
    });

    test('status stream emits syncing and success during sync', () async {
      when(() => mockSharedPrefs.getOrCreateDeviceId())
          .thenAnswer((_) async => 'test-device');
      when(() => mockSharedPrefs.lastSyncAt).thenReturn(null);
      when(() => mockSharedPrefs.authToken).thenReturn(null);
      when(() => mockSharedPrefs.saveLastSyncAt(any()))
          .thenAnswer((_) async {});
      when(() => mockAccountDao.getModifiedSince(any()))
          .thenAnswer((_) async => []);
      when(() => mockCategoryDao.getModifiedSince(any()))
          .thenAnswer((_) async => []);
      when(() => mockTransactionDao.getModifiedSince(any()))
          .thenAnswer((_) async => []);

      final syncManager = SyncManager(
        mockApiService,
        mockSharedPrefs,
        mockAccountDao,
        mockCategoryDao,
        mockTransactionDao,
        mockWebSocketClient,
      );

      await syncManager.initialize();

      expectLater(
        syncManager.statusStream,
        emitsInOrder([
          SyncStatus.syncing,
          SyncStatus.success,
        ]),
      );

      await syncManager.sync();
    });
  });
}
