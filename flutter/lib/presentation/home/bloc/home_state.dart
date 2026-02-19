part of 'home_bloc.dart';

enum HomeTab { dashboard, transactions, accounts, stats, settings }

class HomeState extends Equatable {
  final HomeTab currentTab;

  const HomeState({
    this.currentTab = HomeTab.dashboard,
  });

  HomeState copyWith({
    HomeTab? currentTab,
  }) {
    return HomeState(
      currentTab: currentTab ?? this.currentTab,
    );
  }

  @override
  List<Object?> get props => [currentTab];
}
