package cli

import (
	"io"
)

// CommandExecutor definition to avoid dependency on main package
type CommandExecutor interface {
	RunCommand(name string, args []string, env []string, stdout, stderr io.Writer) error
}

type AICLI interface {
	Install(executor CommandExecutor) error
	Auth(executor CommandExecutor, token string) error
	Run(executor CommandExecutor, prompt string, model string) error
}
