#!/bin/bash

# Docker镜像导出和上传脚本
# 用于将项目所有Docker镜像导出并通过SCP上传到指定服务器

set -e

# 设置环境变量，确保能找到常用命令
export PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:$PATH"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置变量
EXPORT_DIR="./docker-images"

# 获取日期后缀，使用多种方式尝试
get_date_suffix() {
    local date_cmd=""
    
    # 尝试不同的date命令路径
    for cmd in "date" "/bin/date" "/usr/bin/date" "/usr/sbin/date"; do
        if command -v "$cmd" &> /dev/null; then
            date_cmd="$cmd"
            break
        fi
    done
    
    if [ -z "$date_cmd" ]; then
        echo "$(date +%s)" # 使用时间戳作为后备
    else
        "$date_cmd" +%Y%m%d_%H%M%S
    fi
}

DATE_SUFFIX=$(get_date_suffix)
ARCHIVE_NAME="goblog-docker-images-${DATE_SUFFIX}.tar.gz"

# 项目使用的Docker镜像列表
declare -a IMAGES=(
    "postgres:15-alpine"
    "dpage/pgadmin4:latest"
    "alpine:3.18"
    "goblog-goblog:latest"  # 本地构建的镜像
)

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# 显示使用说明
show_usage() {
    echo "用法: $0 [options] <remote_host> <remote_user> <remote_path>"
    echo ""
    echo "参数:"
    echo "  remote_host    远程服务器地址 (例如: 192.168.1.100)"
    echo "  remote_user    远程服务器用户名 (例如: root)"
    echo "  remote_path    远程服务器存储路径 (例如: /opt/docker-images/)"
    echo ""
    echo "选项:"
    echo "  -h, --help     显示此帮助信息"
    echo "  -k, --keep     保留本地导出的镜像文件"
    echo "  -c, --compress 压缩镜像包 (默认开启)"
    echo "  -p, --port     SSH端口 (默认: 22)"
    echo "  -i, --identity SSH私钥文件路径"
    echo "  -P, --password 使用SSH密码认证 (交互式输入)"
    echo "  --password-file 从文件读取SSH密码"
    echo "  --dry-run      只显示将要执行的操作，不实际执行"
    echo ""
    echo "示例:"
    echo "  $0 192.168.1.100 root /opt/docker-images/"
    echo "  $0 -p 2222 -i ~/.ssh/id_rsa 192.168.1.100 deploy /home/deploy/images/"
    echo "  $0 -P 192.168.1.100 root /opt/docker-images/"
    echo "  $0 --password-file ~/.ssh/password 192.168.1.100 root /opt/docker-images/"
    echo "  $0 --dry-run 192.168.1.100 root /opt/docker-images/"
}

# 解析命令行参数
KEEP_LOCAL=false
COMPRESS=true
SSH_PORT=22
SSH_KEY=""
USE_PASSWORD=false
PASSWORD_FILE=""
SSH_PASSWORD=""
DRY_RUN=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_usage
            exit 0
            ;;
        -k|--keep)
            KEEP_LOCAL=true
            shift
            ;;
        -c|--compress)
            COMPRESS=true
            shift
            ;;
        -p|--port)
            SSH_PORT="$2"
            shift 2
            ;;
        -i|--identity)
            SSH_KEY="$2"
            shift 2
            ;;
        -P|--password)
            USE_PASSWORD=true
            shift
            ;;
        --password-file)
            PASSWORD_FILE="$2"
            shift 2
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        -*)
            log_error "未知选项: $1"
            show_usage
            exit 1
            ;;
        *)
            break
            ;;
    esac
done

# 检查必需参数
if [ $# -lt 3 ]; then
    log_error "缺少必需参数"
    show_usage
    exit 1
fi

REMOTE_HOST="$1"
REMOTE_USER="$2"
REMOTE_PATH="$3"

# 处理密码认证
if [ "$USE_PASSWORD" = true ] && [ -n "$PASSWORD_FILE" ]; then
    log_error "不能同时使用 -P 和 --password-file 选项"
    exit 1
fi

if [ "$USE_PASSWORD" = true ] && [ -n "$SSH_KEY" ]; then
    log_error "不能同时使用密码认证和密钥认证"
    exit 1
fi

# 获取SSH密码
if [ "$USE_PASSWORD" = true ]; then
    # 首先检查是否通过环境变量传递了密码
    if [ -n "$SSH_PASS" ]; then
        SSH_PASSWORD="$SSH_PASS"
        log_info "使用环境变量中的SSH密码"
    else
        echo -n "请输入SSH密码: "
        read -s SSH_PASSWORD
        echo
        if [ -z "$SSH_PASSWORD" ]; then
            log_error "密码不能为空"
            exit 1
        fi
    fi
elif [ -n "$PASSWORD_FILE" ]; then
    if [ ! -f "$PASSWORD_FILE" ]; then
        log_error "密码文件不存在: $PASSWORD_FILE"
        exit 1
    fi
    SSH_PASSWORD=$(cat "$PASSWORD_FILE")
    if [ -z "$SSH_PASSWORD" ]; then
        log_error "密码文件为空: $PASSWORD_FILE"
        exit 1
    fi
fi

# 构建SSH选项
SSH_OPTS="-p ${SSH_PORT}"
SCP_OPTS="-P ${SSH_PORT}"

if [ -n "$SSH_KEY" ]; then
    SSH_OPTS="$SSH_OPTS -i $SSH_KEY"
    SCP_OPTS="$SCP_OPTS -i $SSH_KEY"
elif [ -n "$SSH_PASSWORD" ]; then
    # 使用sshpass进行密码认证
    SSH_CMD="sshpass -p '$SSH_PASSWORD' ssh"
    SCP_CMD="sshpass -p '$SSH_PASSWORD' scp"
else
    SSH_CMD="ssh"
    SCP_CMD="scp"
fi

# 检查依赖
check_dependencies() {
    log_step "检查系统依赖..."
    
    # 检查基本命令
    local missing_commands=()
    
    if ! command -v docker &> /dev/null; then
        missing_commands+=("docker")
    fi
    
    if ! command -v scp &> /dev/null; then
        missing_commands+=("scp")
    fi
    
    if ! command -v ssh &> /dev/null; then
        missing_commands+=("ssh")
    fi
    
    # 检查date命令（使用更灵活的方式）
    local date_found=false
    for cmd in "date" "/bin/date" "/usr/bin/date" "/usr/sbin/date"; do
        if command -v "$cmd" &> /dev/null; then
            date_found=true
            break
        fi
    done
    
    if [ "$date_found" = false ]; then
        missing_commands+=("date")
    fi
    
    # 报告缺失的命令
    if [ ${#missing_commands[@]} -gt 0 ]; then
        log_error "缺少以下必需的命令:"
        for cmd in "${missing_commands[@]}"; do
            echo "  - $cmd"
        done
        log_error "请安装缺少的软件包"
        exit 1
    fi
    
    # 检查是否需要sshpass
    if [ -n "$SSH_PASSWORD" ]; then
        if ! command -v sshpass &> /dev/null; then
            log_error "使用密码认证需要安装sshpass"
            log_error "请安装sshpass:"
            log_error "  Ubuntu/Debian: sudo apt install sshpass"
            log_error "  CentOS/RHEL:   sudo dnf install sshpass"
            log_error "  macOS:         brew install sshpass"
            exit 1
        fi
    fi
    
    log_info "所有依赖检查完成"
}

# 检查远程服务器连接
check_remote_connection() {
    log_step "检查远程服务器连接..."
    
    if [ "$DRY_RUN" = true ]; then
        log_info "[DRY-RUN] 跳过远程连接检查"
        return 0
    fi
    
    # 对于密码认证，不能使用BatchMode
    local connect_opts="$SSH_OPTS -o ConnectTimeout=10"
    if [ -z "$SSH_PASSWORD" ]; then
        connect_opts="$connect_opts -o BatchMode=yes"
    fi
    
    if eval "$SSH_CMD $connect_opts \"$REMOTE_USER@$REMOTE_HOST\" \"echo 'Connection test successful'\"" &>/dev/null; then
        log_info "远程服务器连接成功"
    else
        log_error "无法连接到远程服务器 $REMOTE_USER@$REMOTE_HOST"
        log_error "请检查:"
        log_error "  1. 服务器地址是否正确"
        log_error "  2. SSH服务是否运行"
        log_error "  3. 用户名和认证是否正确"
        log_error "  4. 网络连接是否正常"
        exit 1
    fi
}

# 检查镜像是否存在
check_images() {
    log_step "检查Docker镜像..."
    
    local missing_images=()
    
    for image in "${IMAGES[@]}"; do
        if docker image inspect "$image" &>/dev/null; then
            log_info "镜像存在: $image"
        else
            log_warn "镜像不存在: $image"
            missing_images+=("$image")
        fi
    done
    
    if [ ${#missing_images[@]} -gt 0 ]; then
        log_warn "以下镜像不存在，将尝试拉取:"
        for image in "${missing_images[@]}"; do
            echo "  - $image"
        done
        
        read -p "是否要拉取缺失的镜像? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            pull_missing_images "${missing_images[@]}"
        else
            log_error "无法继续，缺少必要的镜像"
            exit 1
        fi
    fi
}

# 拉取缺失的镜像
pull_missing_images() {
    local images=("$@")
    
    log_step "拉取缺失的Docker镜像..."
    
    for image in "${images[@]}"; do
        # 跳过本地构建的镜像
        if [[ "$image" == *"goblog"* ]]; then
            log_warn "跳过本地构建镜像: $image (请先运行 docker-compose build)"
            continue
        fi
        
        log_info "拉取镜像: $image"
        if [ "$DRY_RUN" = false ]; then
            docker pull "$image"
        else
            log_info "[DRY-RUN] docker pull $image"
        fi
    done
}

# 导出Docker镜像
export_images() {
    log_step "导出Docker镜像..."
    
    # 创建导出目录
    if [ "$DRY_RUN" = false ]; then
        mkdir -p "$EXPORT_DIR"
    else
        log_info "[DRY-RUN] mkdir -p $EXPORT_DIR"
    fi
    
    local exported_files=()
    
    for image in "${IMAGES[@]}"; do
        # 检查镜像是否存在
        if [ "$DRY_RUN" = false ] && ! docker image inspect "$image" &>/dev/null; then
            log_warn "跳过不存在的镜像: $image"
            continue
        fi
        
        # 生成文件名 (替换特殊字符)
        local filename=$(echo "$image" | sed 's/[\/:]/_/g').tar
        local filepath="$EXPORT_DIR/$filename"
        
        log_info "导出镜像: $image -> $filename"
        
        if [ "$DRY_RUN" = false ]; then
            docker save "$image" -o "$filepath"
            exported_files+=("$filepath")
        else
            log_info "[DRY-RUN] docker save $image -o $filepath"
            exported_files+=("$filepath")
        fi
    done
    
    # 压缩镜像文件
    if [ "$COMPRESS" = true ]; then
        log_info "压缩镜像文件..."
        
        if [ "$DRY_RUN" = false ]; then
            tar -czf "$ARCHIVE_NAME" -C "$EXPORT_DIR" .
            log_info "压缩完成: $ARCHIVE_NAME"
            
            # 显示文件大小
            local size=$(du -h "$ARCHIVE_NAME" | cut -f1)
            log_info "压缩包大小: $size"
        else
            log_info "[DRY-RUN] tar -czf $ARCHIVE_NAME -C $EXPORT_DIR ."
        fi
    fi
}

# 上传到远程服务器
upload_to_remote() {
    log_step "上传镜像到远程服务器..."
    
    # 确保远程目录存在
    log_info "创建远程目录..."
    if [ "$DRY_RUN" = false ]; then
        eval "$SSH_CMD $SSH_OPTS \"$REMOTE_USER@$REMOTE_HOST\" \"mkdir -p $REMOTE_PATH\""
    else
        log_info "[DRY-RUN] $SSH_CMD $SSH_OPTS $REMOTE_USER@$REMOTE_HOST \"mkdir -p $REMOTE_PATH\""
    fi
    
    if [ "$COMPRESS" = true ]; then
        # 上传压缩包
        local file_to_upload="$ARCHIVE_NAME"
        log_info "上传压缩包: $file_to_upload"
        
        if [ "$DRY_RUN" = false ]; then
            eval "$SCP_CMD $SCP_OPTS \"$file_to_upload\" \"$REMOTE_USER@$REMOTE_HOST:$REMOTE_PATH\""
            
            # 在远程服务器上解压
            log_info "在远程服务器上解压..."
            eval "$SSH_CMD $SSH_OPTS \"$REMOTE_USER@$REMOTE_HOST\" \"cd $REMOTE_PATH && tar -xzf $ARCHIVE_NAME\""
        else
            log_info "[DRY-RUN] $SCP_CMD $SCP_OPTS $file_to_upload $REMOTE_USER@$REMOTE_HOST:$REMOTE_PATH"
            log_info "[DRY-RUN] $SSH_CMD $SSH_OPTS $REMOTE_USER@$REMOTE_HOST \"cd $REMOTE_PATH && tar -xzf $ARCHIVE_NAME\""
        fi
    else
        # 上传单个文件
        log_info "上传镜像文件..."
        if [ "$DRY_RUN" = false ]; then
            eval "$SCP_CMD $SCP_OPTS \"$EXPORT_DIR\"/*.tar \"$REMOTE_USER@$REMOTE_HOST:$REMOTE_PATH\""
        else
            log_info "[DRY-RUN] $SCP_CMD $SCP_OPTS $EXPORT_DIR/*.tar $REMOTE_USER@$REMOTE_HOST:$REMOTE_PATH"
        fi
    fi
    
    log_info "上传完成！"
}

# 创建远程导入脚本
create_import_script() {
    log_step "创建远程导入脚本..."
    
    local import_script="import_images.sh"
    
    if [ "$DRY_RUN" = false ]; then
        cat > "$import_script" << 'EOF'
#!/bin/bash

# Docker镜像导入脚本
# 在远程服务器上运行此脚本来导入镜像

echo "开始导入Docker镜像..."

# 导入所有tar文件
for tar_file in *.tar; do
    if [ -f "$tar_file" ]; then
        echo "导入镜像: $tar_file"
        docker load -i "$tar_file"
    fi
done

echo "镜像导入完成！"

# 显示导入的镜像
echo ""
echo "已导入的镜像列表:"
docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}"

# 清理tar文件（可选）
read -p "是否删除tar文件? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    rm -f *.tar
    echo "tar文件已删除"
fi
EOF
        
        # 上传导入脚本
        log_info "上传导入脚本..."
        eval "$SCP_CMD $SCP_OPTS \"$import_script\" \"$REMOTE_USER@$REMOTE_HOST:$REMOTE_PATH\""
        
        # 给脚本添加执行权限
        eval "$SSH_CMD $SSH_OPTS \"$REMOTE_USER@$REMOTE_HOST\" \"chmod +x $REMOTE_PATH/$import_script\""
        
        # 清理本地脚本
        rm -f "$import_script"
    else
        log_info "[DRY-RUN] 创建并上传导入脚本到远程服务器"
    fi
}

# 清理本地文件
cleanup_local() {
    if [ "$KEEP_LOCAL" = false ]; then
        log_step "清理本地文件..."
        
        if [ "$DRY_RUN" = false ]; then
            rm -rf "$EXPORT_DIR"
            if [ "$COMPRESS" = true ]; then
                rm -f "$ARCHIVE_NAME"
            fi
            log_info "本地文件已清理"
        else
            log_info "[DRY-RUN] 清理本地导出文件"
        fi
    else
        log_info "保留本地文件: $EXPORT_DIR"
        if [ "$COMPRESS" = true ]; then
            log_info "保留压缩包: $ARCHIVE_NAME"
        fi
    fi
}

# 显示总结信息
show_summary() {
    log_step "操作总结"
    
    echo ""
    log_info "镜像导出和上传完成！"
    echo ""
    echo "📁 远程服务器信息:"
    echo "  - 服务器: $REMOTE_USER@$REMOTE_HOST:$SSH_PORT"
    echo "  - 路径: $REMOTE_PATH"
    echo ""
    echo "🐳 导出的镜像:"
    for image in "${IMAGES[@]}"; do
        echo "  - $image"
    done
    echo ""
    echo "📋 在远程服务器上导入镜像:"
    if [ -n "$SSH_PASSWORD" ]; then
        echo "  $SSH_CMD $SSH_OPTS $REMOTE_USER@$REMOTE_HOST"
    else
        echo "  ssh $SSH_OPTS $REMOTE_USER@$REMOTE_HOST"
    fi
    echo "  cd $REMOTE_PATH"
    echo "  ./import_images.sh"
    echo ""
    echo "🚀 启动服务:"
    echo "  docker-compose up -d"
    echo ""
}

# 主函数
main() {
    echo "🐳 Docker镜像导出和上传工具"
    echo "=============================="
    echo ""
    
    log_info "目标服务器: $REMOTE_USER@$REMOTE_HOST:$SSH_PORT"
    log_info "目标路径: $REMOTE_PATH"
    
    if [ "$DRY_RUN" = true ]; then
        log_warn "运行在 DRY-RUN 模式，只显示操作不实际执行"
    fi
    
    echo ""
    
    check_dependencies
    check_remote_connection
    check_images
    export_images
    upload_to_remote
    create_import_script
    cleanup_local
    show_summary
    
    log_info "所有操作完成！"
}

# 错误处理
trap 'log_error "脚本执行失败，请检查错误信息"; exit 1' ERR

# 运行主函数
main "$@" 