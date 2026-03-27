package sshpick

import (
	"strings"
	"testing"
)

func TestColorizeSSHComments(t *testing.T) {
	t.Parallel()
	in := "Host x\n  # note\nUser u"
	out := colorizeSSHComments(in)
	if !strings.Contains(out, sshPreviewComment) {
		t.Fatalf("expected comment color: %q", out)
	}
	if !strings.HasPrefix(strings.TrimSpace(strings.Split(out, "\n")[1]), sshPreviewComment) {
		t.Fatalf("second line should be comment: %q", out)
	}
	if strings.Count(out, "User u") != 1 {
		t.Fatalf("User line preserved once: %q", out)
	}
}

func TestColorizeSSHComments_empty(t *testing.T) {
	t.Parallel()
	if colorizeSSHComments("") != "" {
		t.Fatal("empty in should stay empty")
	}
}
