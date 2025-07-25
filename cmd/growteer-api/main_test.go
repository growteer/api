package main

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/test-go/testify/assert"
)

func TestMain(t *testing.T) {
	assert.NoError(t, os.Setenv("HTTP_PORT", "8888"))

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		assert.NotPanics(t, func() {
			main()
		})
	}()

	go func() {
		defer wg.Done()

		c := time.Tick(time.Second * 3)
		<-c
		assert.NoError(t, sendInterruptSignal(os.Getpid()))
	}()

	wg.Wait()
}
