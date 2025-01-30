package model

type DeviceType int

const (
	PC DeviceType = iota
	MUSIC_BOX
	TOGGLE_SWITCH
)

var DeviceTypeName = map[DeviceType]string{
	PC:            "PC",
	MUSIC_BOX:     "MUSIC_BOX",
	TOGGLE_SWITCH: "TOGGLE_SWITCH",
}

func (t DeviceType) String() string {
	return DeviceTypeName[t]
}

type State int

const (
	DISCONNECTED State = iota
	CONNECTED
	IDLE
	PROCESSING
)

var StateName = map[State]string{
	DISCONNECTED: "DISCONNECTED",
	CONNECTED:    "CONNECTED",
	IDLE:         "IDLE",
	PROCESSING:   "PROCESSING",
}

func (s *State) String() string {
	return StateName[*s]
}

func (s *State) changeState(state State) {
	*s = state
}
