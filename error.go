package ngrok

type Error string

func (e Error) Error() string {
	return string(e)
}

const ErrBadExecutable Error = "the executable does not exist, or is a directory"

// ProcessFailed is returned from a process that starts successfully, but exits before it is deliberately stopped.
type ProcessFailed struct {
	Log       LogMessages
	ExitError error
}

func (p *ProcessFailed) Error() (str string) {
	str = p.Log.Error()
	if str == "" && p.ExitError != nil {
		str = p.ExitError.Error()
	}
	return
}
