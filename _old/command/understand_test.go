package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnderstand(t *testing.T) {
	f := func(text string, action Action, distance int, units DistanceUnits) {
		c := Understand(text)
		assert.Equal(t, c.Action, action, "invalid action type")
		assert.Equal(t, c.Distance, distance, "invalid distance")
		assert.Equal(t, c.Units, units, "invalid units")
	}

	f("unknown action", "", 0, 0)
	f("start", Start, 0, 0)
	f("let's start", Start, 0, 0)
	f("land", Land, 0, 0)
	f("end up", Land, 0, 0)
	f("stop", Hover, 0, 0)
	f("hover", Hover, 0, 0)
	f("left 2 meters", Left, 2, Meters)
	f("turn right 180 degrees", TurnRight, 180, Degrees)
}
