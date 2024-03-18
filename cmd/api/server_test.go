package main

import (
	"context"
	"filmoteka/internal/jsonlog"
	"os"
	"sync"
	"testing"
	"time"
)

func TestServe(t *testing.T) {
	app := &application{
		config: config{
			port: 8080,
			env:  "test",
		},
		logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		wg:     sync.WaitGroup{},
	}

	_, cancel := context.WithCancel(context.Background())

	go func() {
		if err := app.serve(); err != nil {
			t.Errorf("expected server to be closed, got %v", err)
		}
	}()

	time.Sleep(1 * time.Second)

	cancel()
}
