package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/lbwise/logger"
)

// Manual Test with terminal for colors
func main() {
	buf := new(bytes.Buffer)
	l := logger.New(
		logger.WithOutput(buf),
		logger.WithColor(),
		logger.WithPrefix("consensus"),
		logger.WithTimeIncluded(),
	)

	err := l.Warn("This is a colored low level log", logger.WarnSevere)
	err = l.Log("This is a normal blue log")
	err = l.Error(errors.New("something screwed up"), 500)
	fmt.Println(buf.String())
	if err != nil {
		fmt.Println("Error logging:", err)
	}
}
