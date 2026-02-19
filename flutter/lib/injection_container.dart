import 'package:get_it/get_it.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'core/network/api_client.dart';
import 'core/network/websocket_client.dart';
import 'data/datasources/local/database/app_database.dart';
import 'data/datasources/local/database/daos/account_dao.dart';
import 'data/datasources/local/database/daos/category_dao.dart';
import 'data/datasources/local/database/daos/transaction_dao.dart';
import 'data/datasources/local/shared_prefs/app_shared_prefs.dart';
import 'data/datasources/remote/api_service.dart';
import 'data/repositories/account_repository.dart';
import 'data/repositories/auth_repository.dart';
import 'data/repositories/category_repository.dart';
import 'data/repositories/transaction_repository.dart';
import 'presentation/import/bloc/import_bloc.dart';
import 'presentation/stats/bloc/stats_bloc.dart';
import 'sync/sync_manager.dart';

final sl = GetIt.instance;

Future<void> initDependencies() async {
  // Shared Preferences
  final sharedPrefs = await SharedPreferences.getInstance();
  sl.registerSingleton<AppSharedPrefs>(AppSharedPrefs(sharedPrefs));

  // Database
  final db = AppDatabase();
  sl.registerSingleton<AppDatabase>(db);

  // DAOs
  sl.registerSingleton<AccountDao>(AccountDao(db));
  sl.registerSingleton<CategoryDao>(CategoryDao(db));
  sl.registerSingleton<TransactionDao>(TransactionDao(db));

  // API Client & Service
  sl.registerSingleton<ApiClient>(ApiClient());
  sl.registerSingleton<ApiService>(ApiService(sl()));

  // WebSocket Client
  sl.registerSingleton<WebSocketClient>(WebSocketClient());

  // Repositories
  sl.registerSingleton<AuthRepository>(AuthRepository(sl(), sl()));
  sl.registerSingleton<AccountRepository>(AccountRepository(sl(), sl()));
  sl.registerSingleton<CategoryRepository>(CategoryRepository(sl(), sl()));
  sl.registerSingleton<TransactionRepository>(TransactionRepository(sl(), sl()));

  // Sync Manager
  sl.registerSingleton<SyncManager>(SyncManager(
    sl(),
    sl(),
    sl(),
    sl(),
    sl(),
    sl(),
    sl(),
  ));

  // BLoCs (Factories)
  sl.registerFactory<ImportBloc>(() => ImportBloc());
  sl.registerFactory<StatsBloc>(() => StatsBloc());
}
