package input

import (
	"bufio"
	"os"
	"time"
)

const defaultTimeOut = 30 * time.Second

func WaitEnter(timeout time.Duration) {
	if timeout == 0 {
		timeout = defaultTimeOut
	}

	input := make(chan bool)
	done := make(chan bool)

	go func() {
		if bufio.NewScanner(os.Stdin).Scan() {
			input <- true
		}
	}()
	go func() {
		time.Sleep(timeout)
		done <- true
	}()

	select {
	case <-input:
		return
	case <-done:
		return
	}
}
