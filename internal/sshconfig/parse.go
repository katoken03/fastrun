package sshconfig

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Block is one Host stanza (from a "Host" line through the line before the next "Host").
type Block struct {
	Patterns []string
	Lines    []string // raw lines belonging to this stanza (including the Host line)
}

// Candidate is a selectable SSH host alias with metadata for display and preview.
type Candidate struct {
	Alias    string
	HostName string // real HostName from stanza if set, else empty
	Block    *Block
}

// ParseFile reads an OpenSSH config file and returns Host blocks in order.
func ParseFile(path string) ([]Block, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open ssh config: %w", err)
	}
	defer f.Close()

	var blocks []Block
	var cur *Block

	sc := bufio.NewScanner(f)
	// Support long lines
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, 1024*1024)

	for sc.Scan() {
		line := sc.Text()
		t := strings.TrimSpace(line)
		if t == "" || strings.HasPrefix(t, "#") {
			if cur != nil {
				cur.Lines = append(cur.Lines, line)
			}
			continue
		}

		fields := strings.Fields(t)
		if len(fields) >= 2 && strings.EqualFold(fields[0], "Host") {
			b := Block{
				Patterns: append([]string(nil), fields[1:]...),
				Lines:    []string{line},
			}
			blocks = append(blocks, b)
			cur = &blocks[len(blocks)-1]
			continue
		}

		if cur != nil {
			cur.Lines = append(cur.Lines, line)
		}
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("read ssh config: %w", err)
	}
	return blocks, nil
}

func isWildcardPattern(p string) bool {
	return strings.ContainsAny(p, "*?")
}

func hostNameFromBlock(b *Block) string {
	for _, line := range b.Lines {
		t := strings.TrimSpace(line)
		if t == "" || strings.HasPrefix(t, "#") {
			continue
		}
		fields := strings.Fields(t)
		if len(fields) >= 2 && strings.EqualFold(fields[0], "HostName") {
			return fields[1]
		}
	}
	return ""
}

// CandidatesFromBlocks expands blocks into selectable aliases (literals only; wildcard patterns omitted).
func CandidatesFromBlocks(blocks []Block) []Candidate {
	var out []Candidate
	for i := range blocks {
		b := &blocks[i]
		hasLiteral := false
		for _, p := range b.Patterns {
			if !isWildcardPattern(p) {
				hasLiteral = true
				break
			}
		}
		if !hasLiteral {
			continue
		}

		hn := hostNameFromBlock(b)
		for _, p := range b.Patterns {
			if isWildcardPattern(p) {
				continue
			}
			out = append(out, Candidate{
				Alias:    p,
				HostName: hn,
				Block:    b,
			})
		}
	}
	return out
}

// ResolvePath returns the SSH config path: SSH_CONFIG if set, else ~/.ssh/config.
func ResolvePath() (string, error) {
	if p := strings.TrimSpace(os.Getenv("SSH_CONFIG")); p != "" {
		return p, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("user home: %w", err)
	}
	return filepath.Join(home, ".ssh", "config"), nil
}

// BlockText renders a stanza for preview (newline-joined Lines).
func BlockText(b *Block) string {
	if b == nil {
		return ""
	}
	return strings.Join(b.Lines, "\n")
}
