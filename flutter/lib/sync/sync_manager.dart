import 'dart:async';

import '../core/network/websocket_client.dart';
import '../data/datasources/local/database/daos/account_dao.dart';
import '../data/datasources/local/database/daos/category_dao.dart';
import '../data/datasources/local/database/daos/transaction_dao.dart';
import '../data/datasources/local/shared_prefs/app_shared_prefs.dart';
import '../data/datasources/remote/api_service.dart';
import '../data/models/account.dart';
import '../data/models/category.dart';
import '../data/models/transaction.dart';
import 'lww_strategy.dart';

enum SyncStatus { idle, syncing, success, error }

class SyncManager {
  final ApiService _apiService;
  final AppSharedPrefs _sharedPrefs;
  final WebSocketClient? _websocketClient;
  final AccountDao _accountDao;
  final CategoryDao _categoryDao;
  final TransactionDao _transactionDao;
  final LwwStrategy _lwwStrategy;

  final _statusController = StreamController<SyncStatus>.broadcast();
  final _messageController = StreamController<String>.broadcast();
  DateTime? _lastSyncAt;
  String? _deviceId;
  StreamSubscription? _wsSubscription;

  SyncStatus _status = SyncStatus.idle;
  bool _isSyncing = false;

  SyncManager(
    this._apiService,
    this._sharedPrefs,
    this._accountDao,
    this._categoryDao,
    this._transactionDao, [
    this._websocketClient,
  ]) : _lwwStrategy = LwwStrategy();

  Stream<SyncStatus> get statusStream => _statusController.stream;
  Stream<String> get messageStream => _messageController.stream;
  SyncStatus get status => _status;
  bool get isSyncing => _isSyncing;
  DateTime? get lastSyncAt => _lastSyncAt;

  Future<void> initialize() async {
    _deviceId = await _sharedPrefs.getOrCreateDeviceId();
    _lastSyncAt = _sharedPrefs.lastSyncAt;

    // If we have an auth token, connect to WebSocket
    final token = _sharedPrefs.authToken;
    if (token != null && _deviceId != null && _websocketClient != null) {
      _websocketClient!.setCredentials(token, _deviceId!);
      _websocketClient!.connect();
      _subscribeToWebSocket();
    }
  }

  void _subscribeToWebSocket() {
    if (_websocketClient == null) return;

    _wsSubscription = _websocketClient!.messageStream.listen((message) {
      final type = message['type'] as String?;
      if (type == 'sync_available') {
        // Trigger a sync when we get a notification
        if (!_isSyncing) {
          sync();
        }
      }
    });
  }

  Future<void> sync() async {
    if (_isSyncing) return;

    _isSyncing = true;
    _updateStatus(SyncStatus.syncing);
    _addMessage('Starting sync...');

    try {
      // Step 1: Pull changes from server
      await _pullChanges();

      // Step 2: Push local changes to server
      await _pushChanges();

      _lastSyncAt = DateTime.now().toUtc();
      await _sharedPrefs.saveLastSyncAt(_lastSyncAt!);
      _updateStatus(SyncStatus.success);
      _addMessage('Sync completed successfully');
    } catch (e) {
      _updateStatus(SyncStatus.error);
      _addMessage('Sync failed: $e');
    } finally {
      _isSyncing = false;
    }
  }

  Future<void> _pullChanges() async {
    _addMessage('Pulling changes from server...');

    // Get last sync time
    final since = _lastSyncAt ?? DateTime(2020);

    // TODO: Implement actual API calls to pull changes
    // For now, just simulate
    await Future.delayed(const Duration(milliseconds: 500));
  }

  Future<void> _pushChanges() async {
    _addMessage('Pushing changes to server...');

    // Get all modified entities from local DB
    final since = _lastSyncAt ?? DateTime(2020);

    final modifiedAccounts = await _accountDao.getModifiedSince(since);
    final modifiedCategories = await _categoryDao.getModifiedSince(since);
    final modifiedTransactions = await _transactionDao.getModifiedSince(since);

    // TODO: Implement actual API calls to push changes
    // For now, just simulate
    await Future.delayed(const Duration(milliseconds: 500));

    _addMessage('Pushed ${modifiedAccounts.length} accounts, ${modifiedCategories.length} categories, ${modifiedTransactions.length} transactions');
  }

  void _updateStatus(SyncStatus status) {
    _status = status;
    _statusController.add(status);
  }

  void _addMessage(String message) {
    _messageController.add(message);
  }

  Future<void> dispose() async {
    _wsSubscription?.cancel();
    _websocketClient?.dispose();
    _statusController.close();
    _messageController.close();
  }
}
