<div align="center">
	<br>
	<br>
	<br>
	<img alt="fastrun" width="300" src="./f-logo.svg">
	<br>
	<br>
	<h3>The fastest way to run NPM scripts and Makefile targets</h3>
	<br>
</div>

## Install

```bash
brew install katoken03/fastrun/fastrun
```

## What is FastRun?

Instead of typing...

```bash
npm run dev
```

Just type:

```bash
f [Enter]
```

<div align="center">
	<br>
	<img alt="overview" width="400" src="./f-overview.gif">
	<br>
	<br>
</div>

FastRun instantly shows all available commands in your project — npm scripts and Makefile targets — and lets you select one with fuzzy search or arrow keys. No need to remember command names.

## Features

- **Instant fuzzy search** — powered by [fzf](https://github.com/junegunn/fzf)
- **npm & Makefile support** — automatically detects scripts in your project
- **Smart package manager detection** — works with npm, pnpm, bun, and [ni](https://github.com/antfu/ni)
- **Shell history integration** — selected commands are added to your shell history
- **SSH host picker** — `f ssh` (alias `f s`) lists `Host` aliases from `~/.ssh/config` (or `SSH_CONFIG`) in fzf with a config stanza preview, then connects


## Configuration

FastRun reads optional **JSON** config (UTF-8). Use a single root object and set only the keys you need.

**Where it is loaded from (later wins):**

| Order | Path | Role |
| --- | --- | --- |
| 1 | `~/.config/fastrun/config.json` | User-wide defaults |
| 2 | `.fastrun/config.json` (under the directory you run `f` from) | Per-project overrides |

If a file is missing, it is skipped. Keys that are **not** present in a file do not erase values already merged from an earlier file.

**Parameters:**

| Key | Type | Values / notes |
| --- | --- | --- |
| `fzf_position` | string | `"top"` (default) or `"bottom"`. FastRun passes `--reverse` to `fzf` when this is `"top"`. Other values are ignored. |
| `command_color` | string | ANSI color name for the **command name** column in the list: `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`. Other strings fall back to uncolored text. Default is `cyan`. |
| `use_nr` | boolean | `true` — prefer the `nr` command when it is on `PATH` ([ni](https://github.com/antfu/ni)). `false` — never use `nr`; npm scripts use `npm` / `pnpm` / `bun` according to lockfiles. **Omit the key** to keep the classic behavior (try `nr` if on `PATH`). |

**Tip:** Set `use_nr` to `false` if a version manager leaves an `nr` shim on `PATH` but the active Node install does not actually ship `nr`.

**Example — `~/.config/fastrun/config.json`:**

```json
{
  "fzf_position": "bottom",
  "command_color": "green",
  "use_nr": false
}
```

**Example — project-only `.fastrun/config.json`:**

```json
{
  "use_nr": false
}
```

## Usage 

## npm / make

Run `f` in any project directory:

```bash
f
```

FastRun will list all available commands. Use fuzzy search to filter, arrow keys to navigate, and Enter to execute.

### ssh

From any directory:

```bash
f ssh
```

or

```bash
f s
```


Uses the same `fzf_position` as the main launcher. Only literal `Host` patterns are listed (entries that are only wildcards, e.g. `Host *`, are skipped). `Include` is not expanded. Without the shell wrapper, `fastrun ssh` (or `fastrun s`) runs `ssh` directly after you choose a host.

## Shell History Integration

To have selected commands saved to your shell history (accessible with the ↑ key), add the following snippet to your shell config.

### Bash (`.bash_profile` or `.bashrc`)

```bash
f() {
    if [[ "$1" == --* ]]; then
        command f "$@"
        return
    fi
    local cmd=$(command f -t "$@")
    if [ $? -eq 0 ] && [ -n "$cmd" ]; then
        echo -e "\033[36m$cmd\033[0m"
        history -s "$cmd"
        eval "$cmd"
    fi
}
```

### Zsh (`.zshrc`)

```zsh
f() {
    if [[ "$1" == --* ]]; then
        command f "$@"
        return
    fi
    local cmd=$(command f -t "$@")
    if [ $? -eq 0 ] && [ -n "$cmd" ]; then
        echo -e "\033[36m$cmd\033[0m"
        print -s "$cmd"
        eval "$cmd"
    fi
}
```

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

## License

[MIT](LICENSE)
