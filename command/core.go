package command

// Action means what to do
type Action string

// DistanceUnits show units of distance in a command (meters, degrees, seconds)
type DistanceUnits int8

// Command contains action type (what to do) and distance info (how long to do)
type Command struct {
	Action   Action
	Distance int
	Units    DistanceUnits
}

const (
	// Start means take of and hover
	Start Action = "start"
	// Land means stop moving and land
	Land Action = "land"
	// Hover means stop doing anything and hover
	Hover Action = "hover"

	// -- MOVEMENTS -- //

	// Left is negative ox movement
	Left Action = "left"
	// Right is positive ox movement
	Right Action = "right"
	// Back is negative oy movement
	Back Action = "back"
	// Front is positive oy movement
	Front Action = "front"
	// Down is negative oz movement
	Down Action = "down"
	// Up is positive oz movement
	Up Action = "up"

	// -- TURNS -- //

	// TurnLeft means start rotating counter clockwise
	TurnLeft Action = "turn left"
	// TurnRight means start rotating clockwise
	TurnRight Action = "turn right"

	// -- DISTANCE TYPES -- //

	// Degrees is a rotation distance unit
	Degrees DistanceUnits = iota + 1
	// Meters is a movement distance unit
	Meters DistanceUnits = iota
	// Seconds is how long to do an action
	Seconds DistanceUnits = iota
)
