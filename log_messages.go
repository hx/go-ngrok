package ngrok

// LogMessages is a series of messages emitted by a Process.
type LogMessages []*LogMessage

// FinalError returns the last LogMessage containing a non-blank Error in the receiving series.
func (l LogMessages) FinalError() *LogMessage {
	for i := len(l) - 1; i >= 0; i-- {
		msg := l[i]
		if msg.Error != "" {
			return msg
		}
	}
	return nil
}

func (l LogMessages) Error() string {
	if err := l.FinalError(); err != nil {
		return err.Error
	}
	return ""
}
