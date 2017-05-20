package shell

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
)

type ReadLiner interface {
	ReadLine(context.Context) (context.Context, string, error)
}

func (it *Cmd) Loop(readline1 ReadLiner) error {
	sigint := make(chan os.Signal, 1)
	defer close(sigint)
	quit := make(chan struct{}, 1)
	defer close(quit)

	for {
		ctx, cancel := context.WithCancel(context.Background())
		ctx, line, err := readline1.ReadLine(ctx)

		if err != nil {
			cancel()
			return err
		}
		signal.Notify(sigint, os.Interrupt)

		go func(sigint_ chan os.Signal, quit_ chan struct{}, cancel_ func()) {
			for {
				select {
				case <-sigint_:
					cancel_()
					<-quit
					return
				case <-quit:
					cancel_()
					return
				}
			}
		}(sigint, quit, cancel)
		_, err = it.InterpretContext(ctx, line)
		signal.Stop(sigint)
		quit <- struct{}{}

		if err != nil {
			if err == io.EOF {
				break
			}
			if err1, ok := err.(AlreadyReportedError); ok {
				if err1.Err == io.EOF {
					break
				}
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
	return nil
}
