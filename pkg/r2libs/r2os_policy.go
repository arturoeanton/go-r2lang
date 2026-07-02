package r2libs

import (
	"fmt"
	"strings"
	"sync"
)

// CommandPolicy restricts what os.Command/execCmd/runProcess/execWithTimeout/
// execWithEnv are allowed to run. It exists for embedders that run R2Lang
// scripts they don't fully trust: by default (the zero value) every command
// execution builtin behaves exactly as before — unrestricted shell access —
// which matches this project's existing, deliberate design (os.Command is
// shell-injection-shaped on purpose, see CLAUDE.md/session notes). Call
// SetCommandPolicy from the host Go program, before running any untrusted
// script, to opt into restrictions.
type CommandPolicy struct {
	// AllowedCommands, if non-empty, is an allowlist of executable names
	// (the first whitespace-separated token, e.g. "echo", "git") permitted
	// for the NON-shell execution path (os.Command(...).run(), which
	// splits its argument with strings.Fields and calls exec.Command
	// directly — no shell metacharacters are interpreted there). An empty
	// map means unrestricted (default).
	AllowedCommands map[string]bool

	// DisableShell, if true, blocks every builtin that runs its argument
	// through a shell (execCmd, runProcess, execWithTimeout, execWithEnv —
	// all of them do `sh -c cmdLine`). A shell string can chain arbitrary
	// commands via `;`, `&&`, `|`, backticks, `$(...)`, so unlike
	// AllowedCommands there is no safe partial/allowlist mode for these:
	// either the embedder trusts the script with a full shell, or it
	// doesn't. Embedders that need specific host-side operations from a
	// restricted script should expose them via native.callFunc
	// (RegisterNativeFunc) instead, not the shell builtins.
	DisableShell bool
}

var (
	commandPolicyMu sync.RWMutex
	commandPolicy   CommandPolicy // zero value: fully unrestricted (backward compatible default)
)

// SetCommandPolicy installs the active CommandPolicy for every Environment
// registered afterward via RegisterOS in this process. Intended to be
// called once, early, by the host Go program — before running any script
// that isn't fully trusted. Passing the zero value (CommandPolicy{})
// restores the default, unrestricted behavior.
func SetCommandPolicy(policy CommandPolicy) {
	commandPolicyMu.Lock()
	defer commandPolicyMu.Unlock()
	commandPolicy = policy
}

// getCommandPolicy returns a snapshot of the active policy.
func getCommandPolicy() CommandPolicy {
	commandPolicyMu.RLock()
	defer commandPolicyMu.RUnlock()
	return commandPolicy
}

// checkCommandAllowed enforces AllowedCommands against a non-shell command
// line (e.g. os.Command's argument), panicking with a clear message if the
// resolved executable isn't on the allowlist. No-op when the policy's
// AllowedCommands is empty (the default).
func checkCommandAllowed(funcName, cmdLine string) {
	policy := getCommandPolicy()
	if len(policy.AllowedCommands) == 0 {
		return
	}
	parts := strings.Fields(cmdLine)
	if len(parts) == 0 {
		panic(fmt.Sprintf("%s: empty command", funcName))
	}
	if !policy.AllowedCommands[parts[0]] {
		panic(fmt.Sprintf("%s: command %q is not in the allowed command list", funcName, parts[0]))
	}
}

// checkShellAllowed enforces DisableShell for the sh-based execution
// builtins, panicking with a clear message if shell execution has been
// disabled by the embedder.
func checkShellAllowed(funcName string) {
	if getCommandPolicy().DisableShell {
		panic(fmt.Sprintf("%s: shell command execution is disabled by the host's CommandPolicy", funcName))
	}
}
