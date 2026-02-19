import 'package:equatable/equatable.dart';
import 'package:json_annotation/json_annotation.dart';

import '../../sync/lww_strategy.dart';

part 'account.g.dart';

enum AccountType {
  @JsonValue('bank')
  bank,
  @JsonValue('cash')
  cash,
  @JsonValue('alipay')
  alipay,
  @JsonValue('wechat')
  wechat,
  @JsonValue('credit')
  credit,
  @JsonValue('investment')
  investment,
  @JsonValue('other')
  other,
}

@JsonSerializable()
class Account extends Equatable implements LwwEntity {
  final String id;
  @JsonKey(name: 'user_id')
  final String userId;
  final String name;
  final AccountType type;
  final String currency;
  final double balance;
  @JsonKey(name: 'created_at')
  final DateTime createdAt;
  @JsonKey(name: 'updated_at')
  final DateTime updatedAt;
  @JsonKey(name: 'last_modified_at')
  final DateTime lastModifiedAt;
  final int version;
  @JsonKey(name: 'is_deleted')
  final bool isDeleted;

  const Account({
    required this.id,
    required this.userId,
    required this.name,
    required this.type,
    required this.currency,
    required this.balance,
    required this.createdAt,
    required this.updatedAt,
    required this.lastModifiedAt,
    required this.version,
    required this.isDeleted,
  });

  Account copyWith({
    String? id,
    String? userId,
    String? name,
    AccountType? type,
    String? currency,
    double? balance,
    DateTime? createdAt,
    DateTime? updatedAt,
    DateTime? lastModifiedAt,
    int? version,
    bool? isDeleted,
  }) {
    return Account(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      name: name ?? this.name,
      type: type ?? this.type,
      currency: currency ?? this.currency,
      balance: balance ?? this.balance,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
      lastModifiedAt: lastModifiedAt ?? this.lastModifiedAt,
      version: version ?? this.version,
      isDeleted: isDeleted ?? this.isDeleted,
    );
  }

  factory Account.fromJson(Map<String, dynamic> json) => _$AccountFromJson(json);

  Map<String, dynamic> toJson() => _$AccountToJson(this);

  @override
  List<Object?> get props => [
    id, userId, name, type, currency, balance, createdAt,
    updatedAt, lastModifiedAt, version, isDeleted,
  ];
}
