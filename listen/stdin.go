package listen

import (
	"fmt"

	"github.com/francoispqt/onelog"
)

// Console reads user input from stdin
type Console struct {
	logger *onelog.Logger
}

// Close does nothing for Console and exists to satisfy interface
func (in *Console) Close() error {
	return nil
}

// Listen reads user input from stdin
func (in *Console) Listen() string {
	in.logger.Debug("waiting for user input")
	var input string
	fmt.Scanln(&input)
	return input
}

// NewConsole creates a new Console instance to get user input from console
func NewConsole(config ListenConfig, logger *onelog.Logger) (*Console, error) {
	return &Console{logger: logger}, nil
}
