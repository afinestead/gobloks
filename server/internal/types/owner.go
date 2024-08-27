package types

import (
	"strconv"
)

func (o1 Owner) IsSamePlayer(o2 Owner) bool {
	return o1&PLAYER_MASK == o2&PLAYER_MASK
}

func (o Owner) IsVacant() bool {
	return o&RESERVED == RESERVED || o&VACANT == VACANT
}

func (o Owner) IsOrigin() bool {
	return o&ORIGIN == ORIGIN
}

func (o Owner) ToString() string {
	var s string
	if o&RESERVED != 0 {
		s = "#"
	} else if o&VACANT != 0 && o&ORIGIN == 0 {
		s = " "
	} else {
		s = strconv.Itoa(int(o & PLAYER_MASK))
	}

	if o&ORIGIN != 0 {
		return "[" + s + "]"
	} else {
		return " " + s + " "
	}
}
