# fastrun - Claude 開発ガイドライン

## リリース方法

### 手順
1. 変更をコミットする
2. バージョンタグを付ける（例: `v1.0.6`）
3. タグをリモートにプッシュする

```bash
git tag v1.0.x
git push origin main
git push origin v1.0.x
```

### リリースの仕組み
- `v*` 形式のタグをプッシュすると GitHub Actions の Release ワークフローが自動起動する
- ワークフローは [GoReleaser](https://goreleaser.com/) を使ってバイナリをビルドし、GitHub Release を自動作成する
- Homebrew tap へのパブリッシュも自動で行われる（`HOMEBREW_GITHUB_API_TOKEN` シークレット使用）

### 現在の最新バージョン確認
```bash
git tag --sort=-v:refname | head -1
```

---

## git worktree ルール

Git Worktree は以下のディレクトリをベースディレクトリとする。


```
../fastrun.worktrees/
```

こちらのディレクトリがまだ存在しない場合は作成する。


---

## 設定ファイル

- **読み込み順**: `~/.config/fastrun/config.json` → カレントの `.fastrun/config.json`（後者が上書き）。ファイルが無ければスキップ。
- **形式**: JSON。詳細は README の Configuration を参照。

| キー | 概要 |
|------|------|
| `fzf_position` | `"top"` / `"bottom"`（fzf に `--reverse` を付けるか） |
| `command_color` | コマンド名列の色名（`cyan` 等、`internal/ui/fzf.go` の `colorCodes`） |
| `use_nr` | `false` で `nr` を使わない。省略または `true` は従来どおり PATH に `nr` があれば利用 |

---

## コーディング方針

### 単一責任の原則（SRP）を守ること

各ファイル・関数・構造体は **1つの責務** のみを持つこと。

- **プラグイン**（`plugins/npm`, `plugins/make`）: コマンドのパースと実行のみ
- **UI**（`internal/ui`）: ユーザーへの表示と選択のみ
- **runner**（`internal/runner`）: 共通インターフェースの定義のみ

関数が複数の役割を担い始めたら、分割を検討すること。
