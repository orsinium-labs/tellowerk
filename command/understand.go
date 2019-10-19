package command

import (
	"regexp"
	"strconv"
	"strings"
)

// Understand parses text into Command
func Understand(text string) Command {
	rex := regexp.MustCompile("\\w+")
	words := rex.FindAllString(text, -1)
	command := Command{}
	command.Action = getAction(words)
	if command.Action == "" {
		return command
	}
	command.Distance = getDistance(words)
	command.Units = getUnits(words)
	return command
}

func getAction(words []string) Action {
	for i, word := range words {
		if word != "turn" {
			continue
		}
		switch strings.ToLower(words[i+1]) {
		case "left":
			return TurnLeft
		case "right":
			return TurnRight
		}
	}

	for _, word := range words {
		switch strings.ToLower(word) {
		case "start", "begin", "launch":
			return Start
		case "land":
			return Land
		case "halt", "close", "end":
			return Halt
		case "hover", "stop", "wait":
			return Hover
		case "left":
			return Left
		case "right":
			return Right
		case "back":
			return Back
		case "front":
			return Front
		case "down":
			return Down
		case "up":
			return Up
		}
	}
	return ""
}

func getDistance(words []string) int {
	for _, word := range words {
		number, _ := strconv.Atoi(word)
		if number > 0 {
			return number
		}
	}
	return 0
}

func getUnits(words []string) DistanceUnits {
	for _, word := range words {
		switch strings.ToLower(word) {
		case "degrees", "degree", "d":
			return Degrees
		case "meter", "meters", "m":
			return Meters
		case "second", "seconds", "sec":
			return Seconds
		}
	}
	return 0
}
