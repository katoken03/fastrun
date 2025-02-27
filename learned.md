# Fastrun プロジェクトの知見

## シェルヒストリーへのコマンド追加機能

### 問題と解決策

#### 問題
- fastrunで選択したコマンドがシェルのヒストリーに追加されない
- ユーザーは選択したコマンドを再度利用する際に、手動で入力する必要があった

#### 解決策
1. **特殊なマーカーの追加**
   - コマンド表示時に特殊なマーカー（`FASTRUN_CMD:`）を追加
   - これにより、シェル関数でコマンドを簡単に識別できるようになった

2. **シェル関数の実装**
   - `f` コマンドをラップするシェル関数を作成
   - 特殊なマーカーを使用してコマンドを抽出し、ヒストリーに追加
   - マーカー行を除いて出力を表示

3. **シェル関数生成コマンドの追加**
   - `generate-shell-function` コマンドを実装
   - ユーザーが簡単にシェル関数を設定できるよう、設定コマンドも表示

### 技術的な詳細

#### 特殊なマーカーの実装
```go
// DisplayCommand prints the command in the specified color if enabled
func DisplayCommand(cmd string, color string) {
    if code, ok := colorCodes[color]; ok {
        // 特殊なマーカーを追加してコマンドを表示
        fmt.Printf("FASTRUN_CMD:%s\n", cmd)
        fmt.Printf("\033[%sm%s\033[0m\n", code, cmd)
    } else {
        fmt.Printf("FASTRUN_CMD:%s\n", cmd)
        fmt.Println(cmd)
    }
}
```

#### シェル関数の実装（Bash用）
```bash
# ========== fastrun wrapper function for bash
# ========== Add this to your ~/.bash_profile
f() {
    # Run the f command and capture the output
    local cmd_output=$(command f "$@")
    local exit_code=$?
    
    # 特殊なマーカー（FASTRUN_CMD:）を使用してコマンドを抽出
    local cmd=$(echo "$cmd_output" | grep "FASTRUN_CMD:" | sed 's/FASTRUN_CMD://')
    
    # Add command to history if not empty
    if [ -n "$cmd" ]; then
        history -s "$cmd"
    fi
    
    # 特殊なマーカー行を除いて出力を表示
    echo "$cmd_output" | grep -v "FASTRUN_CMD:"
    
    return $exit_code
}
```

### 学んだ教訓

1. **出力の解析に関する問題**
   - 当初、ANSIカラーコードを使用してコマンドを識別しようとしたが、正規表現が複雑になり信頼性が低かった
   - 特殊なマーカーを追加することで、より確実にコマンドを識別できるようになった

2. **シェル関数の再帰呼び出し防止**
   - シェル関数と同名のコマンドを呼び出す際は、`command` プレフィックスを使用して再帰呼び出しを防止する必要がある

3. **シェル関数の追加方法**
   - ヒアドキュメント（`cat << 'EOF'`）を使用することで、シングルクォートやバックスラッシュなどの特殊文字を含むシェル関数を安全に追記できる

4. **ユーザビリティの向上**
   - コマンド生成機能により、ユーザーは設定手順をコピー＆ペーストするだけで済むようになった
   - 目立つコメントを追加することで、シェル設定ファイル内でのシェル関数の識別が容易になった

### 使用方法

1. シェル関数を生成する：
   ```
   f generate-shell-function --shell=bash
   ```

2. 表示されたコマンドをコピー＆ペーストして、シェル関数を `.bash_profile` に追加する

3. シェル設定を再読み込みする：
   ```
   source ~/.bash_profile
   ```

4. `f` コマンドを使用すると、選択したコマンドがヒストリーに追加される
