package r2libs

import (
	"strings"
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// resetCommandPolicy restores the default (unrestricted) CommandPolicy,
// since it's process-global state shared across tests.
func resetCommandPolicy(t *testing.T) {
	t.Helper()
	SetCommandPolicy(CommandPolicy{})
	t.Cleanup(func() { SetCommandPolicy(CommandPolicy{}) })
}

func runR2OSScript(t *testing.T, code string) (result interface{}, panicked interface{}) {
	t.Helper()
	env := r2core.NewEnvironment()
	env.Set("true", true)
	env.Set("false", false)
	env.Set("nil", nil)
	RegisterStd(env)
	RegisterOS(env)

	defer func() {
		panicked = recover()
	}()
	parser := r2core.NewParser(code)
	program := parser.ParseProgram()
	result = program.Eval(env)
	return
}

func TestCommandPolicy_DefaultIsUnrestricted(t *testing.T) {
	resetCommandPolicy(t)
	_, panicked := runR2OSScript(t, `
let cmd = os.Command("echo hello")
let result = cmd.run()
`)
	if panicked != nil {
		t.Fatalf("expected default policy to allow any command, got panic: %v", panicked)
	}
}

func TestCommandPolicy_AllowedCommandsBlocksDisallowed(t *testing.T) {
	resetCommandPolicy(t)
	SetCommandPolicy(CommandPolicy{AllowedCommands: map[string]bool{"echo": true}})

	_, panicked := runR2OSScript(t, `
let cmd = os.Command("rm -rf /tmp/should-not-run")
cmd.run()
`)
	if panicked == nil {
		t.Fatal("expected a command not on the allowlist to panic, got none")
	}
	msg, ok := panicked.(string)
	if !ok || !strings.Contains(msg, "not in the allowed command list") {
		t.Fatalf("expected a clear allowlist panic message, got: %v", panicked)
	}
}

func TestCommandPolicy_AllowedCommandsPermitsListed(t *testing.T) {
	resetCommandPolicy(t)
	SetCommandPolicy(CommandPolicy{AllowedCommands: map[string]bool{"echo": true}})

	_, panicked := runR2OSScript(t, `
let cmd = os.Command("echo hello")
cmd.run()
`)
	if panicked != nil {
		t.Fatalf("expected an allowlisted command to run, got panic: %v", panicked)
	}
}

func TestCommandPolicy_DisableShellBlocksExecCmd(t *testing.T) {
	resetCommandPolicy(t)
	SetCommandPolicy(CommandPolicy{DisableShell: true})

	_, panicked := runR2OSScript(t, `os.execCmd("echo hello")`)
	if panicked == nil {
		t.Fatal("expected execCmd to panic when shell execution is disabled")
	}
	msg, ok := panicked.(string)
	if !ok || !strings.Contains(msg, "shell command execution is disabled") {
		t.Fatalf("expected a clear shell-disabled panic message, got: %v", panicked)
	}
}

func TestCommandPolicy_DisableShellBlocksRunProcess(t *testing.T) {
	resetCommandPolicy(t)
	SetCommandPolicy(CommandPolicy{DisableShell: true})

	_, panicked := runR2OSScript(t, `os.runProcess("echo hello")`)
	if panicked == nil {
		t.Fatal("expected runProcess to panic when shell execution is disabled")
	}
}

func TestCommandPolicy_DisableShellDoesNotAffectNonShellCommand(t *testing.T) {
	resetCommandPolicy(t)
	SetCommandPolicy(CommandPolicy{DisableShell: true})

	// os.Command(...).run() never goes through a shell, so DisableShell
	// alone must not block it.
	_, panicked := runR2OSScript(t, `
let cmd = os.Command("echo hello")
cmd.run()
`)
	if panicked != nil {
		t.Fatalf("expected DisableShell to not affect the non-shell os.Command path, got panic: %v", panicked)
	}
}
