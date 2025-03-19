---
date: 2025-03-19
tags: [fastrun, shell-history, go, cobra]
author: Claude
---

# Fastrun シェルヒストリー問題の解決: `-t`オプション追加

## 要件

Fastrunを使って実行したコマンドをbashのヒストリーに残すための機能改善。

**課題**:
- Fastrunは実行しているbashの子プロセスで起動するため、親のbashプロセスのヒストリーに記録を残せない
- 現在は`f generate-shell-function`で作成したシェル関数を使用して解決しているが、動作が不安定

**解決策**:
- `-t`オプション（`--text-only`）を追加し、選択したコマンドのテキストだけを返すモードを実装
- シェル関数を改良し、返されたテキストをヒストリーに追加してから実行する

## 実施内容

1. **コードの調査と理解**
   - fastrun の基本構造と動作の理解
   - 既存のシェル関数生成機能の分析

2. **`-t`オプションの実装**
   - `cmd/root.go`にテキストモード用のフラグを追加
   - コマンド実行ロジックに分岐を追加し、テキストモードでは選択したコマンドのテキストのみを出力

3. **シェル関数の改善**
   - 新しいシェル関数テンプレートの追加（bash と zsh 用）
   - `generate-shell-function`コマンドの更新

4. **コミットとタグ付け**
   - 変更内容をコミット
   - バージョンタグを追加

## コード変更

### 1. `-t`フラグの追加 (`cmd/root.go`)

```go
var textOnly bool

func init() {
    // テキストモードのフラグを追加
    rootCmd.Flags().BoolVarP(&textOnly, "text-only", "t", false, "Print selected command text only without execution")
}
```

### 2. テキストモード処理の追加 (`cmd/root.go`)

```go
// テキストのみモードの場合は、コマンドのテキストを出力して終了
if textOnly {
    // 実行コマンドだけをテキストとして出力
    for _, r := range runners {
        cmds, err := r.ParseCommands(cwd)
        if err != nil {
            continue
        }
        for _, cmd := range cmds {
            if cmd.Name == selectedCmd.Name {
                fmt.Println(cmd.ExecuteCommand)
                return nil
            }
        }
    }
    return fmt.Errorf("command not found: %s", selectedCmd.Name)
}
```

### 3. 新しいシェル関数 (`cmd/generate.go`)

**Bash用**:
```bash
f() {
    local cmd=$(command f -t "$@")
    if [ $? -eq 0 ] && [ -n "$cmd" ]; then
        # コマンドをヒストリーに追加
        history -s "$cmd"
        # コマンドを実行
        eval "$cmd"
    fi
}
```

**Zsh用**:
```bash
f() {
    local cmd=$(command f -t "$@")
    if [ $? -eq 0 ] && [ -n "$cmd" ]; then
        # コマンドをヒストリーに追加
        print -s "$cmd"
        # コマンドを実行
        eval "$cmd"
    fi
}
```

## 得られた知見

1. **子プロセスとシェルヒストリーの関係**
   - シェルのヒストリーは基本的に親プロセスでのみ管理されており、子プロセスでの操作は記録されない
   - コマンドをヒストリーに追加するには、親プロセスのコンテキストで`history -s`や`print -s`を実行する必要がある

2. **テキストモードと実行モードの分離**
   - CLIツールにテキストモードを追加することで、より柔軟な使用方法が可能になる
   - 分離によって、複雑な問題を単純なテキスト処理に変換できる

3. **シェル関数とワンライナー**
   - 複数のコマンドを組み合わせた操作をシェル関数にまとめることで、使い勝手を向上できる
   - `eval`を使用する際は、セキュリティに注意が必要

## 技術詳細

- **コミットID**: `0bc0d8f11d2e56c36c530b386bebc55c6f8fe3c9`
- **変更ファイル**: 
  - `cmd/root.go` - テキストモードフラグとロジックの追加
  - `cmd/generate.go` - 新しいシェル関数の実装

## 使用方法

1. **ビルド**:
   ```bash
   go build
   ```

2. **シェル関数の生成と設定**:
   ```bash
   f generate-shell-function
   ```

3. **使用方法**:
   ```bash
   f
   # (コマンドを選択)
   # 選択したコマンドがヒストリーに追加されて実行される
   ```

## 今後の課題

- エラーハンドリングの改善
- タブ補完機能の追加
- より多くのシェル（fish等）のサポート
