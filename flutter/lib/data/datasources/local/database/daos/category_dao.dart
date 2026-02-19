import 'package:drift/drift.dart';

import '../app_database.dart';
import '../tables/category_table.dart';

part 'category_dao.g.dart';

@DriftAccessor(tables: [Categories])
class CategoryDao extends DatabaseAccessor<AppDatabase> with _$CategoryDaoMixin {
  CategoryDao(AppDatabase db) : super(db);

  Future<List<Category>> getAllCategories() {
    return (select(categories)
          ..where((c) => c.isDeleted.equals(false))
          ..orderBy([(c) => OrderingTerm(expression: c.name)]))
        .get();
  }

  Stream<List<Category>> watchAllCategories() {
    return (select(categories)
          ..where((c) => c.isDeleted.equals(false))
          ..orderBy([(c) => OrderingTerm(expression: c.name)]))
        .watch();
  }

  Future<List<Category>> getCategoriesByType(String type) {
    return (select(categories)
          ..where((c) => c.isDeleted.equals(false) & c.type.equals(type))
          ..orderBy([(c) => OrderingTerm(expression: c.name)]))
        .get();
  }

  Stream<List<Category>> watchCategoriesByType(String type) {
    return (select(categories)
          ..where((c) => c.isDeleted.equals(false) & c.type.equals(type))
          ..orderBy([(c) => OrderingTerm(expression: c.name)]))
        .watch();
  }

  Future<Category?> getCategoryById(String id) {
    return (select(categories)..where((c) => c.id.equals(id))).getSingleOrNull();
  }

  Future<void> insertCategory(CategoriesCompanion category) {
    return into(categories).insert(category, mode: InsertMode.insertOrReplace);
  }

  Future<void> insertCategories(List<CategoriesCompanion> categoryList) {
    return batch((batch) {
      batch.insertAll(categories, categoryList, mode: InsertMode.insertOrReplace);
    });
  }

  Future<void> updateCategory(CategoriesCompanion category) {
    return (update(categories)..where((c) => c.id.equals(category.id.value))).write(category);
  }

  Future<void> deleteCategory(String id) {
    return (update(categories)..where((c) => c.id.equals(id)))
        .write(CategoriesCompanion(isDeleted: const Value(true)));
  }

  Future<List<Category>> getModifiedSince(DateTime since) {
    return (select(categories)
          ..where((c) => c.lastModifiedAt.isBiggerThanValue(since)))
        .get();
  }
}
