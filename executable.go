package ngrok

import (
	"errors"
	"os"
	"os/exec"
	"regexp"
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

var versionPattern = regexp.MustCompile(`\d+(\.\d+)+`)

// Version returns the numeric part of the executable's --version output.
func (e Executable) Version() (string, error) {
	if stat, err := os.Stat(e.Path()); err != nil || stat.IsDir() {
		return "", ErrBadExecutable
	}
	var (
		result, err = e.command("--version").CombinedOutput()
		str         = strings.TrimSpace(string(result))
	)
	if match := versionPattern.FindString(str); match != "" {
		str = match
	} else {
		err = errors.New(str)
		str = ""
	}
	return str, err
}

func (e Executable) command(args ...string) *exec.Cmd {
	return exec.Command(e.Path(), args...)
}
