package main

import (
	"log/slog"
	"os"
)

const file = "app.log"

var logger *slog.Logger

func init() {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}
	logger = slog.New(slog.NewTextHandler(f, nil))
}
