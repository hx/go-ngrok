package ngrok

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// Executable represents an ngrok executable.
type Executable string

// IsAvailable is true when the Executable exists and has the necessary permissions to be executed.
func (e Executable) IsAvailable() bool {
	_, err := e.Version()
	return err == nil
}

// Path returns the path to the Executable.
func (e Executable) Path() string {
	return string(e)
}

// String returns the path to the Executable.
func (e Executable) String() string {
	return string(e)
}

// NewProcess creates a new Process with the given arguments.
func (e Executable) NewProcess(args ...string) *Process {
	return &Process{cmd: e.command(append(args, "--log", "stdout", "--log-format", "json")...)}
}

// Version returns the numeric part of the executable's --version output.
func (e Executable) Version() (*Version, error) {
	if stat, err := os.Stat(e.Path()); err != nil || stat.IsDir() {
		return nil, ErrBadExecutable
	}

	result, err := e.command("--version").CombinedOutput()
	if err != nil {
		return nil, err
	}

	str := strings.TrimSpace(string(result))
	if v := ParseVersion(str); v != nil {
		return v, nil
	}

	return nil, errors.New(str)
}

func (e Executable) command(args ...string) *exec.Cmd {
	return exec.Command(e.Path(), args...)
}
