# 启动指南

## 前置要求

- Docker Desktop (Windows/Mac) 或 Docker Engine (Linux)
- Go 1.21+
- Flutter 3.16+ (Dart SDK >=3.2.0 <4.0.0)

## Docker Desktop与 WSL2集成
- https://learn.microsoft.com/zh-cn/windows/wsl/tutorials/wsl-containers

## Go安装
```bash
curl -O -L https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
# 删除已有安装并解压到usr/local
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
# 设置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
```


## Flutter版本管理

本项目需要Flutter 3.16或更高版本。你可以通过以下方式安装指定版本：

### 方式1: 使用FVM (Flutter Version Management) - 推荐

FVM可以方便地管理多个Flutter版本：

```bash
# 安装FVM (如果还没有安装)
# macOS/Linux:
brew tap leoafarias/fvm
brew install fvm

# Windows (使用Chocolatey):
choco install fvm

# 或者通过Pub安装:
dart pub global activate fvm

# 在项目中安装并使用指定Flutter版本
cd flutter
fvm install 3.16.9
fvm use 3.16.9

# 后续使用fvm flutter命令代替flutter
fvm flutter pub get
fvm flutter pub run build_runner build --delete-conflicting-outputs
fvm flutter run -d chrome
```

### 方式2: 直接安装指定版本

从Flutter官网下载安装：
- 访问 https://flutter.dev/docs/development/tools/sdk/releases
- 下载Flutter 3.16.x版本
- 解压并配置PATH环境变量

```bash
# 验证Flutter版本
flutter --version
flutter doctor
```

## 步骤 1: 启动数据库和Redis

```bash
# 在项目根目录下运行
docker-compose up -d
```

等待PostgreSQL和Redis完全启动（约10-30秒）。

## 步骤 2: 启动Go服务器

```bash
cd server
cp config.yaml.example config.yaml  # 如果还没有配置文件
go mod download
go run cmd/server/main.go
```

服务器将在 `http://localhost:8080` 启动。

## 步骤 3: 启动Flutter Web应用

打开一个新的终端窗口：

### 使用标准Flutter命令
```bash
cd flutter
flutter pub get
flutter pub run build_runner build --delete-conflicting-outputs
flutter run -d chrome
```

### 使用FVM (如果安装了FVM)
```bash
cd flutter
fvm flutter pub get
fvm flutter pub run build_runner build --delete-conflicting-outputs
fvm flutter run -d chrome
```

Flutter Web应用将在Chrome浏览器中自动打开。

## 开发端口

- Go Server: http://localhost:8080
- Flutter Web: http://localhost:xxxx (Flutter会自动分配端口)
- PostgreSQL: localhost:5432
- Redis: localhost:6379

## 停止服务

```bash
# 停止Go服务器和Flutter应用: Ctrl+C

# 停止数据库和Redis
docker-compose down

# 停止并删除数据 volumes (谨慎使用)
docker-compose down -v
```

## 故障排除

### Docker未运行
确保Docker Desktop已启动。

### 端口被占用
如果端口被占用，可以修改 `docker-compose.yml` 和 `server/config.yaml` 中的端口配置。

### Flutter没有找到
确保Flutter已添加到PATH中，运行 `flutter doctor` 检查安装状态。
