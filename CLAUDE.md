# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the **account** project - a complete standalone server-client application for daily accounting and personal financial management.

### 核心功能

**历史记账初始化**：
- 支持用户上传银行流水表格、支付宝微信账单表格
- CSV 解析器支持：支付宝、微信支付、银行对账单、通用 CSV
- 智能账户匹配和重复检测
- 导入预览功能，确认后再导入

**每日增量记账**：
- **自动模式** (规划中)：要求用户授权支付宝、微信、京东的账单 API
- **手动模式** (UI 占位)：用户输入金额，备注信息，选择所用账户等
- 快速添加按钮和导入入口集成在交易页面

**统计功能**：
- 支持按日期范围筛选（本周/本月/本季度/本年/近6个月/自定义）
- 按类别聚合账目，收入和支出分类饼图
- 月收支差趋势折线图
- 交互式图表展示

## Architecture

- **Pattern**: Server-client architecture
- **Clients**: Flutter (single codebase for Android + Web)
- **Server**: Go + PostgreSQL + Redis

**Key Features**:
- Data synchronization between multiple clients
- 完善的客户端本地和服务端版本管理系统，支持多端同步修改的数据同步
- Last-Write-Wins (LWW) 时间戳冲突解决策略
- 完善且可靠的 client 侧本地数据库管理 (Drift/SQLite)
- 鲁棒且高性能的 server 侧数据库管理 (PostgreSQL)
- 高性能的 server，包括同步速度，较低的网络带宽要求
- WebSocket 实时同步通知

## Project Status

**COMPLETE - All 9 phases implemented ✓**

- Phase 1: Server Foundation - Completed ✓
- Phase 2: Core Data Model & APIs - Completed ✓
- Phase 3: Sync Engine - Completed ✓
- Phase 4: Flutter App Foundation - Completed ✓
- Phase 5: Flutter Core Features - Completed ✓
- Phase 6: Flutter Sync Integration - Completed ✓
- Phase 7: Import Service - Completed ✓
- Phase 8: Statistics & Reporting - Completed ✓
- Phase 9: Polish & Optimization - Completed ✓

### Compilation Fixes (2026-02-19)

Fixed the following compilation errors for Flutter 3.16.9:

1. **Repository Imports**: Fixed missing app_database.dart imports in repositories
2. **ColorScheme**: Removed `surfaceContainerHighest` (not available in 3.16.9), added required `background` and `onBackground` parameters
3. **SyncManager**: Fixed constructor parameter mismatch (7 → 6 parameters)
4. **TransactionDao**: Fixed `FutureOr<double>` type issue in fold operations by using explicit loops
5. **ApiService**: Added missing dio import for FormData and MultipartFile
6. **Import Widgets**: Resolved naming conflicts by:
   - Renaming widget classes: `ImportPreview` → `ImportPreviewWidget`, `ImportResult` → `ImportResultWidget`
   - Using prefix import `as models` for data models
7. **Const Constructors**: Removed `const` keyword for non-const constructor calls
8. **fl_chart**: Changed `getTooltipColor` → `tooltipBgColor` for compatibility with older version
9. **Assets**: Removed missing asset directory references from pubspec.yaml

Linux desktop build now succeeds with exit code 0.

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

## Project Structure

```
account/
├── server/                    # Go backend server
│   ├── cmd/server/main.go     # Entry point
│   ├── internal/
│   │   ├── api/               # HTTP handlers, routes
│   │   ├── business/          # Services and models
│   │   ├── data/              # Repositories and database
│   │   └── sync/              # Sync engine
│   └── pkg/                   # Utilities (auth, config, logger)
├── flutter/                   # Flutter mobile and web app
│   ├── lib/
│   │   ├── core/              # Constants, theme, network, utils
│   │   ├── data/              # Models, repositories, API service
│   │   ├── domain/            # (Clean architecture layer placeholder)
│   │   ├── presentation/      # UI pages, BLoCs, widgets
│   │   ├── sync/              # Sync manager, LWW strategy
│   │   ├── main.dart          # App entry
│   │   ├── app.dart           # App initialization
│   │   └── injection_container.dart  # Dependency injection
├── CLAUDE.md                  # This file (project specifications)
└── README.md                  # User-facing README
```

## Development Commands

### Server

```bash
cd server
go mod download
go run cmd/server/main.go
```

### Flutter App

**IMPORTANT**: All Flutter commands must be executed from the `flutter/` directory using `fvm` (Flutter Version Manager):

```bash
cd flutter
fvm flutter pub get
fvm flutter pub run build_runner build --delete-conflicting-outputs
fvm flutter run          # Android
fvm flutter run -d chrome  # Web
fvm flutter run -d linux   # Linux desktop
```

To add desktop platform support (e.g., Linux):
```bash
cd flutter
fvm flutter create --platforms=linux --project-name=account .
```

## Key Architectural Patterns

- **Clean Architecture**: Data → Domain → Presentation layers
- **BLoC Pattern**: Events in, States out for predictable state management
- **Repository Pattern**: Data access abstraction
- **Offline-First**: Write to local DB first, sync in background
- **LWW Sync**: Simple timestamp-based conflict resolution
- **Dependency Injection**: get_it for service location

## Sync Protocol

1. **Handshake**: Client sends last sync timestamp and device ID
2. **Pull**: Server sends all entities modified since last sync
3. **Push**: Client sends all locally modified entities since last sync
4. **LWW Resolution**: For any entity modified both sides, later `last_modified_at` wins
5. **Commit**: Both sides apply all changes, update sync timestamp
6. **Notification**: WebSocket notifies other devices to sync

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


## Notification
- execute `proxy_http` bash command before executing any ohter bash command which needs downloading or online request