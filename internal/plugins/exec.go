// story: e03s03
package plugins

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/danielvm-git/big-release/internal/algorithm"
)

// CommandRunner defines the interface for running commands.
type CommandRunner interface {
	Run(name string, args ...string) (string, error)
}

// OSCommandRunner is the production command runner using os/exec.
type OSCommandRunner struct{}

// Run executes a command and returns its output.
func (r *OSCommandRunner) Run(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command %q failed: %w\noutput: %s", name, err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

// ExecPlugin executes custom shell commands during release.
type ExecPlugin struct {
	runner CommandRunner
}

// NewExecPlugin creates a new ExecPlugin.
func NewExecPlugin() *ExecPlugin {
	return &ExecPlugin{
		runner: &OSCommandRunner{},
	}
}

// Name returns the plugin name.
func (p *ExecPlugin) Name() string {
	return "exec"
}

// VerifyConditions verifies that exec commands are configured.
func (p *ExecPlugin) VerifyConditions(ctx *algorithm.Context) error {
	if ctx.Config == nil {
		return fmt.Errorf("config is nil")
	}
	publisher, ok := ctx.Config.Publishers["exec"]
	if !ok || !publisher.Enabled {
		return fmt.Errorf("exec publisher not configured or not enabled")
	}
	commands, ok := publisher.Options["commands"]
	if !ok || strings.TrimSpace(commands) == "" {
		return fmt.Errorf("no commands configured in exec publisher options")
	}
	return nil
}

// AnalyzeCommits is not applicable for the exec plugin.
func (p *ExecPlugin) AnalyzeCommits(ctx *algorithm.Context) (algorithm.ReleaseType, error) {
	return "", nil
}

// GenerateNotes is not applicable for the exec plugin.
func (p *ExecPlugin) GenerateNotes(ctx *algorithm.Context) (string, error) {
	return "", nil
}

func (p *ExecPlugin) parseCommands(ctx *algorithm.Context) []string {
	publisher := ctx.Config.Publishers["exec"]
	commands := publisher.Options["commands"]
	lines := strings.Split(commands, "\n")
	var cmds []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			cmds = append(cmds, line)
		}
	}
	return cmds
}

// shellSplit splits a command line into tokens respecting double-quoted arguments.
func shellSplit(line string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false
	for _, r := range line {
		switch {
		case r == '"':
			inQuotes = !inQuotes
		case (r == ' ' || r == '\t') && !inQuotes:
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	return parts
}

func (p *ExecPlugin) runCommand(line string, idx int) error {
	parts := shellSplit(line)
	out, err := p.runner.Run(parts[0], parts[1:]...)
	if err != nil {
		return fmt.Errorf("command %d (%s) failed: %w", idx+1, line, err)
	}
	_ = out
	return nil
}

// Prepare runs all configured exec commands.
func (p *ExecPlugin) Prepare(ctx *algorithm.Context) error {
	if ctx.DryRun {
		return nil
	}
	commands := p.parseCommands(ctx)
	for i, line := range commands {
		if err := p.runCommand(line, i); err != nil {
			return err
		}
	}
	return nil
}

// Publish is not applicable for the exec plugin.
func (p *ExecPlugin) Publish(ctx *algorithm.Context) (*algorithm.Release, error) {
	return nil, nil
}

// Success is called after a successful release.
func (p *ExecPlugin) Success(ctx *algorithm.Context) error {
	return nil
}

// Fail is called on release failure.
func (p *ExecPlugin) Fail(ctx *algorithm.Context, err error) error {
	return nil
}

func init() {
	Register(NewExecPlugin())
}
