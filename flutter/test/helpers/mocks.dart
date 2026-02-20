import 'package:mocktail/mocktail.dart';
import 'package:account/data/datasources/remote/api_service.dart';
import 'package:account/core/network/websocket_client.dart';
import 'package:account/data/datasources/local/shared_prefs/app_shared_prefs.dart';
import 'package:account/sync/sync_manager.dart';
import 'package:account/data/datasources/local/database/daos/transaction_dao.dart';
import 'package:account/data/datasources/local/database/daos/account_dao.dart';
import 'package:account/data/datasources/local/database/daos/category_dao.dart';
import 'package:account/core/network/api_client.dart';
import 'package:account/data/repositories/auth_repository.dart';

class MockApiService extends Mock implements ApiService {}
class MockWebSocketClient extends Mock implements WebSocketClient {}
class MockAppSharedPrefs extends Mock implements AppSharedPrefs {}
class MockSyncManager extends Mock implements SyncManager {}
class MockTransactionDao extends Mock implements TransactionDao {}
class MockAccountDao extends Mock implements AccountDao {}
class MockCategoryDao extends Mock implements CategoryDao {}
class MockApiClient extends Mock implements ApiClient {}
class MockAuthRepository extends Mock implements AuthRepository {}
