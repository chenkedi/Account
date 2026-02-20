import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:account/data/repositories/transaction_repository.dart';
import 'package:account/data/datasources/local/database/app_database.dart';
import 'package:account/core/network/api_client.dart';
import '../../helpers/mocks.dart';

void main() {
  group('TransactionRepository', () {
    late MockTransactionDao mockTransactionDao;
    late MockApiClient mockApiClient;
    late TransactionRepository transactionRepository;

    setUp(() {
      mockTransactionDao = MockTransactionDao();
      mockApiClient = MockApiClient();
      transactionRepository = TransactionRepository(
        mockTransactionDao,
        mockApiClient,
      );
    });

    test('watchAllTransactions delegates to dao', () {
      // Create empty list for mocking - we just need to verify the interaction
      when(() => mockTransactionDao.watchAllTransactions(limit: 100))
          .thenAnswer((_) => Stream.value([]));

      final stream = transactionRepository.watchAllTransactions();

      expect(stream, emits(isEmpty));
      verify(() => mockTransactionDao.watchAllTransactions(limit: 100)).called(1);
    });

    test('getAllTransactions delegates to dao', () async {
      // Create empty list for mocking - we just need to verify the interaction
      when(() => mockTransactionDao.getAllTransactions(limit: 100, offset: 0))
          .thenAnswer((_) async => []);

      final result = await transactionRepository.getAllTransactions();

      expect(result, isEmpty);
      verify(() => mockTransactionDao.getAllTransactions(limit: 100, offset: 0)).called(1);
    });

    test('getTransactionById delegates to dao', () async {
      const id = 'test-id';
      // Return null for mocking - we just need to verify the interaction
      when(() => mockTransactionDao.getTransactionById(id))
          .thenAnswer((_) async => null);

      final result = await transactionRepository.getTransactionById(id);

      expect(result, isNull);
      verify(() => mockTransactionDao.getTransactionById(id)).called(1);
    });
  });
}
