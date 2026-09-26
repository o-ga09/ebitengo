# Wave Duel (仮)

2人のプレイヤーが**波を生成・発射**し、フィールド上で発生する**合成波の干渉現象**を利用して相手にダメージを与える対戦ゲーム。
ブラウザ（WebAssembly）とデスクトップで動作し、将来的に Agones + ゲーム API サーバによるオンライン対戦に対応する。

- ゲーム仕様: [docs/requirements.md](docs/requirements.md)
- オンライン対戦設計: [docs/online-design.md](docs/online-design.md)

## セットアップ

```bash
go mod tidy
```

### ブラウザで遊ぶ

```bash
make serve   # WASM をビルドして http://localhost:8080 で配信
```

ビルドだけなら `make wasm`（`web/` に `game.wasm` と `wasm_exec.js` が出力される）。

### デスクトップで遊ぶ

```bash
make run     # = go run ./cmd
```

## ロードマップ

完了した項目には ✅ を付ける。フェーズの詳細は [docs/online-design.md §11](docs/online-design.md#11-開発フェーズ) を参照。

### P0: 準備

- ✅ ローカル 2P ターン制プロトタイプ（波の描画・発射・ダメージ・勝敗）
- ✅ オンライン対戦化の設計ドキュメント
- ✅ ブラウザ（WASM）で起動できる

### P1: リアルタイム化（ローカル）

- ⬜ ebiten に依存しない `internal/sim` の切り出し
- ⬜ リアルタイムルール v0（上下移動・波パケット・被弾判定・エネルギー・制限時間）
- ⬜ `sim` のユニットテスト
- ⬜ ブラウザで 2P リアルタイム対戦が遊べる

### P2: ネットワーク対戦（Agones なし）

- ⬜ Protocol Buffers でメッセージを定義
- ⬜ `cmd/gameserver`（`-local` モード）
- ⬜ クライアントのオンラインモード（予測・補正・補間・時刻同期）
- ⬜ 遅延・ロスを注入した検証

### P3: Agones + ゲートウェイ

- ⬜ おうちk3s に Agones を導入
- ⬜ Fleet / FleetAutoscaler
- ⬜ GameServer への Agones SDK 組み込み
- ⬜ `cmd/gateway`（WebSocket 中継）

### P4: ゲーム API

- ⬜ ゲストログイン
- ⬜ 簡易マッチング（Redis キュー）と GameServer の割り当て
- ⬜ 試合結果の保存（MySQL）

### P5: 公開・運用

- ⬜ CI/CD（GAR へのイメージ push、infra マニフェストの自動更新）
- ⬜ ブラウザ版クライアントの配信（`wave-duel.o-ga09.com`）
- ⬜ Cloudflare Tunnel 経由の外部公開
- ⬜ 監視ダッシュボード

### P6: ゲーム要素の拡充

- ⬜ 共鳴ゲージ・必殺波
- ⬜ 環境波
- ⬜ モバイルブラウザのタッチ操作
- ⬜ エフェクト・サウンド
