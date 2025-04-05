# go-sort-tabelog

Tabelog（食べログ）上の店舗をエリアとキーワードで検索し、評価順に上位店舗を返す Go 製の LINE Bot です。
CLI ツールとしても使用可能です。

## 機能概要

- LINE Bot として動作し、ユーザーからのメッセージに応じて食べログ検索を実行
- CLI ツールで指定したエリア・キーワードに応じて上位 5 件の飲食店情報を取得
- 店舗情報は Agouti を利用してブラウザ操作によりスクレイピング
- LINE 通知やアラート送信機能あり（alert-linebot 連携）

---

## 使用技術

- Go
- [Cobra](https://github.com/spf13/cobra): CLI フレームワーク
- [Agouti](https://github.com/sclevine/agouti): Web スクレイピング用のヘッドレスブラウザ制御
- [LINE Messaging API](https://developers.line.biz/ja/services/messaging-api/)
- [logrus](https://github.com/sirupsen/logrus): ログ出力
- Gorilla Mux: HTTP ルーティング

## 構成

クリーンアーキテクチャをベースに実装しています。

```
.
├── cli/ # CLIツール実装
│ └── cmd/ # CLIコマンド (root, search)
├── cmd/ # Webサーバ起動エントリポイント
├── entity/ # ドメインモデル（Shopなど）
├── interface/ # 外部インターフェース (LINE, AlertBot, WebDriver)
│ ├── controller/ # HTTPハンドラ
│ └── gateway/ # 外部サービス連携
├── server/ # サーバ起動処理
├── usecase/ # アプリケーションロジック
├── Makefile # ビルド定義
├── go.mod # Go module定義
└── README.md # このファイル
```

## 環境変数

| 変数名                        | 内容                            |
| ----------------------------- | ------------------------------- |
| `SSH_CERT_PATH`               | TLS 証明書ファイルのパス        |
| `SSH_KEY_PATH`                | TLS 鍵ファイルのパス            |
| `TABELOG_SORT_CHANNEL_SECRET` | LINE チャンネルシークレット     |
| `TABELOG_SORT_CHANNEL_TOKEN`  | LINE チャンネルアクセストークン |

## ビルド方法

```bash
# Webサーババイナリのビルド
make build

# CLIバイナリのビルド
make build-cli
```

## 実行方法

### サーバ起動

```bash
./bin/tabelogbot --addr 0.0.0.0 --port 443
```

### CLI 利用例

```bash
./bin/tabelogbot-cli search --area 渋谷 --keyword ラーメン
```

## 作者

jgvkmea
https://github.com/jgvkmea
