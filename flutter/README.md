# Account Flutter App

A cross-platform personal financial management app built with Flutter. Supports both Android and Web from a single codebase.

## Features

- Offline-first with local SQLite database (Drift)
- Real-time data sync with Go server
- Bank, Alipay, WeChat bill import
- Statistics and reporting with charts
- Multi-device sync with Last-Write-Wins conflict resolution

## Getting Started

### Prerequisites

- Flutter 3.16+
- Dart 3.2+
- For Android: Android Studio / Android SDK
- For Web: Chrome / any modern browser

### Setup

1. Install dependencies:
   ```bash
   flutter pub get
   ```

2. Generate code (Drift, JSON serialization):
   ```bash
   flutter pub run build_runner build --delete-conflicting-outputs
   ```

3. Run the app:

   For Android:
   ```bash
   flutter run
   ```

   For Web:
   ```bash
   flutter run -d chrome
   ```

## Project Structure

```
lib/
├── main.dart                  # App entry point
├── app.dart                   # App initialization
├── injection_container.dart   # Dependency injection
├── core/                      # Core utilities and constants
│   ├── constants/
│   ├── errors/
│   ├── network/
│   ├── theme/
│   ├── usecase/
│   └── utils/
├── data/                      # Data layer
│   ├── datasources/
│   │   ├── local/            # Local DB (Drift)
│   │   └── remote/           # API calls
│   ├── models/               # Data models
│   └── repositories/         # Repositories
├── domain/                    # Domain layer
│   ├── entities/
│   ├── repositories/
│   └── usecases/
├── presentation/              # UI layer
│   ├── auth/
│   ├── home/
│   ├── dashboard/
│   ├── transactions/
│   ├── accounts/
│   ├── stats/
│   ├── import/
│   ├── settings/
│   ├── sync/
│   └── widgets/
└── sync/                      # Sync engine
    ├── sync_manager.dart
    ├── lww_strategy.dart
    └── sync_service.dart
```

## Code Generation

When modifying Drift tables, entities, or JSON serializable classes:

```bash
# One-time build
flutter pub run build_runner build --delete-conflicting-outputs

# Watch mode for development
flutter pub run build_runner watch
```

## Development

### Architecture Principles

- **Clean Architecture**: Separation of concerns across data/domain/presentation layers
- **Bloc Pattern**: Predictable state management
- **Offline-First**: Always write to local DB first, sync in background
- **Repository Pattern**: Single source of truth, mediates between local and remote

### Sync Strategy

Last-Write-Wins (LWW) based on `last_modified_at` timestamp:
- Entities compared by their last modification time
- Later timestamp always wins
- Deletes are soft and win if timestamp is later

## Tests

```bash
# Unit tests
flutter test

# Integration tests
flutter test integration_test/app_test.dart
```
