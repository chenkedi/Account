# Account Web Client

基于 React + TypeScript + Vite + Chakra UI 构建的个人财务管理应用。

## 技术栈

- **框架**: React 18 + TypeScript 5
- **构建工具**: Vite 5
- **UI 组件库**: Chakra UI
- **状态管理**: Zustand
- **本地数据库**: Dexie.js (IndexedDB)
- **路由**: React Router 6
- **表单**: React Hook Form + Zod
- **图表**: Recharts
- **跨平台**: Capacitor 5 (Web → Android)

## 项目结构

```
web/
├── src/
│   ├── core/               # 核心基础设施
│   │   ├── constants/      # 常量定义
│   │   ├── theme/          # 主题配置
│   │   └── utils/          # 工具函数
│   ├── data/               # 数据层
│   │   ├── models/         # 数据模型 (Zod)
│   │   ├── database/       # Dexie 数据库
│   │   ├── api/            # API 客户端
│   │   ├── repositories/   # 数据仓库
│   │   └── storage/        # 本地存储
│   ├── sync/               # 同步引擎
│   │   ├── sync-manager.ts
│   │   ├── lww-strategy.ts
│   │   └── websocket-client.ts
│   ├── store/              # Zustand 状态管理
│   ├── presentation/       # UI 层
│   │   ├── components/     # 组件
│   │   ├── pages/          # 页面
│   │   ├── providers/      # Context Providers
│   │   └── router/         # 路由配置
│   ├── App.tsx
│   └── main.tsx
├── capacitor.config.ts
├── package.json
└── vite.config.ts
```

## 开发指南

### 前置要求

- Node.js 18+
- npm 或 yarn
- Go 后端服务 (在 `../server/` 目录)

### 安装依赖

```bash
cd web
npm install
```

### 启动开发服务器

**1. 启动 Go 后端服务** (在新终端中):

```bash
cd server
go run cmd/server/main.go
```

**2. 启动 React 开发服务器**:

```bash
cd web
npm run dev
```

然后在浏览器中访问: `http://localhost:5173`

### 可用命令

```bash
npm run dev          # 启动开发服务器
npm run build        # 构建生产版本
npm run preview      # 预览生产构建
npm run lint         # 代码检查
npm run format       # 代码格式化
npm run test         # 运行测试
npm run typecheck    # TypeScript 类型检查
```

## Capacitor (Android 打包)

### 初始化 Capacitor

```bash
npm run build
npm run cap:init
npm run cap:add:android
```

### 同步到 Android

```bash
npm run build
npm run cap:sync
npm run cap:open:android
```

## WSL2 开发注意事项

### 网络配置

确保 Go 后端监听在 `0.0.0.0:8080` 而不是 `localhost:8080`，以便 WSL2 可以访问。

### 文件监听

Vite 配置已启用 `usePolling: true` 以支持 WSL2 文件系统监听。

### 浏览器测试

可以直接在 Windows 的浏览器中访问 `http://localhost:5173` 来测试应用。

## 功能特性

- [x] 用户认证 (登录/注册)
- [x] 账户管理
- [x] 交易记录
- [x] 分类管理
- [ ] 离线优先数据存储
- [ ] LWW 同步策略
- [ ] WebSocket 实时同步
- [ ] 账单导入 (支付宝/微信/银行)
- [ ] 统计报表
- [ ] 图表展示

## License

MIT
