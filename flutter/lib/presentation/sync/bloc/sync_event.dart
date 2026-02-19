part of 'sync_bloc.dart';

sealed class SyncEvent extends Equatable {
  const SyncEvent();

  @override
  List<Object?> get props => [];
}

final class SyncStarted extends SyncEvent {
  const SyncStarted();
}

final class SyncRequested extends SyncEvent {
  const SyncRequested();
}

final class SyncStatusUpdated extends SyncEvent {
  final String message;
  final int? progress;

  const SyncStatusUpdated(this.message, [this.progress]);

  @override
  List<Object?> get props => [message, progress];
}

final class SyncCompleted extends SyncEvent {
  final DateTime lastSyncAt;

  const SyncCompleted(this.lastSyncAt);

  @override
  List<Object?> get props => [lastSyncAt];
}

final class SyncFailed extends SyncEvent {
  final String error;

  const SyncFailed(this.error);

  @override
  List<Object?> get props => [error];
}
