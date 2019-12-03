package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"
	"time"
)

// nolint: unparam
// Goal of this test is to instrument the main program
// so that we can extract the coverage for the end to end
// tests.
// This test is not expected to run during the unit test
// (eg make run)
func TestMain(t *testing.T) {
	if !strings.Contains(os.Args[0], "with-cov") {
		// regular tests, we do not want to run anything
		return
	}
	var (
		args []string
	)

	for k, arg := range os.Args {
		fmt.Printf("arg:%d - %s\n\n", k, arg)
		switch {
		case strings.HasPrefix(arg, "DEVEL"):
		case strings.HasPrefix(arg, "-test"):
		default:
			args = append(args, arg)
		}
	}

	waitCh := make(chan int, 1)

	os.Args = args
	go func() {

		main()
		close(waitCh)

	}()
	go func() {
		// wait a bit before print a message on how to run the e2e test
		time.Sleep(2 * time.Second)
		fmt.Println("** Coverage instructions: Once server is started, please run the e2e test to workout the coverage")
		fmt.Println("** Coverage instructions: Then Ctr-C when tests are done to flush the coverage profile")
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP)
	select {
	case s := <-signalCh:
		fmt.Printf("RECEIVED SIGNAL: %s\n", s)
		return
	case <-waitCh:
		return
	}
}
