import 'dart:async';

import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../sync/sync_manager.dart' as sync_manager;

part 'sync_event.dart';
part 'sync_state.dart';

class SyncBloc extends Bloc<SyncEvent, SyncState> {
  final sync_manager.SyncManager _syncManager;
  StreamSubscription? _statusSubscription;
  StreamSubscription? _messageSubscription;

  SyncBloc({
    required sync_manager.SyncManager syncManager,
  })  : _syncManager = syncManager,
        super(const SyncState()) {
    on<SyncStarted>(_onSyncStarted);
    on<SyncRequested>(_onSyncRequested);
    on<SyncStatusUpdated>(_onSyncStatusUpdated);
    on<SyncCompleted>(_onSyncCompleted);
    on<SyncFailed>(_onSyncFailed);

    _subscribeToSyncManager();
  }

  void _subscribeToSyncManager() {
    _statusSubscription = _syncManager.statusStream.listen((status) {
      switch (status) {
        case sync_manager.SyncStatus.idle:
          break;
        case sync_manager.SyncStatus.syncing:
          add(const SyncStarted());
          break;
        case sync_manager.SyncStatus.success:
          add(SyncCompleted(_syncManager.lastSyncAt ?? DateTime.now()));
          break;
        case sync_manager.SyncStatus.error:
          break;
      }
    });

    _messageSubscription = _syncManager.messageStream.listen((message) {
      add(SyncStatusUpdated(message));
    });
  }

  Future<void> _onSyncStarted(
    SyncStarted event,
    Emitter<SyncState> emit,
  ) async {
    emit(state.copyWith(
      status: SyncStatus.syncing,
      errorMessage: null,
    ));
  }

  Future<void> _onSyncRequested(
    SyncRequested event,
    Emitter<SyncState> emit,
  ) async {
    if (_syncManager.isSyncing) return;

    emit(state.copyWith(
      status: SyncStatus.syncing,
      errorMessage: null,
    ));

    try {
      await _syncManager.sync();
    } catch (e) {
      add(SyncFailed(e.toString()));
    }
  }

  Future<void> _onSyncStatusUpdated(
    SyncStatusUpdated event,
    Emitter<SyncState> emit,
  ) async {
    emit(state.copyWith(
      message: event.message,
      progress: event.progress,
    ));
  }

  Future<void> _onSyncCompleted(
    SyncCompleted event,
    Emitter<SyncState> emit,
  ) async {
    emit(state.copyWith(
      status: SyncStatus.success,
      lastSyncAt: event.lastSyncAt,
      message: 'Sync completed',
    ));

    // Reset to idle after a delay
    await Future.delayed(const Duration(seconds: 2));
    if (state.isSuccess) {
      emit(state.copyWith(status: SyncStatus.idle));
    }
  }

  Future<void> _onSyncFailed(
    SyncFailed event,
    Emitter<SyncState> emit,
  ) async {
    emit(state.copyWith(
      status: SyncStatus.error,
      errorMessage: event.error,
    ));
  }

  @override
  Future<void> close() {
    _statusSubscription?.cancel();
    _messageSubscription?.cancel();
    return super.close();
  }
}
