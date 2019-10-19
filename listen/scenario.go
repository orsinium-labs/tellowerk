package listen

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/francoispqt/onelog"
)

// Scenario reads user input from stdin
type Scenario struct {
	logger *onelog.Logger
	lines  []string
	pos    int
}

// Close does nothing for Scenario and exists to satisfy interface
func (s *Scenario) Close() error {
	return nil
}

// Listen reads user input from stdin
func (s *Scenario) Listen() string {
	s.logger.Debug("eading scenarion")
	for s.pos < len(s.lines) {
		line := s.lines[s.pos]
		s.pos++
		if line == "" {
			continue
		}
		if strings.HasSuffix(line, "sleep ") {
			secs, _ := strconv.Atoi(line[6:])
			time.Sleep(time.Duration(secs) * time.Second)
			continue
		}
		return line
	}
	return "end"
}

// NewScenario creates a new Scenario instance to get commands from scenario file
func NewScenario(config ListenConfig, logger *onelog.Logger) (*Scenario, error) {
	s := &Scenario{logger: logger, lines: make([]string, 0)}

	file, err := os.Open(config.Scenario)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() { // internally, it advances token based on sperator
		s.lines = append(s.lines, scanner.Text())
	}
	return s, nil
}
