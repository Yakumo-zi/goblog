# GoBlog 生产环境部署指南

本文档介绍如何将GoBlog项目部署到生产服务器，包括Docker镜像的导出、传输和部署流程。

## 📋 部署概览

GoBlog项目支持两种主要的部署方式：

1. **离线部署** - 导出Docker镜像并传输到服务器
2. **在线部署** - 直接在服务器上构建和运行

## 🚀 方式一：离线部署（推荐）

### 1. 导出Docker镜像

使用内置的镜像导出工具将所有必需的Docker镜像导出并上传到服务器：

```bash
# 基本用法
./scripts/export_and_upload_images.sh <server_ip> <username> <remote_path>

# 示例
./scripts/export_and_upload_images.sh 192.168.1.100 root /opt/goblog/

# 或者使用Makefile
make docker-export HOST=192.168.1.100 USER=root PATH=/opt/goblog/
```

### 2. 支持的选项

```bash
# 指定SSH端口
./scripts/export_and_upload_images.sh -p 2222 192.168.1.100 deploy /home/deploy/goblog/

# 使用私钥认证
./scripts/export_and_upload_images.sh -i ~/.ssh/id_rsa 192.168.1.100 root /opt/goblog/

# 保留本地导出文件
./scripts/export_and_upload_images.sh -k 192.168.1.100 root /opt/goblog/

# 预览模式（不实际执行）
./scripts/export_and_upload_images.sh --dry-run 192.168.1.100 root /opt/goblog/
```

### 3. 在服务器上导入镜像

脚本会自动在服务器上创建导入脚本，登录服务器后运行：

```bash
ssh root@192.168.1.100
cd /opt/goblog/
./import_images.sh
```

## 🔧 方式二：在线部署

### 1. 准备服务器环境

确保服务器已安装：
- Docker
- Docker Compose
- Git (如果从源码部署)

### 2. 传输项目文件

```bash
# 方法A: 使用Git
git clone <your-repo-url> goblog
cd goblog

# 方法B: 使用SCP传输项目文件
scp -r . user@server:/opt/goblog/
```

### 3. 构建和运行

```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose -f docker-compose.prod.yml up -d
```

## 🌐 生产环境配置

### 环境变量配置

创建 `.env` 文件配置生产环境参数：

```bash
# 数据库配置
POSTGRES_DB=goblog_prod
POSTGRES_USER=goblog_user
POSTGRES_PASSWORD=your_secure_password_here

# 应用配置
APP_PORT=8080
JWT_SECRET=your_very_secure_jwt_secret_key_here
LOG_LEVEL=warn
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your_secure_admin_password

# 备份配置
BACKUP_SCHEDULE=0 2 * * *  # 每天凌晨2点
BACKUP_RETENTION_DAYS=30   # 保留30天备份

# Nginx配置 (如果使用)
NGINX_HTTP_PORT=80
NGINX_HTTPS_PORT=443
```

### 使用生产环境配置

```bash
# 使用生产环境Docker Compose文件
docker-compose -f docker-compose.prod.yml up -d

# 或者包含Nginx反向代理
docker-compose -f docker-compose.prod.yml --profile with-nginx up -d
```

## 🛡️ 安全配置

### 1. 防火墙设置

```bash
# 开放必要端口
ufw allow 22      # SSH
ufw allow 80      # HTTP
ufw allow 443     # HTTPS
ufw allow 8080    # API (如果直接暴露)

# 限制数据库端口访问 (仅本地)
ufw deny 5432
```

### 2. SSL/TLS配置

如果使用Nginx，创建SSL证书配置：

```bash
# 创建SSL证书目录
mkdir -p ssl/

# 使用Let's Encrypt (示例)
certbot certonly --standalone -d yourdomain.com
cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem ssl/
cp /etc/letsencrypt/live/yourdomain.com/privkey.pem ssl/
```

### 3. 定期更新

```bash
# 创建定期更新脚本
cat > update.sh << 'EOF'
#!/bin/bash
cd /opt/goblog
docker-compose -f docker-compose.prod.yml pull
docker-compose -f docker-compose.prod.yml up -d
docker image prune -f
EOF

chmod +x update.sh

# 添加到crontab (每周检查更新)
echo "0 3 * * 0 /opt/goblog/update.sh" | crontab -
```

## 📊 监控和维护

### 1. 健康检查

```bash
# 检查服务状态
docker-compose -f docker-compose.prod.yml ps

# 查看日志
docker-compose -f docker-compose.prod.yml logs -f

# 检查API健康状态
curl http://localhost:8080/health
```

### 2. 备份管理

```bash
# 手动创建备份
docker exec goblog-postgres-prod pg_dump -U goblog_user goblog_prod > backup.sql

# 查看自动备份
docker exec goblog-postgres-backup-prod ls -la /backups/

# 恢复备份
docker exec -i goblog-postgres-prod psql -U goblog_user goblog_prod < backup.sql
```

### 3. 性能优化

```bash
# 清理未使用的Docker资源
docker system prune -f

# 查看资源使用情况
docker stats

# 数据库性能调优
docker exec goblog-postgres-prod psql -U goblog_user -d goblog_prod -c "
SELECT schemaname,tablename,attname,n_distinct,most_common_vals 
FROM pg_stats 
WHERE schemaname = 'public';"
```

## 🔄 镜像导出工具详细说明

### 导出的镜像列表

工具会自动导出以下镜像：
- `postgres:15-alpine` - PostgreSQL数据库
- `dpage/pgadmin4:latest` - pgAdmin管理界面
- `alpine:3.18` - 备份服务基础镜像
- `goblog-goblog:latest` - 应用程序镜像

### 工具特性

- ✅ **自动检测** - 自动检测缺失的镜像
- ✅ **压缩传输** - 支持压缩以减少传输时间
- ✅ **SSH认证** - 支持密码和私钥认证
- ✅ **干运行模式** - 支持预览操作而不实际执行
- ✅ **自动导入** - 自动创建远程导入脚本
- ✅ **错误处理** - 完善的错误检查和处理

### 故障排除

#### 连接问题
```bash
# 测试SSH连接
ssh -p 22 user@server "echo 'Connection test'"

# 检查防火墙
ufw status

# 验证SSH密钥
ssh-add -l
```

#### 镜像问题
```bash
# 检查本地镜像
docker images

# 重新构建应用镜像
docker-compose build --no-cache

# 手动拉取镜像
docker pull postgres:15-alpine
```

#### 空间问题
```bash
# 检查磁盘空间
df -h

# 清理Docker空间
docker system prune -a
```

## 📞 支持

如果在部署过程中遇到问题：

1. 检查日志：`docker-compose logs`
2. 验证配置：`docker-compose config`
3. 查看资源：`docker stats`
4. 网络测试：`docker network ls`

## 🔗 相关文档

- [README.md](README.md) - 项目概览和快速开始
- [API文档](README.md#api文档) - API接口说明
- [开发指南](README.md#开发工具) - 本地开发环境 