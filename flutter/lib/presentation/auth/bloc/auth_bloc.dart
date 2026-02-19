import 'dart:async';

import 'package:bloc/bloc.dart';
import 'package:equatable/equatable.dart';

import '../../../../data/models/user.dart';
import '../../../../data/repositories/auth_repository.dart';
import '../../../../sync/sync_manager.dart';
import '../../../../injection_container.dart' as di;

part 'auth_event.dart';
part 'auth_state.dart';

class AuthBloc extends Bloc<AuthEvent, AuthState> {
  final AuthRepository _authRepository;
  final SyncManager _syncManager;

  AuthBloc({
    required AuthRepository authRepository,
  })  : _authRepository = authRepository,
        _syncManager = di.sl<SyncManager>(),
        super(const AuthState()) {
    on<AuthCheckRequested>(_onAuthCheckRequested);
    on<AuthLoginRequested>(_onAuthLoginRequested);
    on<AuthRegisterRequested>(_onAuthRegisterRequested);
    on<AuthLogoutRequested>(_onAuthLogoutRequested);
  }

  Future<void> _onAuthCheckRequested(
    AuthCheckRequested event,
    Emitter<AuthState> emit,
  ) async {
    emit(state.copyWith(status: AuthStatus.loading));

    final isAuthenticated = await _authRepository.isAuthenticated();

    if (isAuthenticated) {
      final user = _authRepository.currentUser;
      if (user != null) {
        // Initialize sync manager on auth check success
        unawaited(_syncManager.initialize());
        emit(state.copyWith(
          status: AuthStatus.authenticated,
          user: user,
        ));
      } else {
        emit(state.copyWith(status: AuthStatus.unauthenticated));
      }
    } else {
      emit(state.copyWith(status: AuthStatus.unauthenticated));
    }
  }

  Future<void> _onAuthLoginRequested(
    AuthLoginRequested event,
    Emitter<AuthState> emit,
  ) async {
    emit(state.copyWith(
      status: AuthStatus.loading,
      errorMessage: null,
    ));

    try {
      final response = await _authRepository.login(event.email, event.password);
      // Initialize sync manager on login
      unawaited(_syncManager.initialize());
      emit(state.copyWith(
        status: AuthStatus.authenticated,
        user: response.user,
      ));
    } catch (e) {
      emit(state.copyWith(
        status: AuthStatus.unauthenticated,
        errorMessage: e.toString(),
      ));
    }
  }

  Future<void> _onAuthRegisterRequested(
    AuthRegisterRequested event,
    Emitter<AuthState> emit,
  ) async {
    emit(state.copyWith(
      status: AuthStatus.loading,
      errorMessage: null,
    ));

    try {
      final response = await _authRepository.register(event.email, event.password);
      // Initialize sync manager on register
      unawaited(_syncManager.initialize());
      emit(state.copyWith(
        status: AuthStatus.authenticated,
        user: response.user,
      ));
    } catch (e) {
      emit(state.copyWith(
        status: AuthStatus.unauthenticated,
        errorMessage: e.toString(),
      ));
    }
  }

  Future<void> _onAuthLogoutRequested(
    AuthLogoutRequested event,
    Emitter<AuthState> emit,
  ) async {
    emit(state.copyWith(status: AuthStatus.loading));
    await _authRepository.logout();
    await _syncManager.dispose();
    emit(const AuthState(status: AuthStatus.unauthenticated));
  }
}
