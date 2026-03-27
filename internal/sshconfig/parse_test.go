package sshconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseFile_hostBlocksAndComments(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "config")
	content := `
# top comment

Host foo bar
  HostName real.example.com
  User u

Host *
  ForwardAgent yes

Host wild *.ex
  HostName x

Host onlypat *.example.org
`
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	blocks, err := ParseFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 4 {
		t.Fatalf("blocks: got %d want 4", len(blocks))
	}

	// First block: two patterns
	if len(blocks[0].Patterns) != 2 || blocks[0].Patterns[0] != "foo" || blocks[0].Patterns[1] != "bar" {
		t.Fatalf("patterns: %#v", blocks[0].Patterns)
	}
	if !strings.Contains(BlockText(&blocks[0]), "User u") {
		t.Fatalf("block text: %q", BlockText(&blocks[0]))
	}

	cands := CandidatesFromBlocks(blocks)
	// Host * -> no literals -> skip entire block
	// wild *.ex -> literals: wild; *.ex skipped
	// onlypat *.example.org -> literal onlypat; *.example.org skipped
	var names []string
	for _, c := range cands {
		names = append(names, c.Alias)
	}
	want := []string{"foo", "bar", "wild", "onlypat"}
	if len(names) != len(want) {
		t.Fatalf("candidates: got %#v want %#v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("candidates: got %#v want %#v", names, want)
		}
	}

	var fooc *Candidate
	for i := range cands {
		if cands[i].Alias == "foo" {
			fooc = &cands[i]
			break
		}
	}
	if fooc == nil || fooc.HostName != "real.example.com" {
		t.Fatalf("foo HostName: %#v", fooc)
	}
}

func TestResolvePath_SSH_CONFIG(t *testing.T) {
	dir := t.TempDir()
	want := filepath.Join(dir, "myconfig")
	t.Setenv("SSH_CONFIG", want)
	got, err := ResolvePath()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("ResolvePath: got %q want %q", got, want)
	}
}
