import 'package:equatable/equatable.dart';
import 'package:json_annotation/json_annotation.dart';

import '../../sync/lww_strategy.dart';

part 'category.g.dart';

enum CategoryType {
  @JsonValue('income')
  income,
  @JsonValue('expense')
  expense,
}

@JsonSerializable()
class Category extends Equatable implements LwwEntity {
  final String id;
  @JsonKey(name: 'user_id')
  final String userId;
  final String name;
  final CategoryType type;
  @JsonKey(name: 'parent_id')
  final String? parentId;
  final String? icon;
  @JsonKey(name: 'created_at')
  final DateTime createdAt;
  @JsonKey(name: 'updated_at')
  final DateTime updatedAt;
  @JsonKey(name: 'last_modified_at')
  final DateTime lastModifiedAt;
  final int version;
  @JsonKey(name: 'is_deleted')
  final bool isDeleted;

  const Category({
    required this.id,
    required this.userId,
    required this.name,
    required this.type,
    this.parentId,
    this.icon,
    required this.createdAt,
    required this.updatedAt,
    required this.lastModifiedAt,
    required this.version,
    required this.isDeleted,
  });

  Category copyWith({
    String? id,
    String? userId,
    String? name,
    CategoryType? type,
    String? parentId,
    String? icon,
    DateTime? createdAt,
    DateTime? updatedAt,
    DateTime? lastModifiedAt,
    int? version,
    bool? isDeleted,
  }) {
    return Category(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      name: name ?? this.name,
      type: type ?? this.type,
      parentId: parentId ?? this.parentId,
      icon: icon ?? this.icon,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
      lastModifiedAt: lastModifiedAt ?? this.lastModifiedAt,
      version: version ?? this.version,
      isDeleted: isDeleted ?? this.isDeleted,
    );
  }

  factory Category.fromJson(Map<String, dynamic> json) => _$CategoryFromJson(json);

  Map<String, dynamic> toJson() => _$CategoryToJson(this);

  @override
  List<Object?> get props => [
    id, userId, name, type, parentId, icon, createdAt,
    updatedAt, lastModifiedAt, version, isDeleted,
  ];
}
