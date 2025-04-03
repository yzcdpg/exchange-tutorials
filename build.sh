#!/bin/bash

# 脚本名称: build.sh
# 用途: 编译 Go 实现的加密货币交易所项目

# 定义颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # 无颜色

# 项目相关变量
PROJECT_NAME="crypto-exchange"
BINARY_NAME="exchange"
MAIN_FILE="main.go"
GO_VERSION="1.21"  # 假设使用 Go 1.21
OUTPUT_DIR="./bin"

# 检查命令是否存在
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# 检查和安装 Go
check_go() {
    if ! command_exists go; then
        echo -e "${RED}Go 未安装，正在安装 Go ${GO_VERSION}...${NC}"
        if [[ "$OSTYPE" == "darwin"* ]]; then
            brew install go@${GO_VERSION}
        elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
            wget https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz
            sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
            echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
            source ~/.bashrc
            rm go${GO_VERSION}.linux-amd64.tar.gz
        else
            echo -e "${RED}不支持的操作系统，请手动安装 Go${NC}"
            exit 1
        fi
    fi
    echo -e "${GREEN}Go 已安装: $(go version)${NC}"
}

# 初始化项目结构
init_project() {
    if [ ! -d "$PROJECT_NAME" ]; then
        echo -e "${GREEN}初始化项目结构...${NC}"
        mkdir -p "$PROJECT_NAME/$OUTPUT_DIR"
        cd "$PROJECT_NAME" || exit
        go mod init "$PROJECT_NAME"
    else
        cd "$PROJECT_NAME" || exit
    fi
}

# 安装依赖
install_dependencies() {
    echo -e "${GREEN}安装项目依赖...${NC}"
    go get github.com/gorilla/mux
    go get github.com/gorilla/websocket
    go get github.com/dgrijalva/jwt-go
    go get github.com/lib/pq
    go get golang.org/x/time/rate
    go mod tidy
}

# 编译项目
build_project() {
    echo -e "${GREEN}编译项目...${NC}"
    if [ ! -f "$MAIN_FILE" ]; then
        echo -e "${RED}错误: 未找到 $MAIN_FILE${NC}"
        exit 1
    fi
    go build -o "$OUTPUT_DIR/$BINARY_NAME" "$MAIN_FILE"
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}编译成功！二进制文件位于: $OUTPUT_DIR/$BINARY_NAME${NC}"
    else
        echo -e "${RED}编译失败${NC}"
        exit 1
    fi
}

# 初始化数据库（可选）
init_database() {
    echo -e "${GREEN}初始化 PostgreSQL 数据库...${NC}"
    if ! command_exists psql; then
        echo -e "${RED}PostgreSQL 未安装，请手动安装${NC}"
        return
    fi
    psql -U postgres -c "CREATE DATABASE exchange;" 2>/dev/null
    psql -U postgres -d exchange -f ../db_schema.sql 2>/dev/null
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}数据库初始化完成${NC}"
    else
        echo -e "${RED}数据库初始化失败，请检查配置${NC}"
    fi
}

# 主函数
main() {
    echo -e "${GREEN}开始构建 $PROJECT_NAME 项目...${NC}"

    # 检查和安装 Go
    check_go

    # 初始化项目结构
    init_project

    # 安装依赖
    install_dependencies

    # 编译项目
    build_project

    # 可选：初始化数据库
    read -p "是否初始化数据库？(y/n): " init_db
    if [ "$init_db" == "y" ]; then
        init_database
    fi

    echo -e "${GREEN}构建完成！${NC}"
}

# 执行主函数
main
