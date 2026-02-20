import 'package:bloc_test/bloc_test.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:account/presentation/auth/bloc/auth_bloc.dart';
import 'package:account/data/repositories/auth_repository.dart';
import 'package:account/sync/sync_manager.dart';
import '../../helpers/mocks.dart';
import '../../helpers/test_data.dart';

void main() {
  group('AuthBloc', () {
    late MockAuthRepository mockAuthRepository;
    late MockSyncManager mockSyncManager;

    setUp(() {
      mockAuthRepository = MockAuthRepository();
      mockSyncManager = MockSyncManager();
    });

    group('AuthCheckRequested', () {
      blocTest<AuthBloc, AuthState>(
        'emits authenticated when user is authenticated',
        build: () {
          when(() => mockAuthRepository.isAuthenticated())
              .thenAnswer((_) async => true);
          when(() => mockAuthRepository.currentUser)
              .thenReturn(TestData.createUser());
          when(() => mockSyncManager.initialize())
              .thenAnswer((_) async {});
          return AuthBloc(
            authRepository: mockAuthRepository,
            syncManager: mockSyncManager,
          );
        },
        act: (bloc) => bloc.add(const AuthCheckRequested()),
        expect: () => [
          const AuthState(status: AuthStatus.loading),
          predicate<AuthState>(
            (state) => state.status == AuthStatus.authenticated && state.user != null,
          ),
        ],
        verify: (_) {
          verify(() => mockSyncManager.initialize()).called(1);
        },
      );

      blocTest<AuthBloc, AuthState>(
        'emits unauthenticated when user is not authenticated',
        build: () {
          when(() => mockAuthRepository.isAuthenticated())
              .thenAnswer((_) async => false);
          return AuthBloc(
            authRepository: mockAuthRepository,
            syncManager: mockSyncManager,
          );
        },
        act: (bloc) => bloc.add(const AuthCheckRequested()),
        expect: () => [
          const AuthState(status: AuthStatus.loading),
          const AuthState(status: AuthStatus.unauthenticated),
        ],
      );
    });

    group('AuthLoginRequested', () {
      blocTest<AuthBloc, AuthState>(
        'emits authenticated on successful login',
        build: () {
          final user = TestData.createUser(email: 'test@example.com');
          final authResponse = TestData.createAuthResponse(user: user);
          when(() => mockAuthRepository.login('test@example.com', 'password'))
              .thenAnswer((_) async => authResponse);
          when(() => mockSyncManager.initialize())
              .thenAnswer((_) async {});
          return AuthBloc(
            authRepository: mockAuthRepository,
            syncManager: mockSyncManager,
          );
        },
        act: (bloc) => bloc.add(const AuthLoginRequested(
          email: 'test@example.com',
          password: 'password',
        )),
        expect: () => [
          const AuthState(status: AuthStatus.loading),
          predicate<AuthState>(
            (state) =>
                state.status == AuthStatus.authenticated &&
                state.user?.email == 'test@example.com',
          ),
        ],
      );

      blocTest<AuthBloc, AuthState>(
        'emits unauthenticated with error on login failure',
        build: () {
          when(() => mockAuthRepository.login('test@example.com', 'wrong'))
              .thenThrow(Exception('Invalid credentials'));
          return AuthBloc(
            authRepository: mockAuthRepository,
            syncManager: mockSyncManager,
          );
        },
        act: (bloc) => bloc.add(const AuthLoginRequested(
          email: 'test@example.com',
          password: 'wrong',
        )),
        expect: () => [
          const AuthState(status: AuthStatus.loading),
          predicate<AuthState>(
            (state) =>
                state.status == AuthStatus.unauthenticated &&
                state.errorMessage != null,
          ),
        ],
      );
    });

    group('AuthRegisterRequested', () {
      blocTest<AuthBloc, AuthState>(
        'emits authenticated on successful register',
        build: () {
          final user = TestData.createUser(email: 'new@example.com');
          final authResponse = TestData.createAuthResponse(user: user);
          when(() => mockAuthRepository.register('new@example.com', 'password'))
              .thenAnswer((_) async => authResponse);
          when(() => mockSyncManager.initialize())
              .thenAnswer((_) async {});
          return AuthBloc(
            authRepository: mockAuthRepository,
            syncManager: mockSyncManager,
          );
        },
        act: (bloc) => bloc.add(const AuthRegisterRequested(
          email: 'new@example.com',
          password: 'password',
        )),
        expect: () => [
          const AuthState(status: AuthStatus.loading),
          predicate<AuthState>(
            (state) =>
                state.status == AuthStatus.authenticated &&
                state.user?.email == 'new@example.com',
          ),
        ],
      );

      blocTest<AuthBloc, AuthState>(
        'emits unauthenticated with error on register failure',
        build: () {
          when(() => mockAuthRepository.register('existing@example.com', 'password'))
              .thenThrow(Exception('Email already exists'));
          return AuthBloc(
            authRepository: mockAuthRepository,
            syncManager: mockSyncManager,
          );
        },
        act: (bloc) => bloc.add(const AuthRegisterRequested(
          email: 'existing@example.com',
          password: 'password',
        )),
        expect: () => [
          const AuthState(status: AuthStatus.loading),
          predicate<AuthState>(
            (state) =>
                state.status == AuthStatus.unauthenticated &&
                state.errorMessage != null,
          ),
        ],
      );
    });

    group('AuthLogoutRequested', () {
      blocTest<AuthBloc, AuthState>(
        'emits unauthenticated on logout',
        build: () {
          when(() => mockAuthRepository.logout())
              .thenAnswer((_) async {});
          when(() => mockSyncManager.dispose())
              .thenAnswer((_) async {});
          return AuthBloc(
            authRepository: mockAuthRepository,
            syncManager: mockSyncManager,
          );
        },
        act: (bloc) => bloc.add(const AuthLogoutRequested()),
        expect: () => [
          const AuthState(status: AuthStatus.loading),
          const AuthState(status: AuthStatus.unauthenticated),
        ],
        verify: (_) {
          verify(() => mockAuthRepository.logout()).called(1);
          verify(() => mockSyncManager.dispose()).called(1);
        },
      );
    });
  });
}
