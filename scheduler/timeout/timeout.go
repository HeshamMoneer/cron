package timeout

import (
	"fmt"
	"time"
)

var (
	ErrTimedOut = fmt.Errorf("function call timed out")
)

func Run(fn func() error, timeout time.Duration) error {
	c := make(chan error, 1)
	go func() {
		c <- fn()
		close(c)
	}()
	t := time.NewTimer(timeout)
	defer t.Stop()
	select {
	case err := <-c:
		return err
	case <-t.C:
		return ErrTimedOut
	}
}
