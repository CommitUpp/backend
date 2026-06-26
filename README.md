# Sheater Backend

## Backend

<p>
  <img src="https://skillicons.dev/icons?i=go" height="50" alt="Go" />
  <img src="https://skillicons.dev/icons?i=docker" height="50" alt="Docker" />
  <img src="https://skillicons.dev/icons?i=redis" height="50" alt="Redis" />
  <img src="https://skillicons.dev/icons?i=supabase" height="50" alt="Supabase" />
  <img src="https://skillicons.dev/icons?i=prometheus" height="50" alt="Prometheus" />
  <img src="https://skillicons.dev/icons?i=grafana" height="50" alt="Grafana" />
  <img src="https://skillicons.dev/icons?i=kubernetes" height="50" alt="Kubernetes" />
  <img src="https://skillicons.dev/icons?i=githubactions" height="50" alt="GitHub Actions" />
</p>

<p>
  <img src="https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white" alt="Go 1.26" />
  <img src="https://img.shields.io/badge/Echo-v4-00ADD8" alt="Echo v4" />
  <img src="https://img.shields.io/badge/gRPC-Protocol%20Buffers-244C5A" alt="gRPC and Protocol Buffers" />
  <img src="https://img.shields.io/badge/OpenAPI-oapi--codegen-6BA539" alt="OpenAPI oapi-codegen" />
  <img src="https://img.shields.io/badge/GHCR-Container%20Registry-181717?logo=github" alt="GitHub Container Registry" />
</p>

Sheater backend は REST API と認証専用 gRPC サービスで構成されています。

| Service | Role | Port |
| ------- | ---- | ---- |
| `api` | REST API / OpenAPI handler / Supabase DB・PostgREST 経由のアプリデータ操作 | `8080` |
| `api-watcher` | API と同一 image を使うログ監視用コンテナ | `10100:8080` |
| `auth` | access token 検証 / Redis token cache / Supabase Auth 連携 | `50051`, `25000` |
| `redis` | auth token cache | `6379` |

### Architecture

```text
Frontend
  |
  | HTTP + Bearer access token
  v
api  -- gRPC -->  auth  -- token cache --> Redis
 |                 |
 |                 | Supabase Auth
 |                 v
 |              Supabase Auth API
 |
 +-- pgx --------> Supabase DB
 |
 +-- PostgREST --> Supabase REST API
```

* `api` は Echo + OpenAPI 生成コードで REST endpoint を提供します。
* `auth` は gRPC で access token を検証し、Redis に token cache を保存します。
* `/health` は Kubernetes の liveness/readiness probe 用に public endpoint として扱います。
* `/metrics` は Prometheus scrape 用です。本番では外部公開せず、クラスタ内部から Prometheus が取得します。

### Local Development

```bash
task dep
task gen
task dup
```

よく使うコマンド:

| Command | Purpose |
| ------- | ------- |
| `task gen` | OpenAPI / Protocol Buffers の生成コード更新 |
| `task dep` | `api` / `auth` の Go modules 整理 |
| `task dup` | Docker Compose を background 起動 |
| `task bup` | Docker Compose を foreground 起動 |
| `task down` | Docker Compose 停止 |
| `task logsi` | `api` logs |
| `task logsh` | `auth` logs |

### CI/CD

`main` への push で GitHub Actions が実行されます。

1. OpenAPI / Protocol Buffers のコード生成
2. `api` / `auth` の test
3. GHCR へ Docker image push
4. `k3s-manifests` の image tag 更新
5. Argo CD が manifest 変更を検知して k3s に反映

## Development Environment

<p>
  <img src="https://skillicons.dev/icons?i=docker" height="50" alt="Docker" />
  <img src="https://skillicons.dev/icons?i=vscode" height="50" alt="VSCode" />
  <img src="https://skillicons.dev/icons?i=git" height="50" alt="Git" />
  <img src="https://skillicons.dev/icons?i=github" height="50" alt="GitHub" />
</p>

* Docker
* Docker Compose
* VSCode
* Git
* GitHub

---

# 📋 Issue Management

Issue は以下のテンプレートを使用します。

| Type    | Purpose        |
| ------- | -------------- |
| Feature | 新機能            |
| Bug     | 不具合報告          |
| Fix     | 既存機能・設定修正      |
| Task    | 調査・環境構築・ドキュメント |

# 📝 Commit Convention

本プロジェクトでは Conventional Commits を採用します。

## Commit Types

| Type        | Purpose          |
| ----------- | ---------------- |
| `feat:`     | 新機能追加            |
| `fix:`      | バグ修正             |
| `docs:`     | README・ドキュメント修正  |
| `refactor:` | リファクタリング（動作変更なし） |
| `test:`     | テスト追加・修正         |
| `chore:`    | ビルド・設定・依存関係更新    |
| `style:`    | コードフォーマット・Lint対応 |
| `perf:`     | パフォーマンス改善        |

## Examples

```text
feat: create room API

feat: implement websocket sync

fix: room join validation error

docs: update README architecture

refactor: simplify room manager

test: add websocket unit tests

chore: add docker compose

style: format backend code

perf: optimize room lookup
```

## Commit Policy

* 1コミット = 1目的
* コミットメッセージは英語でも日本語でもok
* Issue に紐づく場合は Issue 番号を記載

例:

```text
feat: implement room creation (#12)

fix: websocket reconnect bug (#15)
```

---
