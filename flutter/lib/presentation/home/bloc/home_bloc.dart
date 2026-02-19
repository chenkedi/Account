import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

part 'home_state.dart';

final class HomeTabChanged extends Equatable {
  final HomeTab tab;

  const HomeTabChanged(this.tab);

  @override
  List<Object?> get props => [tab];
}

class HomeBloc extends Bloc<HomeTabChanged, HomeState> {
  HomeBloc() : super(const HomeState()) {
    on<HomeTabChanged>(_onTabChanged);
  }

  void _onTabChanged(
    HomeTabChanged event,
    Emitter<HomeState> emit,
  ) {
    emit(state.copyWith(currentTab: event.tab));
  }
}
