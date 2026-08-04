#!/bin/bash

# iSCSI Web Panel 一键安装脚本
# 支持 Ubuntu/Debian 和 CentOS/RHEL

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查 root 权限
check_root() {
    if [ "$EUID" -ne 0 ]; then
        echo -e "${RED}错误: 请使用 sudo 运行此脚本${NC}"
        exit 1
    fi
}

# 检测操作系统
detect_os() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS=$NAME
        VER=$VERSION_ID
    else
        echo -e "${RED}错误: 无法检测操作系统${NC}"
        exit 1
    fi
}

# 安装依赖
install_dependencies() {
    echo -e "${GREEN}=== 安装依赖 ===${NC}"
    
    case "$OS" in
        *Ubuntu*|*Debian*)
            apt-get update
            apt-get install -y \
                tgt \
                sqlite3 \
                wget \
                curl \
                git \
                build-essential
            
            # 启动 tgt 服务
            systemctl enable tgt
            systemctl start tgt
            ;;
        *CentOS*|*Red\ Hat*|*Fedora*)
            yum install -y epel-release
            yum install -y \
                scsi-target-utils \
                sqlite \
                wget \
                curl \
                git \
                gcc \
                make
            
            # 启动 tgtd 服务
            systemctl enable tgtd
            systemctl start tgtd
            ;;
        *)
            echo -e "${RED}不支持的操作系统: $OS${NC}"
            exit 1
            ;;
    esac
    
    echo -e "${GREEN}✓ 依赖安装完成${NC}"
}

# 安装 Go
install_go() {
    if command -v go &> /dev/null; then
        echo -e "${GREEN}✓ Go 已安装: $(go version)${NC}"
        return
    fi
    
    echo -e "${GREEN}=== 安装 Go ===${NC}"
    
    GO_VERSION="1.21.5"
    ARCH=$(uname -m)
    
    case "$ARCH" in
        x86_64)
            GO_ARCH="amd64"
            ;;
        aarch64|arm64)
            GO_ARCH="arm64"
            ;;
        *)
            echo -e "${RED}不支持的架构: $ARCH${NC}"
            exit 1
            ;;
    esac
    
    cd /tmp
    wget -q "https://go.dev/dl/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
    rm -rf /usr/local/go
    tar -C /usr/local -xzf "go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
    rm "go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
    
    # 设置环境变量
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile.d/go.sh
    export PATH=$PATH:/usr/local/go/bin
    
    echo -e "${GREEN}✓ Go 安装完成${NC}"
}

# 编译项目
build_project() {
    echo -e "${GREEN}=== 编译项目 ===${NC}"
    
    INSTALL_DIR="/opt/iscsi-web-panel"
    
    # 创建安装目录
    mkdir -p "$INSTALL_DIR"
    
    # 复制文件
    cp -r . "$INSTALL_DIR/"
    cd "$INSTALL_DIR"
    
    # 下载 Go 依赖
    export PATH=$PATH:/usr/local/go/bin
    go mod download
    
    # 编译
    go build -o iscsi-web-panel -ldflags="-s -w" main.go
    
    # 创建数据目录
    mkdir -p data
    
    # 设置权限
    chmod +x iscsi-web-panel
    
    echo -e "${GREEN}✓ 编译完成${NC}"
}

# 创建 systemd 服务
create_service() {
    echo -e "${GREEN}=== 创建系统服务 ===${NC}"
    
    cat > /etc/systemd/system/iscsi-web-panel.service << 'EOF'
[Unit]
Description=iSCSI Web Panel
After=network.target tgt.service
Wants=tgt.service

[Service]
Type=simple
User=root
WorkingDirectory=/opt/iscsi-web-panel
ExecStart=/opt/iscsi-web-panel/iscsi-web-panel
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

# 环境变量
Environment=LISTEN_ADDR=:3005
Environment=DATA_DIR=/opt/iscsi-web-panel/data
Environment=DB_PATH=/opt/iscsi-web-panel/data/iscsi.db

# 安全设置
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/iscsi-web-panel/data

[Install]
WantedBy=multi-user.target
EOF

    # 重新加载 systemd
    systemctl daemon-reload
    
    # 启用并启动服务
    systemctl enable iscsi-web-panel
    systemctl start iscsi-web-panel
    
    echo -e "${GREEN}✓ 服务创建完成${NC}"
}

# 配置防火墙
configure_firewall() {
    echo -e "${GREEN}=== 配置防火墙 ===${NC}"
    
    # UFW (Ubuntu/Debian)
    if command -v ufw &> /dev/null; then
        ufw allow 3005/tcp
        echo -e "${GREEN}✓ UFW 规则已添加${NC}"
    fi
    
    # firewalld (CentOS/RHEL)
    if command -v firewall-cmd &> /dev/null; then
        firewall-cmd --permanent --add-port=3005/tcp
        firewall-cmd --reload
        echo -e "${GREEN}✓ firewalld 规则已添加${NC}"
    fi
    
    # iptables (通用)
    if command -v iptables &> /dev/null && ! command -v ufw &> /dev/null && ! command -v firewall-cmd &> /dev/null; then
        iptables -I INPUT -p tcp --dport 3005 -j ACCEPT
        if command -v iptables-save &> /dev/null; then
            iptables-save > /etc/iptables.rules
        fi
        echo -e "${GREEN}✓ iptables 规则已添加${NC}"
    fi
}

# 显示安装信息
show_info() {
    IP=$(hostname -I | awk '{print $1}')
    
    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}      iSCSI Web Panel 安装完成！${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo -e "访问地址: ${YELLOW}http://${IP}:3005${NC}"
    echo ""
    echo -e "默认账号:"
    echo -e "  用户名: ${YELLOW}admin${NC}"
    echo -e "  密码: ${YELLOW}admin123${NC}"
    echo ""
    echo -e "${RED}⚠️  重要: 请立即修改默认密码！${NC}"
    echo ""
    echo -e "服务管理命令:"
    echo -e "  查看状态: ${YELLOW}systemctl status iscsi-web-panel${NC}"
    echo -e "  重启服务: ${YELLOW}systemctl restart iscsi-web-panel${NC}"
    echo -e "  查看日志: ${YELLOW}journalctl -u iscsi-web-panel -f${NC}"
    echo ""
}

# 主函数
main() {
    echo -e "${GREEN}=== iSCSI Web Panel 一键安装 ===${NC}"
    echo ""
    
    check_root
    detect_os
    
    echo -e "操作系统: ${YELLOW}$OS${NC}"
    echo ""
    
    install_dependencies
    install_go
    build_project
    create_service
    configure_firewall
    show_info
}

# 运行主函数
main "$@"
