import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import 'app.dart';
import 'injection_container.dart' as di;
import 'presentation/app_bloc_observer.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // Initialize dependencies
  await di.initDependencies();

  // Set up bloc observer for debugging
  Bloc.observer = const AppBlocObserver();

  runApp(const App());
}
