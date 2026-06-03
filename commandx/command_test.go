package commandx

import (
	"strings"
	"testing"
	"time"
)

func TestRunCommand_Argv(t *testing.T) {
	res := RunCommand([]string{"go", "version"})
	if !res.Success || res.Status != 0 {
		t.Fatalf("expected success, got success=%v status=%d err=%v stderr=%q",
			res.Success, res.Status, res.Err, res.Stderr.String())
	}
	if !strings.Contains(res.Stdout.String(), "go version") {
		t.Fatalf("unexpected stdout: %q", res.Stdout.String())
	}
}

func TestRunBashCommand_PipeAndRedirect(t *testing.T) {
	res := RunBashCommand("echo hello | tr 'a-z' 'A-Z'")
	if !res.Success {
		t.Fatalf("expected success, got err=%v stderr=%q", res.Err, res.Stderr.String())
	}
	if strings.TrimSpace(res.Stdout.String()) != "HELLO" {
		t.Fatalf("unexpected stdout: %q", res.Stdout.String())
	}
}

func TestRunCommand_NonZeroExit(t *testing.T) {
	res := RunBashCommand("exit 7")
	if res.Success {
		t.Fatalf("expected failure")
	}
	if res.Status != 7 {
		t.Fatalf("expected status 7, got %d", res.Status)
	}
}

func TestRunCommand_WithDir(t *testing.T) {
	res := RunCommand([]string{"pwd"}, WithDir("/tmp"))
	if !res.Success {
		t.Fatalf("expected success, got err=%v", res.Err)
	}
	got := strings.TrimSpace(res.Stdout.String())
	if got != "/tmp" && got != "/private/tmp" {
		t.Fatalf("expected /tmp, got %q", got)
	}
}

func TestRunCommand_RedirectStderr(t *testing.T) {
	res := RunBashCommand("echo err >&2; echo out", WithRedirectStderr())
	if !res.Success {
		t.Fatalf("expected success, got err=%v", res.Err)
	}
	out := res.Stdout.String()
	if !strings.Contains(out, "err") || !strings.Contains(out, "out") {
		t.Fatalf("expected both err and out in stdout, got stdout=%q stderr=%q",
			out, res.Stderr.String())
	}
	if res.Stderr.String() != "" {
		t.Fatalf("expected empty stderr, got %q", res.Stderr.String())
	}
}

func TestRunCommand_Timeout(t *testing.T) {
	start := time.Now()
	res := RunBashCommand("sleep 5", WithTimeout(300*time.Millisecond))
	elapsed := time.Since(start)
	if res.Success {
		t.Fatalf("expected failure on timeout")
	}
	if elapsed > 2*time.Second {
		t.Fatalf("timeout did not interrupt: elapsed=%v", elapsed)
	}
}

func TestRunCommand_ArgvWithSpecialChars(t *testing.T) {
	res := RunCommand([]string{"echo", "hello world", "a'b", "$VAR"})
	if !res.Success {
		t.Fatalf("expected success, got err=%v stderr=%q", res.Err, res.Stderr.String())
	}
	got := strings.TrimSpace(res.Stdout.String())
	want := "hello world a'b $VAR"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestRunCommand_Empty(t *testing.T) {
	res := RunCommand(nil)
	if res.Success {
		t.Fatalf("expected failure on empty command")
	}
	if res.Err == nil {
		t.Fatalf("expected error")
	}
}
