package runner

import "fmt"

// ANSI color codes
var colorCodes = map[string]string{
    "black":   "30",
    "red":     "31",
    "green":   "32",
    "yellow":  "33",
    "blue":    "34",
    "magenta": "35",
    "cyan":    "36",
    "white":   "37",
}

// DisplayCommand prints the command in the specified color if enabled
func DisplayCommand(cmd string, color string) {
    if code, ok := colorCodes[color]; ok {
        fmt.Printf("\033[%sm%s\033[0m\n", code, cmd)
    } else {
        fmt.Println(cmd)
    }
}
