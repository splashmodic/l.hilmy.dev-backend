package exithandler

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

type exitHandler []func()

func New(fns ...func()) *exitHandler {
	exitHandler := exitHandler{}
	for _, fn := range fns {
		exitHandler = append(exitHandler, fn)
	}
	return &exitHandler
}

func (e *exitHandler) Add(newFns ...func()) *exitHandler {
	exitHandler := *e
	for _, fn := range newFns {
		exitHandler = append(exitHandler, fn)
	}
	return &exitHandler
}

func (e *exitHandler) Watch() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGSEGV)
	go func() {
		<-c
		log.Println("clearing...")
		for _, fn := range *e {
			fn()
		}
		os.Exit(0)
	}()
}
