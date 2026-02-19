import 'package:equatable/equatable.dart';
import 'package:json_annotation/json_annotation.dart';

part 'sync_state.g.dart';

@JsonSerializable()
class SyncState extends Equatable {
  @JsonKey(name: 'user_id')
  final String userId;
  @JsonKey(name: 'device_id')
  final String deviceId;
  @JsonKey(name: 'last_sync_at')
  final DateTime lastSyncAt;
  @JsonKey(name: 'sync_token')
  final String? syncToken;
  @JsonKey(name: 'created_at')
  final DateTime createdAt;
  @JsonKey(name: 'updated_at')
  final DateTime updatedAt;

  const SyncState({
    required this.userId,
    required this.deviceId,
    required this.lastSyncAt,
    this.syncToken,
    required this.createdAt,
    required this.updatedAt,
  });

  SyncState copyWith({
    String? userId,
    String? deviceId,
    DateTime? lastSyncAt,
    String? syncToken,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return SyncState(
      userId: userId ?? this.userId,
      deviceId: deviceId ?? this.deviceId,
      lastSyncAt: lastSyncAt ?? this.lastSyncAt,
      syncToken: syncToken ?? this.syncToken,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  factory SyncState.fromJson(Map<String, dynamic> json) =>
      _$SyncStateFromJson(json);

  Map<String, dynamic> toJson() => _$SyncStateToJson(this);

  @override
  List<Object?> get props => [
    userId, deviceId, lastSyncAt, syncToken, createdAt, updatedAt,
  ];
}
