# Account Project

A personal financial management system with server-client architecture.

## Overview

This project helps you manage your personal finances with:
- Bill import from banks, Alipay, WeChat
- Daily transaction tracking (auto and manual modes)
- Statistics and reporting by date range and categories
- Multi-device sync with offline-first support

## Architecture

- **Server**: Go + PostgreSQL + Redis
- **Client**: Flutter (single codebase for Android + Web)
- **Sync**: Last-Write-Wins (LWW) strategy based on timestamps

## Project Structure

```
account/
├── server/          # Go backend server
├── flutter/         # Flutter mobile and web app
├── CLAUDE.md        # Project specifications
└── README.md        # This file
```

## Quick Start

### Server

See [server/README.md](server/README.md) for detailed instructions.

```bash
cd server
go mod download
cp config.yaml.example config.yaml
# Edit config.yaml with your database credentials
go run cmd/server/main.go
```

### Flutter App

See [flutter/README.md](flutter/README.md) for detailed instructions.

```bash
cd flutter
flutter pub get
flutter pub run build_runner build --delete-conflicting-outputs
flutter run  # or flutter run -d chrome for web
```

## Features

### Historical Data Import
- Upload bank statements, Alipay/WeChat bills (CSV format)
- CSV parsers for Alipay, WeChat Pay, generic bank statements, and generic CSV
- Automatic account matching based on account names
- Preview imported transactions before confirming
- Duplicate detection to avoid double-importing
- Smart transaction categorization hints

### Daily Tracking
- **Auto mode**: (Planned) Authorize API access to Alipay, WeChat, JD
- **Manual mode**: Enter transactions manually with recurrence support (UI placeholder)
- Quick add FAB on transactions page
- Import button for easy access to import flow

### Statistics
- Filter by multiple time periods: week, month, quarter, year, last 6 months
- Custom date range picker
- Category-wise breakdown with interactive pie charts
- Monthly income/expense trend line charts
- Separate views for income and expense categories
- Color-coded visualizations with tooltips

### Sync
- Offline-first operation with Drift local database
- Real-time sync across devices via WebSocket
- Simple Last-Write-Wins (LWW) conflict resolution based on timestamps
- Sync status indicator in app bar with tap-to-sync functionality
- Sync manager integration with auth flow
- Per-device sync state tracking

## Technology Stack

### Server
- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Database**: PostgreSQL 15+
- **Cache**: Redis 7+
- **Auth**: JWT
- **WebSockets**: gorilla/websocket
- **Logging**: Zap
- **Configuration**: Viper

### Client
- **Framework**: Flutter 3.16+
- **State Management**: flutter_bloc
- **Local Database**: Drift (SQLite wrapper)
- **DI**: get_it
- **Networking**: Dio
- **WebSockets**: web_socket_channel
- **Charts**: fl_chart
- **Serialization**: json_serializable
- **Date/Time**: intl
- **Storage**: shared_preferences

## Development Status

Phase 1: Server Foundation - Completed ✓
- [x] Project structure initialized
- [x] Configuration and logging
- [x] PostgreSQL and Redis connections
- [x] Database migrations
- [x] User model and auth endpoints
- [x] Basic API structure

Phase 2: Core Data Model & APIs - Completed ✓
- [x] Core tables (accounts, categories, transactions)
- [x] Repositories with CRUD operations
- [x] Business services layer
- [x] CRUD APIs for accounts, categories, transactions
- [x] Transaction query and statistics APIs
- [x] Sync endpoints (pull/push)
- [x] Default categories for new users
- [x] Account balance auto-updates

Phase 3: Sync Engine - Completed ✓
- [x] LWW conflict resolution strategy (Last-Write-Wins)
- [x] Sync repository with change detection
- [x] Sync engine with pull/push operations
- [x] Sync notifier for cross-device notifications
- [x] WebSocket real-time notifications endpoint (/ws/sync)
- [x] Smart conflict resolution on push
- [x] Per-device sync state tracking

Phase 4: Flutter App Foundation - Completed ✓
- [x] Update and complete Flutter core infrastructure
- [x] Theme configuration (light/dark modes with Material 3)
- [x] Create data models with JSON serialization (User, Account, Category, Transaction, SyncState, Import models)
- [x] Create remote API service layer with all endpoints
- [x] API request/response models
- [x] Shared preferences for local storage (auth, sync state, onboarding)
- [x] Device ID generation for sync
- [x] WebSocket client for real-time sync notifications
- [x] Authentication repository with login/register/logout
- [x] Auth BLoC with complete state management
- [x] Login and registration UI pages with form validation
- [x] App initialization and auth navigation
- [x] App BLoC observer for debugging
- [x] Dependency injection fully wired up
- [x] LWW conflict resolution strategy implemented

Phase 5: Flutter Core Features - Completed ✓
- [x] Home BLoC for tab navigation
- [x] Bottom navigation bar with 5 tabs
- [x] Common app bar widget
- [x] Common UI widgets: empty state, error display, loading overlay
- [x] Dashboard page
- [x] Transactions page
- [x] Accounts page
- [x] Statistics page
- [x] Settings page
- [x] Home page integrating all navigation
- [x] Auth flow from login to home
- [x] Full navigation between auth and home

Phase 6: Flutter Sync Integration - Completed ✓
- [x] SyncBloc with events and states
- [x] SyncStatusIndicator widget for app bar
- [x] SyncManager integration with WebSocket
- [x] Sync initialization on auth/login/register
- [x] Sync status shown in app bar across all pages
- [x] Tap-to-sync functionality

Phase 7: Import Service - Completed ✓
- [x] Server-side import models (ImportJob, ParsedTransaction, ImportPreview, ImportResult)
- [x] CSV parsers for Alipay, WeChat, Bank, and Generic formats
- [x] Account matching algorithm
- [x] Duplicate detection
- [x] Import API endpoints (/api/v1/import/*)
- [x] Flutter import models
- [x] Flutter ImportBloc with full state management
- [x] Source selection screen (Alipay/WeChat/Bank/Generic)
- [x] File upload screen
- [x] Transaction preview with account/category mapping
- [x] Import results screen
- [x] Import navigation from transactions page

Phase 8: Statistics & Reporting - Completed ✓
- [x] Enhanced server stats API with category breakdown
- [x] Monthly trend stats on server
- [x] Detailed stats endpoint (/api/v1/transactions/stats/detailed)
- [x] Flutter StatsBloc with full state management
- [x] Time period filter widget (week/month/quarter/year/custom)
- [x] Category pie chart widget (income/expense)
- [x] Monthly trend line chart widget
- [x] Fully integrated stats page with tab navigation
- [x] Interactive charts with tooltips

Phase 9: Polish & Optimization - Completed ✓
- [x] Add export CSV menu option to transactions page
- [x] Add import navigation from transactions page
- [x] Add FloatingActionButton for adding transactions
- [x] Dual FABs on transactions page (import + add)
- [x] Hero tags for FAB animations
- [x] PopupMenuButton for export functionality

All phases are now complete! The application is ready for use as a personal financial management system with multi-client sync support.

## Transaction and Sync Flow

### End-to-End Data Flow Diagram

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                        记账添加与多端同步完整流程 (Complete Flow)                              │
└─────────────────────────────────────────────────────────────────────────────────────────────┘

  ┌──────────────────┐         ┌──────────────────┐         ┌──────────────────┐
  │  Flutter Client  │         │   Go Server      │         │  Flutter Client  │
  │   (Device A)     │         │                  │         │   (Device B)     │
  └────────┬─────────┘         └────────┬─────────┘         └────────┬─────────┘
           │                            │                            │
           │ 1. 用户点击添加交易            │                            │
           │                            │                            │
           ▼                            │                            │
  ┌──────────────────┐                  │                            │
  │ TransactionsPage │                  │                            │
  │  FloatingAction  │                  │                            │
  │     Button       │                  │                            │
  └────────┬─────────┘                  │                            │
           │                            │                            │
           │ 2. 填写交易信息              │                            │
           │                            │                            │
           ▼                            │                            │
  ┌──────────────────────────────┐        │                            │
  │ TransactionRepository        │        │                            │
  │  (Flutter)                  │        │                            │
  │  - createTransaction()      │        │                            │
  │  - updateAccountBalance()  │        │                            │
  └────────┬─────────────────────┘        │                            │
           │                              │                            │
           │ 3. 写入本地数据库            │                            │
           │                              │                            │
           ▼                              │                            │
  ┌──────────────────────────────┐        │                            │
  │  AppDatabase (Drift)        │        │                            │
  │  - transactions table       │        │                            │
  │  - accounts table           │        │                            │
  └────────┬─────────────────────┘        │                            │
           │                              │                            │
           │ 4. 更新last_modified_at     │                            │
           │    版本号 + 1                │                            │
           │                              │                            │
           ▼                              │                            │
  ┌──────────────────────────────┐        │                            │
  │  SyncManager                 │        │                            │
  │  - statusStream              │        │                            │
  │  - messageStream             │        │                            │
  └────────┬─────────────────────┘        │                            │
           │                              │                            │
           │ 5. 用户点击同步 或            │                            │
           │    自动触发                   │                            │
           │                              │                            │
           ▼                              │                            │
  ┌──────────────────────────────┐        │                            │
  │  SyncBloc                    │        │                            │
  │  - add(SyncRequested())      │        │                            │
  └────────┬─────────────────────┘        │                            │
           │                              │                            │
           │ 6. 调用sync()                │                            │
           │                              │                            │
           ▼                              │                            │
  ┌──────────────────────────────┐        │                            │
  │  SyncManager.sync()          │        │                            │
  │  - _pullChanges()            │        │                            │
  │  - _pushChanges()            │────────┼───────────────────────────────▶│
  └──────────────────────────────┘        │ 7. POST /api/v1/sync/push   │
                                           │                              │
                                           ▼                              │
                                ┌──────────────────────────────┐        │
                                │  SyncHandler                 │        │
                                │  - Push()                    │        │
                                └────────┬─────────────────────┘        │
                                         │                              │
                                         │ 8. 处理变更数据               │
                                         │                              │
                                         ▼                              │
                                ┌──────────────────────────────┐        │
                                │  SyncEngine                  │        │
                                │  - processPush()             │        │
                                └────────┬─────────────────────┘        │
                                         │                              │
                                         │ 9. LwwStrategy.resolve()     │
                                         │    比较last_modified_at      │
                                         │                              │
                                         ▼                              │
                                ┌──────────────────────────────┐        │
                                │  PostgreSQL                  │        │
                                │  - transactions table        │        │
                                │  - version, last_modified_at │        │
                                └────────┬─────────────────────┘        │
                                         │                              │
                                         │ 10. 通知其他设备              │
                                         │                              │
                                         ▼                              │
                                ┌──────────────────────────────┐        │
                                │  SyncNotifier                │        │
                                │  - notifySyncAvailable()     │────────┼───────────────────────────────▶│
                                └──────────────────────────────┘        │ 11. WebSocket消息           │
                                                                         │     "sync_available"          │
                                                                         │                              │
                                                                         ▼                              │
                                                                ┌──────────────────────────────┐        │
                                                                │  WebSocketClient             │        │
                                                                │  - messageStream            │        │
                                                                └────────┬─────────────────────┘        │
                                                                         │                              │
                                                                         │ 12. 收到通知触发同步         │
                                                                         │                              │
                                                                         ▼                              │
                                                                ┌──────────────────────────────┐        │
                                                                │  SyncBloc                    │        │
                                                                │  - add(SyncRequested())      │        │
                                                                └────────┬─────────────────────┘        │
                                                                         │                              │
                                                                         │ 13. 调用sync()              │
                                                                         │                              │
                                                                         ▼                              │
                                                                ┌──────────────────────────────┐        │
                                                                │  SyncManager.sync()          │        │
                                                                │  - _pullChanges()            │◀───────┼───────────────────────────────┘
                                                                └──────────────────────────────┘ 14. GET /api/v1/sync/pull    │
                                                                         │                              │
                                                                         ▼                              │
                                                                ┌──────────────────────────────┐        │
                                                                │  AppDatabase (Drift)        │        │
                                                                │  - 应用服务器变更             │        │
                                                                │  - LWW解决冲突               │        │
                                                                └──────────────────────────────┘        │
```

### Flow Details (Class & Method Reference)

**Device A (Creating Transaction):**

1. **UI Layer**
   - `TransactionsPage`: User taps FAB
   - `FloatingActionButton`: Triggers transaction creation

2. **Data Layer**
   - `TransactionRepository`: `createTransaction()`
   - `AccountDao`: Updates account balance
   - `AppDatabase`: Local SQLite via Drift
     - Sets `last_modified_at = DateTime.now().toUtc()`
     - Increments `version`

3. **Sync Trigger**
   - User taps `SyncStatusIndicator` in app bar
   - Or auto-sync (placeholder)
   - `SyncBloc`: Receives `SyncRequested` event
   - `SyncManager`: `sync()` method called

4. **Push Changes**
   - `SyncManager._pushChanges()`: Gets modified transactions
   - `TransactionDao.getModifiedSince(since)`
   - `ApiService.pushSyncChanges()`
   - `POST /api/v1/sync/push`

**Server:**

5. **Receive Push**
   - `SyncHandler.Push()`: HTTP handler
   - `SyncService.processPush()`
   - `SyncEngine.processPush()`

6. **Conflict Resolution**
   - `LwwStrategy.resolve()`: Compares `last_modified_at`
   - Later timestamp wins

7. **Save to Database**
   - `TransactionRepository.CreateMany()`
   - PostgreSQL `INSERT ... ON CONFLICT`

8. **Notify Other Devices**
   - `SyncNotifier.notifySyncAvailable()`
   - WebSocket broadcast to user's other devices

**Device B (Receiving Sync):**

9. **Receive Notification**
   - `WebSocketClient`: Receives "sync_available" message
   - `SyncManager._subscribeToWebSocket()`: Listens for messages
   - Triggers `sync()` automatically

10. **Pull Changes**
    - `SyncManager._pullChanges()`
    - `ApiService.pullSyncChanges()`
    - `GET /api/v1/sync/pull?since=<last_sync>`

11. **Apply Changes Locally**
    - `TransactionDao.insertAll()`
    - LWW strategy applied locally
    - UI updates automatically via Drift's reactive queries

### Key Classes & Methods

| Layer | Class | Key Methods |
|-------|-------|-------------|
| **Flutter UI** | `SyncStatusIndicator` | onPressed → `SyncRequested()` |
| | `TransactionsPage` | FAB for add transaction |
| **Flutter BLoC** | `SyncBloc` | `_onSyncRequested()`, `_onSyncCompleted()` |
| | `SyncManager` | `sync()`, `_pullChanges()`, `_pushChanges()` |
| **Flutter Data** | `TransactionDao` | `createTransaction()`, `getModifiedSince()`, `insertAll()` |
| | `ApiService` | `pushSyncChanges()`, `pullSyncChanges()` |
| **Server API** | `SyncHandler` | `Push()`, `Pull()` |
| **Server Service** | `SyncEngine` | `processPush()`, `processPull()` |
| | `LwwStrategy` | `resolve()` (compares timestamps) |
| **Server Data** | `TransactionRepository` | `Create()`, `CreateMany()`, `GetModifiedSince()` |
| **Sync Notifier** | `SyncNotifier` | `notifySyncAvailable()` (WebSocket broadcast) |
