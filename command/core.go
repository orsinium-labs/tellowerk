package command

type Action string
type DistanceUnits int8

type Command struct {
	Action   Action
	Distance int
	Units    DistanceUnits
}

const (
	Start Action = "start" // take of and hover
	Land  Action = "land"  // stop moving and land
	Hover Action = "hover" // stop doing anything and hover

	// movements
	Left  Action = "left"
	Right Action = "right"
	Back  Action = "back"
	Front Action = "front"
	Up    Action = "up"
	Down  Action = "down"

	// turns
	TurnLeft  Action = "turn left"  // counter clockwise
	TurnRight Action = "turn right" // clockwise

	// distance types
	Degrees DistanceUnits = iota + 1
	Meters  DistanceUnits = iota
	Seconds DistanceUnits = iota
)
