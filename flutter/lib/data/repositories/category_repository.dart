import '../datasources/local/database/app_database.dart';
import '../datasources/local/database/daos/category_dao.dart';
import '../../core/network/api_client.dart';

class CategoryRepository {
  final CategoryDao _categoryDao;
  final ApiClient _apiClient;

  CategoryRepository(this._categoryDao, this._apiClient);

  Stream<List<Category>> watchAllCategories() => _categoryDao.watchAllCategories();

  Stream<List<Category>> watchCategoriesByType(String type) =>
      _categoryDao.watchCategoriesByType(type);

  Future<List<Category>> getAllCategories() => _categoryDao.getAllCategories();

  Future<List<Category>> getCategoriesByType(String type) =>
      _categoryDao.getCategoriesByType(type);

  Future<Category?> getCategoryById(String id) => _categoryDao.getCategoryById(id);

  Future<void> addCategory(CategoriesCompanion category) =>
      _categoryDao.insertCategory(category);

  Future<void> updateCategory(CategoriesCompanion category) =>
      _categoryDao.updateCategory(category);

  Future<void> deleteCategory(String id) => _categoryDao.deleteCategory(id);
}
