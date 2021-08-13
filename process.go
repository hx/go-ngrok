package ngrok

import (
	"context"
	"os"
	"os/exec"
	"sync/atomic"
)

// Process represents an Executable process.
type Process struct {
	Dispatcher

	cmd      *exec.Cmd
	stream   messageStream
	stopping int64
	done     chan struct{}
	err      error
	log      LogMessages
}

// Start starts the process, and returns when it is running. Use Wait to rejoin the process.
func (p *Process) Start() (err error) {
	p.done = make(chan struct{})
	p.stream.dispatcher = &p.Dispatcher
	p.cmd.Stdout = &p.stream
	p.Subscribe(func(message *LogMessage) { p.log = append(p.log, message) })
	if err = p.cmd.Start(); err != nil {
		return
	}
	go func() {
		p.err = p.cmd.Wait()
		if p.err != nil || (p.log.FinalError() != nil && atomic.LoadInt64(&p.stopping) == 0) {
			p.err = &ProcessFailed{Log: p.log, ExitError: p.err}
		}
		p.stream.dispatcher.release()
		close(p.done)
	}()
	return
}

// StopContext interrupts the process. If ctx expires before the process completes, the process will be killed. On
// operating systems that do not support the interrupt signal (such as Windows), the process will be killed immediately.
func (p *Process) StopContext(ctx context.Context) (err error) {
	atomic.StoreInt64(&p.stopping, 1)
	if err = p.cmd.Process.Signal(os.Interrupt); err != nil {
		err = p.cmd.Process.Kill()
	} else {
		select {
		case <-p.done:
		case <-ctx.Done():
			err = p.cmd.Process.Kill()
		}
	}

	if err == nil {
		err = p.Wait()
	}

	return
}

// Stop interrupts the process, and returns its exit error, if any.
func (p *Process) Stop() error {
	return p.StopContext(context.Background())
}

// Wait blocks and returns when the process exits.
func (p *Process) Wait() error {
	<-p.done
	return p.err
}
