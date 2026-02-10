.PHONY: help run build test clean migrate lint docker-build docker-run

# 默认目标
help:
	@echo "可用命令:"
	@echo "  make run       使用 air 热重载运行应用"
	@echo "  make run-simple 直接运行应用"
	@echo "  make build     构建应用"
	@echo "  make test      运行测试"
	@echo "  make clean     清理构建文件"
	@echo "  make lint      运行代码检查"
	@echo "  make docker-build 构建 Docker 镜像"
	@echo "  make docker-run   使用 Docker 运行应用"

# 使用 air 热重载运行
run:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air 未安装，请先安装: go install github.com/cosmtrek/air@latest"; \
		echo "或者运行: make run-simple"; \
	fi

# 直接运行
run-simple:
	go run main.go

# 构建应用
build:
	go build -o myblog-backend main.go

# 运行测试
test:
	go test ./... -v

# 清理
clean:
	rm -f myblog-backend
	rm -f tmp/main
	rm -f *.db
	rm -rf tmp/
	rm -f *.log

# 代码检查
lint:
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint 未安装，跳过检查"; \
	fi

# 初始化数据库（如果需要）
migrate:
	go run main.go --migrate

# Docker 构建
docker-build:
	docker build -t myblog-backend .

# Docker 运行
docker-run:
	docker run -p 8080:8080 --env-file .env myblog-backend