part of 'sync_bloc.dart';

// Use SyncStatus from sync_manager
typedef SyncStatus = sync_manager.SyncStatus;

class SyncState extends Equatable {
  final SyncStatus status;
  final String? message;
  final int? progress;
  final DateTime? lastSyncAt;
  final String? errorMessage;

  const SyncState({
    this.status = SyncStatus.idle,
    this.message,
    this.progress,
    this.lastSyncAt,
    this.errorMessage,
  });

  SyncState copyWith({
    SyncStatus? status,
    String? message,
    int? progress,
    DateTime? lastSyncAt,
    String? errorMessage,
  }) {
    return SyncState(
      status: status ?? this.status,
      message: message ?? this.message,
      progress: progress ?? this.progress,
      lastSyncAt: lastSyncAt ?? this.lastSyncAt,
      errorMessage: errorMessage ?? this.errorMessage,
    );
  }

  bool get isSyncing => status == SyncStatus.syncing;
  bool get hasError => status == SyncStatus.error;
  bool get isSuccess => status == SyncStatus.success;

  @override
  List<Object?> get props => [status, message, progress, lastSyncAt, errorMessage];
}
