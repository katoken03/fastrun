---
date: 2025-08-28
tags: [package-manager, npm, pnpm, bun, lockfile-detection]
---

# Fastrun パッケージマネージャー自動判定機能追加

## 要件

### 背景
Fastrунで`f`コマンドを実行してpackage.jsonのスクリプトを選択・実行する際、常に`npm run`コマンドが使用されていました。しかし、プロジェクトによってはpnpmやbunなどの異なるパッケージマネージャーを使用している場合があり、適切なパッケージマネージャーで実行したいという要望がありました。

### 目的
- ロックファイルの存在に基づいてパッケージマネージャーを自動判定
- `pnpm-lock.yaml` → pnpm
- `bun.lockb` → bun  
- どちらも存在しない → npm（デフォルト）

## 実施内容

### 1. Runner構造体の拡張
```go
// plugins/npm/plugin.go
type Runner struct{
    packageManager string
}
```

### 2. パッケージマネージャー検出ロジック追加
```go
// detectPackageManager detects the appropriate package manager based on lock files
func (r *Runner) detectPackageManager(path string) string {
    // Check for pnpm-lock.yaml
    if _, err := os.Stat(filepath.Join(path, "pnpm-lock.yaml")); err == nil {
        return "pnpm"
    }
    
    // Check for bun.lockb
    if _, err := os.Stat(filepath.Join(path, "bun.lockb")); err == nil {
        return "bun"
    }
    
    // Default to npm
    return "npm"
}
```

### 3. ParseCommands()メソッドの修正
- ロックファイル検出を追加
- ExecuteCommandの生成を動的に変更

```go
// Detect the appropriate package manager
r.packageManager = r.detectPackageManager(path)

var commands []runner.Command
scripts.ForEach(func(key, value gjson.Result) bool {
    cmd := runner.Command{
        Name:           key.String(),
        Description:    value.String(),
        ExecuteCommand: fmt.Sprintf("%s run %s", r.packageManager, key.String()),
    }
    commands = append(commands, cmd)
    return true
})
```

### 4. RunCommand()メソッドの修正
- ハードコードされた"npm"を動的なpackageManagerに変更

```go
func (r *Runner) RunCommand(cmd runner.Command) error {
    fullCmd := fmt.Sprintf("%s run %s", r.packageManager, cmd.Name)
    runner.DisplayCommand(fullCmd, "cyan")

    c := exec.Command(r.packageManager, "run", cmd.Name)
    c.Stdout = os.Stdout
    c.Stderr = os.Stderr
    c.Stdin = os.Stdin
    return c.Run()
}
```

## 主要な変更ファイル

- **plugins/npm/plugin.go**: パッケージマネージャー自動判定機能の実装

## 得られた知見

### 設計面
1. **既存のプラグインアーキテクチャを活用**: 既存のCommandRunnerインターフェースを変更することなく、npm runnerの内部ロジックのみを拡張することで実装
2. **状態管理**: Runner構造体にpackageManagerフィールドを追加し、ParseCommands()で判定結果を保存、RunCommand()で使用する設計

### 技術面
1. **ロックファイル検出**: `os.Stat()`を使用してファイルの存在確認を行う簡潔な実装
2. **Go言語のファイル操作**: `filepath.Join()`を使用したクロスプラットフォーム対応のパス結合
3. **既存コードとの整合性**: DisplayCommandやexec.Commandの既存の使用パターンを維持

### UX改善
- ユーザーは何も設定することなく、プロジェクトに適したパッケージマネージャーでスクリプトが実行される
- `--text-only`モードでも正しいコマンド（`pnpm run dev`など）が出力される

## 技術詳細

### 判定ロジックの優先順位
1. `pnpm-lock.yaml`の存在確認
2. `bun.lockb`の存在確認  
3. デフォルトは`npm`

### 対応パッケージマネージャー
- **npm**: Node.js標準のパッケージマネージャー
- **pnpm**: 高速で効率的なパッケージマネージャー
- **bun**: JavaScriptランタイム兼パッケージマネージャー

この実装により、Fastrунがより多くのJavaScriptプロジェクトで適切に動作するようになりました。