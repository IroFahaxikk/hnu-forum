# 生产服务器更新手册

本文档适用于东坡湖畔生产服务器，部署目录为：

```text
/opt/hnu-forum/deploy/production
```

生产环境使用以下文件：

```text
/opt/hnu-forum/deploy/production/docker-compose.yml
/opt/hnu-forum/deploy/production/.env
```

## 重要警告

不要在 `/opt/hnu-forum` 仓库根目录直接执行 `docker compose up`、`docker compose pull` 或 `docker compose stop`。仓库不再保留上游的根目录 Compose 配置，生产环境只使用 `deploy/production/docker-compose.yml`。

生产环境的所有 Compose 命令都应当在 `/opt/hnu-forum/deploy/production` 目录执行，并显式指定环境文件和配置文件：

```bash
docker compose --env-file .env -f docker-compose.yml <命令>
```

不要执行 `docker compose down -v`，不要删除 `/srv/hnu-forum` 或 `/var/lib/docker/volumes` 下的文件。

## GitHub 构建说明

代码合并到 `main` 后，GitHub Actions 会构建并推送以下镜像：

```text
ghcr.io/irofahaxikk/hnu-forum:production
ghcr.io/irofahaxikk/hnu-forum:sha-<完整提交哈希>
```

只有在 GitHub Actions 中的 `Build Production Image` 工作流执行成功后，才能开始更新服务器。

## 标准更新流程

### 1. 登录并进入生产目录

```bash
ssh root@149.88.72.99
cd /opt/hnu-forum
git status --short
git pull --ff-only origin main
cd /opt/hnu-forum/deploy/production
```

如果 `git status --short` 显示服务器存在未提交的代码修改，应先停止更新并确认修改来源，不要使用 `git reset --hard` 强行覆盖。同步仓库可以确保服务器同时取得本次 PR 对生产 Compose 等部署文件的修改；应用程序本身仍使用 GitHub Actions 构建的镜像。

确认当前目录和生产文件：

```bash
pwd
ls -la .env docker-compose.yml
```

`pwd` 必须输出：

```text
/opt/hnu-forum/deploy/production
```

### 2. 检查更新前状态

```bash
docker compose --env-file .env -f docker-compose.yml ps
```

`answer`、`mysql` 和 `redis` 应当处于运行或健康状态。

记录当前镜像对应的 Git 提交，便于回滚：

```bash
CURRENT_IMAGE_ID=$(docker inspect hnu-forum-answer-1 --format '{{.Image}}')
PREVIOUS_REVISION=$(docker image inspect "$CURRENT_IMAGE_ID" --format '{{index .Config.Labels "org.opencontainers.image.revision"}}')
echo "更新前版本：$PREVIOUS_REVISION"
```

请将输出的提交哈希保存在本次更新记录中。

### 3. 备份数据库

```bash
mkdir -p /root/hnu-emergency-backups
BACKUP_FILE="/root/hnu-emergency-backups/answer-$(date +%Y%m%d-%H%M%S).sql"

docker exec hnu-forum-mysql-1 sh -lc \
  'exec mysqldump --no-tablespaces -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" --single-transaction --routines --triggers "$MYSQL_DATABASE"' \
  > "$BACKUP_FILE"

test -s "$BACKUP_FILE" \
  && grep -q 'CREATE TABLE' "$BACKUP_FILE" \
  && echo "数据库备份成功：$BACKUP_FILE"
```

出现 `Using a password on the command line interface can be insecure` 警告是 MySQL 客户端的常规提示。只要备份文件存在且不为空，即可继续。

### 4. 拉取并更新应用容器

```bash
docker compose --env-file .env -f docker-compose.yml pull answer
docker compose --env-file .env -f docker-compose.yml up -d --no-deps answer
```

`--no-deps` 表示只重建 `answer`，不会重建 MySQL 和 Redis。

### 5. 检查容器和数据挂载

```bash
docker compose --env-file .env -f docker-compose.yml ps

docker inspect hnu-forum-answer-1 \
  --format '镜像={{.Config.Image}}{{println}}{{range .Mounts}}挂载={{.Source}} -> {{.Destination}}{{println}}{{end}}'
```

结果必须包含：

```text
镜像=ghcr.io/irofahaxikk/hnu-forum:production
挂载=/srv/hnu-forum/answer -> /data
```

如果镜像不是 `ghcr.io/irofahaxikk/hnu-forum:production`，或者挂载不是 `/srv/hnu-forum/answer -> /data`，应立即停止更新并检查执行目录和 Compose 参数。

### 6. 检查日志和服务

```bash
docker compose --env-file .env -f docker-compose.yml logs --tail=100 answer
curl -I http://127.0.0.1:9080
```

正常日志应包含以下信息，且不应持续出现数据库或 Redis 连接错误：

```text
upgrade done
config file path: /data/conf/config.yaml
Answer is starting
```

随后访问 [https://www.dongpolakeside.com](https://www.dongpolakeside.com)，检查以下内容：

- 首页能够正常打开；
- 原有帖子、回答和评论存在；
- 登录、发帖和上传功能正常；
- 浏览器强制刷新后使用的是新版本静态资源。

浏览器强制刷新快捷键：Windows/Linux 使用 `Ctrl + F5`，macOS 使用 `Command + Shift + R`。

## 查看当前部署版本

```bash
CURRENT_IMAGE_ID=$(docker inspect hnu-forum-answer-1 --format '{{.Image}}')
docker image inspect "$CURRENT_IMAGE_ID" \
  --format '镜像={{.Id}}{{println}}Git提交={{index .Config.Labels "org.opencontainers.image.revision"}}{{println}}构建时间={{.Created}}'
```

将输出的 Git 提交与 GitHub Actions 对应工作流中的提交进行核对。

## 应用镜像回滚

如果新版本无法正常工作，并且已经记录了更新前的完整提交哈希，可以使用 GitHub Actions 推送的不可变 SHA 标签回滚。

将下面变量的值替换为更新前记录的完整提交哈希：

```bash
cd /opt/hnu-forum/deploy/production
PREVIOUS_REVISION="更新前完整提交哈希"

IMAGE_TAG="sha-${PREVIOUS_REVISION}" \
  docker compose --env-file .env -f docker-compose.yml pull answer

IMAGE_TAG="sha-${PREVIOUS_REVISION}" \
  docker compose --env-file .env -f docker-compose.yml up -d --no-deps answer
```

回滚后重新执行容器状态、数据挂载、日志和网页检查。

应用镜像回滚不会自动回滚数据库。如果新版本已经执行数据库结构迁移，不要直接覆盖或导入数据库，应先停止 `answer`，保留现场并根据具体迁移内容制定恢复方案。

## 常用只读检查

查看数据库中的内容数量：

```bash
docker exec hnu-forum-mysql-1 sh -lc \
  'exec mysql -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE" -e "SELECT COUNT(*) AS questions FROM question; SELECT COUNT(*) AS answers FROM answer; SELECT COUNT(*) AS comments FROM comment;"'
```

查看生产数据目录但不修改内容：

```bash
find /srv/hnu-forum/answer -maxdepth 3 -type f -printf '%p | %s bytes | %TY-%Tm-%Td %TH:%TM:%TS\n'
```

查看 Docker 卷：

```bash
docker volume ls
```
