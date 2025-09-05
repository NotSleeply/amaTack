# 编译环境变量
GO=go
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOTEST=$(GO) test
GOGET=$(GO) get

# 项目名称和输出目录
BINARY_NAME=tamazon.exe
BIN_DIR=bin

# 其他可执行文件名称
QTACK_NAME=QTack2.0.exe
STACK_NAME=STack2.0.exe
MTACK_NAME=MTack3.0.exe

# 默认目标
all: clean build

# 构建主程序
build:
	$(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME) .

# 构建快速版本 - 搜索深度2跳3
qtack:
	$(GOBUILD) -ldflags "-X main.searchMode=qtack" -o $(BIN_DIR)/$(QTACK_NAME) ./main.go

# 构建慢速版本 - 搜索深度2跳4
stack:
	$(GOBUILD) -ldflags "-X main.searchMode=stack" -o $(BIN_DIR)/$(STACK_NAME) ./main.go

# 构建最新版本MTack3.0
mtack:
	$(GOBUILD) -ldflags "-X main.searchMode=mtack" -o $(BIN_DIR)/$(MTACK_NAME) ./main.go

# 构建所有版本
all-versions: clean build qtack stack mtack

# 测试
test:
	$(GOTEST) -v ./...

# 清理编译产物
clean:
	$(GOCLEAN)
	rm -f $(BIN_DIR)/$(BINARY_NAME)
	rm -f $(BIN_DIR)/$(QTACK_NAME)
	rm -f $(BIN_DIR)/$(STACK_NAME)
	rm -f $(BIN_DIR)/$(MTACK_NAME)

# 安装依赖
deps:
	$(GOGET) -v ./...

# 运行程序
run:
	$(GO) run .

# 显示帮助信息
help:
	@echo "以下是可用的make命令:"
	@echo "  make build      - 构建标准版本"
	@echo "  make qtack      - 构建快速版本 (QTack2.0)"
	@echo "  make stack      - 构建慢速版本 (STack2.0)"
	@echo "  make mtack      - 构建新版本 (MTack3.0)"
	@echo "  make all        - 清理并构建标准版本"
	@echo "  make all-versions - 构建所有版本"
	@echo "  make clean      - 清理编译产物"
	@echo "  make test       - 运行测试"
	@echo "  make deps       - 安装依赖"
	@echo "  make run        - 直接运行程序"

.PHONY: build qtack stack mtack all all-versions clean test deps run help