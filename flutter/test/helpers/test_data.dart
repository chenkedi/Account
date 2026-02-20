import 'package:uuid/uuid.dart';
import 'package:account/data/models/transaction.dart';
import 'package:account/data/models/account.dart';
import 'package:account/data/models/category.dart';
import 'package:account/data/models/user.dart';
import 'package:account/data/datasources/remote/models/responses/auth_responses.dart';

const _uuid = Uuid();

class TestData {
  static String createId() => _uuid.v4();

  static User createUser({
    String? id,
    String email = 'test@example.com',
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    final now = DateTime.now().toUtc();
    return User(
      id: id ?? createId(),
      email: email,
      createdAt: createdAt ?? now,
      updatedAt: updatedAt ?? now,
    );
  }

  static AuthResponse createAuthResponse({
    String accessToken = 'test-access-token',
    String tokenType = 'Bearer',
    int expiresIn = 3600,
    User? user,
  }) {
    return AuthResponse(
      accessToken: accessToken,
      tokenType: tokenType,
      expiresIn: expiresIn,
      user: user ?? createUser(),
    );
  }

  static Account createAccount({
    String? id,
    String? userId,
    String name = 'Test Account',
    AccountType type = AccountType.bank,
    String currency = 'CNY',
    double balance = 1000.0,
    DateTime? createdAt,
    DateTime? updatedAt,
    DateTime? lastModifiedAt,
    int version = 1,
    bool isDeleted = false,
  }) {
    final now = DateTime.now().toUtc();
    return Account(
      id: id ?? createId(),
      userId: userId ?? createId(),
      name: name,
      type: type,
      currency: currency,
      balance: balance,
      createdAt: createdAt ?? now,
      updatedAt: updatedAt ?? now,
      lastModifiedAt: lastModifiedAt ?? now,
      version: version,
      isDeleted: isDeleted,
    );
  }

  static Category createCategory({
    String? id,
    String? userId,
    String name = 'Test Category',
    CategoryType type = CategoryType.expense,
    String? parentId,
    String? icon,
    DateTime? createdAt,
    DateTime? updatedAt,
    DateTime? lastModifiedAt,
    int version = 1,
    bool isDeleted = false,
  }) {
    final now = DateTime.now().toUtc();
    return Category(
      id: id ?? createId(),
      userId: userId ?? createId(),
      name: name,
      type: type,
      parentId: parentId,
      icon: icon,
      createdAt: createdAt ?? now,
      updatedAt: updatedAt ?? now,
      lastModifiedAt: lastModifiedAt ?? now,
      version: version,
      isDeleted: isDeleted,
    );
  }

  static Transaction createTransaction({
    String? id,
    String? userId,
    String? accountId,
    String? categoryId,
    TransactionType type = TransactionType.expense,
    double amount = 100.0,
    String currency = 'CNY',
    String? note = 'Test transaction',
    DateTime? transactionDate,
    DateTime? createdAt,
    DateTime? updatedAt,
    DateTime? lastModifiedAt,
    int version = 1,
    bool isDeleted = false,
  }) {
    final now = DateTime.now().toUtc();
    return Transaction(
      id: id ?? createId(),
      userId: userId ?? createId(),
      accountId: accountId ?? createId(),
      categoryId: categoryId,
      type: type,
      amount: amount,
      currency: currency,
      note: note,
      transactionDate: transactionDate ?? now,
      createdAt: createdAt ?? now,
      updatedAt: updatedAt ?? now,
      lastModifiedAt: lastModifiedAt ?? now,
      version: version,
      isDeleted: isDeleted,
    );
  }
}
