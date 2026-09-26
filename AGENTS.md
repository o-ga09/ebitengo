# Wave Duel — Agent Instructions

2人対戦アクションパズルゲーム。波の重ね合わせ原理をコアメカニクスとした Ebitengine (Go) 製ゲーム。

詳細仕様: [docs/requirements.md](docs/requirements.md)
オンライン対戦設計: [docs/online-design.md](docs/online-design.md)

---

## クイックスタート

```bash
go mod tidy   # 依存関係インストール
make serve    # ブラウザ版（WASM）をビルドして http://localhost:8080 で配信
make run      # デスクトップ版を起動
make test     # テスト
```

---

## 開発ルール

- **main への直接 push を許可する。** このリポジトリではブランチ・PR を作らず、main にコミットして push してよい（`o-ga09/infra` は push が即本番反映のため対象外。push 前に必ず確認する）
- **ロードマップは README.md で管理する。** 項目を完了したら、同じ変更（コミット）の中で README.md の該当項目を `⬜` → `✅` に更新する。新しい作業が発生したら該当フェーズに項目を追加する
- **ブラウザ（WASM）を主要ターゲットとする。** 変更後は `GOOS=js GOARCH=wasm go build ./cmd` が通ることを確認する。WASM で使えない API（`os/exec`、ファイル I/O、生の TCP/UDP ソケットなど）をクライアントのコードで使わない

---

## アーキテクチャ

### ファイル責務

| ファイル | 責務 |
|---|---|
| `cmd/main.go` | エントリーポイント。`ebiten.RunGame()` 呼び出しのみ |
| `internal/game.go` | ゲーム状態管理・メインループ (`Update` / `Draw`) |
| `internal/wave.go` | 波の定義・計算 (`ValueAt`, `compositeAmplitude`) |
| `internal/player.go` | プレイヤー状態・操作入力処理 |
| `internal/field.go` | フィールド描画・環境波ギミック |
| `internal/battle.go` | 戦闘ロジック（HP・ダメージ・共鳴ゲージ） |
| `internal/ui.go` | HP バー・ゲージ・波形プレビュー描画 |
| `internal/input_pc.go` | PC キーボード入力処理 |
| `internal/input_mobile.go` | タッチ・フリック入力処理 |

## コーディング規約

- **ロジックは `internal/` 配下のみ**。`cmd/main.go` にゲームロジックを書かない
- `internal/` 直下のパッケージ名はすべて `package game`（`input_pc.go` / `input_mobile.go` はビルドタグで分離可）
- 外部公開 API は最小限にとどめ、`internal/` のパッケージ間インターフェースのみ公開する
- テストは `_test.go` に記述し、`wave.go` の計算ロジックから優先的にカバーする
