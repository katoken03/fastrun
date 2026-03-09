# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=f
PREFIX=/usr/local

# Build-time variables
VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

.PHONY: all build clean test install uninstall

# ビルド(ファイル名 f でビルドします)
build:
	$(GOBUILD) -o $(BINARY_NAME) -v

# クリーンアップ
clean:
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).test

# クリーンアップ + ビルド
all: clean build

# テスト実行
test:
	$(GOTEST) -v ./...

# インストール（システムへの配置）
install: build
	mkdir -p $(PREFIX)/bin
	cp $(BINARY_NAME) $(PREFIX)/bin/$(BINARY_NAME)
	# 設定ディレクトリの作成
	mkdir -p $(HOME)/.config/fastrun

# アンインストール
uninstall:
	rm -f $(PREFIX)/bin/$(BINARY_NAME)

# 開発用のセットアップ（依存関係の解決）
dev-setup:
	$(GOCMD) mod tidy
	$(GOCMD) mod verify

# バージョン情報の表示
version:
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
