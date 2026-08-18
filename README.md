<div align="center">
  <a href="https://www.dongpolakeside.com">
    <img alt="东坡湖畔" src="docs/img/logo.svg" height="99px">
  </a>

  # 东坡湖畔

  面向海南大学校友和在校学生的非官方交流论坛

  **访问地址：[www.dongpolakeside.com](https://www.dongpolakeside.com)**
</div>

## 项目简介

东坡湖畔是一个服务于海南大学校友和在校学生的在线交流社区，旨在为大家提供一个分享校园资讯、交流学习经验、互助答疑和保持校友联系的平台。

本项目基于开源问答社区项目 [Apache Answer](https://github.com/apache/answer) 进行二次开发，并结合社区需求对功能、界面及部署方式进行调整。感谢 Apache Answer 社区及所有上游贡献者提供的优秀开源项目。

## 重要声明

- 本项目及东坡湖畔网站均为**非官方论坛**，由社区独立开发和维护。
- 本项目与海南大学及其下属院系、部门、学生组织或校友组织不存在隶属、授权或官方合作关系。
- 网站中的用户观点及内容仅代表发布者本人，不代表海南大学或项目维护者的立场。
- “海南大学”等名称仅用于说明本社区所服务和讨论的对象，不表示任何官方背书。

## 技术基础

项目沿用 Apache Answer 的主要技术架构：

- 后端：Go
- 前端：React
- 插件系统：Apache Answer Plugin System

上游项目文档：[answer.apache.org](https://answer.apache.org)

## 本地开发

本项目不使用上游 `apache/answer` 镜像进行本地开发。准备好
`answer-data/conf/config.yaml` 中的本地 MySQL 配置后，启动后端：

```bash
CACHE_TYPE=memory go run ./cmd/answer upgrade -C "$PWD/answer-data"
CACHE_TYPE=memory go run ./cmd/answer run -C "$PWD/answer-data"
```

在另一个终端启动前端：

```bash
cd ui
pnpm install
pnpm start
```

## 生产部署与更新

生产环境使用 GitHub Actions 构建的 GHCR 镜像和独立的 Docker Compose 配置。服务器更新、数据备份、运行状态检查和应用回滚方法，请参阅 [生产服务器更新手册](deploy/production/README.md)。

> 请勿在生产服务器的仓库根目录直接运行 `docker compose`。生产环境必须使用 `deploy/production/docker-compose.yml` 和对应的 `.env` 文件。

## 从源码构建

### 环境要求

- Golang >= 1.23
- Node.js >= 20
- pnpm >= 9
- [mockgen](https://github.com/uber-go/mock) >= 0.6.0
- [wire](https://github.com/google/wire/) >= 0.5.0

### 构建命令

```bash
# 生成构建所需代码；可先运行 make check 检查相关工具是否已安装
make generate

# 安装前端依赖并构建前端
make ui

# 安装后端依赖并构建项目
make build
```

## 开源许可

本项目的二次开发与分发遵循仓库中的 [Apache License 2.0](LICENSE)。Apache Answer 的名称、商标及相关权益归其各自权利人所有。
