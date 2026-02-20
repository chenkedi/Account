import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:account/data/repositories/auth_repository.dart';
import 'package:account/data/datasources/remote/api_service.dart';
import 'package:account/data/datasources/local/shared_prefs/app_shared_prefs.dart';
import 'package:account/data/datasources/remote/models/responses/auth_responses.dart';
import 'package:account/data/models/user.dart';
import '../../helpers/mocks.dart';
import '../../helpers/test_data.dart';

void main() {
  group('AuthRepository', () {
    late MockApiService mockApiService;
    late MockAppSharedPrefs mockSharedPrefs;
    late AuthRepository authRepository;

    setUp(() {
      mockApiService = MockApiService();
      mockSharedPrefs = MockAppSharedPrefs();
      authRepository = AuthRepository(mockApiService, mockSharedPrefs);
    });

    group('login', () {
      test('successfully logs in and saves auth data', () async {
        const email = 'test@example.com';
        const password = 'password123';
        final user = TestData.createUser(email: email);
        final authResponse = TestData.createAuthResponse(user: user);

        when(() => mockApiService.login(email, password))
            .thenAnswer((_) async => authResponse);
        when(() => mockSharedPrefs.saveAuthToken(any()))
            .thenAnswer((_) async {});
        when(() => mockSharedPrefs.saveUserId(any()))
            .thenAnswer((_) async {});
        when(() => mockSharedPrefs.saveUserEmail(any()))
            .thenAnswer((_) async {});

        final result = await authRepository.login(email, password);

        expect(result.accessToken, 'test-access-token');
        expect(result.user.email, email);
        verify(() => mockApiService.login(email, password)).called(1);
        verify(() => mockSharedPrefs.saveAuthToken('test-access-token')).called(1);
        verify(() => mockSharedPrefs.saveUserId(user.id)).called(1);
        verify(() => mockSharedPrefs.saveUserEmail(email)).called(1);
      });

      test('throws when login fails', () async {
        const email = 'test@example.com';
        const password = 'wrong-password';

        when(() => mockApiService.login(email, password))
            .thenThrow(Exception('Invalid credentials'));

        expect(
          () => authRepository.login(email, password),
          throwsA(isA<Exception>()),
        );
      });
    });

    group('register', () {
      test('successfully registers and saves auth data', () async {
        const email = 'new@example.com';
        const password = 'password123';
        final user = TestData.createUser(email: email);
        final authResponse = TestData.createAuthResponse(
          accessToken: 'new-token-456',
          user: user,
        );

        when(() => mockApiService.register(email, password))
            .thenAnswer((_) async => authResponse);
        when(() => mockSharedPrefs.saveAuthToken(any()))
            .thenAnswer((_) async {});
        when(() => mockSharedPrefs.saveUserId(any()))
            .thenAnswer((_) async {});
        when(() => mockSharedPrefs.saveUserEmail(any()))
            .thenAnswer((_) async {});

        final result = await authRepository.register(email, password);

        expect(result.accessToken, 'new-token-456');
        expect(result.user.email, email);
        verify(() => mockApiService.register(email, password)).called(1);
      });
    });

    group('logout', () {
      test('clears auth data', () async {
        when(() => mockSharedPrefs.clearAuthToken())
            .thenAnswer((_) async {});
        when(() => mockSharedPrefs.clearUser())
            .thenAnswer((_) async {});

        await authRepository.logout();

        verify(() => mockSharedPrefs.clearAuthToken()).called(1);
        verify(() => mockSharedPrefs.clearUser()).called(1);
      });
    });

    group('isAuthenticated', () {
      test('returns true when token exists', () async {
        when(() => mockSharedPrefs.authToken).thenReturn('valid-token');

        final result = await authRepository.isAuthenticated();

        expect(result, true);
      });

      test('returns false when token is null', () async {
        when(() => mockSharedPrefs.authToken).thenReturn(null);

        final result = await authRepository.isAuthenticated();

        expect(result, false);
      });

      test('returns false when token is empty', () async {
        when(() => mockSharedPrefs.authToken).thenReturn('');

        final result = await authRepository.isAuthenticated();

        expect(result, false);
      });
    });

    group('currentUser', () {
      test('returns user when userId and userEmail exist', () {
        const userId = 'user-123';
        const email = 'user@example.com';

        when(() => mockSharedPrefs.userId).thenReturn(userId);
        when(() => mockSharedPrefs.userEmail).thenReturn(email);

        final user = authRepository.currentUser;

        expect(user, isNotNull);
        expect(user!.id, userId);
        expect(user.email, email);
      });

      test('returns null when userId is null', () {
        when(() => mockSharedPrefs.userId).thenReturn(null);
        when(() => mockSharedPrefs.userEmail).thenReturn('user@example.com');

        final user = authRepository.currentUser;

        expect(user, isNull);
      });

      test('returns null when userEmail is null', () {
        when(() => mockSharedPrefs.userId).thenReturn('user-123');
        when(() => mockSharedPrefs.userEmail).thenReturn(null);

        final user = authRepository.currentUser;

        expect(user, isNull);
      });
    });
  });
}
