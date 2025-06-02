# Toggl Daily Report CLI Tool

Toggl Track APIを利用して日報を自動生成するCLIツール

## インストール

```bash
go install github.com/toggl-daily-report
```

## 設定

### 方法1: 環境変数
```bash
export TOGGL_API_TOKEN="your_toggl_api_token"
```

### 方法2: 設定ファイル
実行ファイルと同じディレクトリに`.toggl-daily-report.json`を作成:
```json
{
  "api_token": "your_toggl_api_token",
  "workspace_id": "your_workspace_id"
}
```

APIトークンはTogglの[プロフィール設定](https://track.toggl.com/profile)から取得できます。

## 使用方法

```bash
# 当日分の日報を生成
toggl-daily-report

# 特定日付の日報を生成
toggl-daily-report --date 2025-06-01

# プロジェクトでフィルタリング
toggl-daily-report --project "ProjectName"

# 複数オプションの組み合わせ
toggl-daily-report --date 2025-06-01 --project "ProjectA"

# ファイルに保存
toggl-daily-report > daily-report.md
```

## ビルド

```bash
go build -o toggl-daily-report
```