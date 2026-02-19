import 'package:equatable/equatable.dart';
import 'package:json_annotation/json_annotation.dart';

import '../../sync/lww_strategy.dart';

part 'transaction.g.dart';

enum TransactionType {
  @JsonValue('income')
  income,
  @JsonValue('expense')
  expense,
  @JsonValue('transfer')
  transfer,
}

@JsonSerializable()
class Transaction extends Equatable implements LwwEntity {
  final String id;
  @JsonKey(name: 'user_id')
  final String userId;
  @JsonKey(name: 'account_id')
  final String accountId;
  @JsonKey(name: 'category_id')
  final String? categoryId;
  final TransactionType type;
  final double amount;
  final String currency;
  final String? note;
  @JsonKey(name: 'transaction_date')
  final DateTime transactionDate;
  @JsonKey(name: 'created_at')
  final DateTime createdAt;
  @JsonKey(name: 'updated_at')
  final DateTime updatedAt;
  @JsonKey(name: 'last_modified_at')
  final DateTime lastModifiedAt;
  final int version;
  @JsonKey(name: 'is_deleted')
  final bool isDeleted;

  const Transaction({
    required this.id,
    required this.userId,
    required this.accountId,
    this.categoryId,
    required this.type,
    required this.amount,
    required this.currency,
    this.note,
    required this.transactionDate,
    required this.createdAt,
    required this.updatedAt,
    required this.lastModifiedAt,
    required this.version,
    required this.isDeleted,
  });

  Transaction copyWith({
    String? id,
    String? userId,
    String? accountId,
    String? categoryId,
    TransactionType? type,
    double? amount,
    String? currency,
    String? note,
    DateTime? transactionDate,
    DateTime? createdAt,
    DateTime? updatedAt,
    DateTime? lastModifiedAt,
    int? version,
    bool? isDeleted,
  }) {
    return Transaction(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      accountId: accountId ?? this.accountId,
      categoryId: categoryId ?? this.categoryId,
      type: type ?? this.type,
      amount: amount ?? this.amount,
      currency: currency ?? this.currency,
      note: note ?? this.note,
      transactionDate: transactionDate ?? this.transactionDate,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
      lastModifiedAt: lastModifiedAt ?? this.lastModifiedAt,
      version: version ?? this.version,
      isDeleted: isDeleted ?? this.isDeleted,
    );
  }

  factory Transaction.fromJson(Map<String, dynamic> json) => _$TransactionFromJson(json);

  Map<String, dynamic> toJson() => _$TransactionToJson(this);

  @override
  List<Object?> get props => [
    id, userId, accountId, categoryId, type, amount, currency, note,
    transactionDate, createdAt, updatedAt, lastModifiedAt, version, isDeleted,
  ];
}
