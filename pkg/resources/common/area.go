package common

import (
	"encoding/json"
	"fmt"
)

type Area int

const (
	LivingRoom Area = iota
	Kitchen
	Dining
	Bedroom
	KidsBedroom
	Bathroom
	Nursery
	Recreation
	Office
	Gym
	Hallway
	Toilet
	FrontDoor
	Garage
	Terrace
	Garden
	Driveway
	Carport
	Home
	Downstairs
	Upstairs
	TopFloor
	Attic
	GuestRoom
	Staircase
	Lounge
	ManCave
	Computer
	Studio
	Music
	TV
	Reading
	Closet
	Storage
	LaundryRoom
	Balcony
	Porch
	Barbecue
	Pool
	Other
	InvalidArea
)

var AreaNames = [...]string{
	"living_room",
	"kitchen",
	"dining",
	"bedroom",
	"kids_bedroom",
	"bathroom",
	"nursery",
	"recreation",
	"office",
	"gym",
	"hallway",
	"toilet",
	"front_door",
	"garage",
	"terrace",
	"garden",
	"driveway",
	"carport",
	"home",
	"downstairs",
	"upstairs",
	"top_floor",
	"attic",
	"guest_room",
	"staircase",
	"lounge",
	"man_cave",
	"computer",
	"studio",
	"music",
	"tv",
	"reading",
	"closet",
	"storage",
	"laundry_room",
	"balcony",
	"porch",
	"barbecue",
	"pool",
	"other",
	"invalid_area",
}

// String returns the original snake-case token (or "Area(<n>)" if out of range).
func (a *Area) String() string {
	var s string
	if a == nil {
		return AreaNames[InvalidArea]
	}
	index := int(*a)
	if index < 0 || index >= len(AreaNames) {
		return fmt.Sprintf("Area(%d)", a)
	}
	s = AreaNames[index]
	return s
}

// UnmarshalJSON implements json.Unmarshaler interface for Area
func (a *Area) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := ParseArea(s)
	if err != nil {
		return err
	}
	*a = parsed
	return nil
}

// MarshalJSON implements json.Marshaler interface for Area
func (a *Area) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

// ParseArea converts a string like "living_room" back into the enum.
func ParseArea(s string) (Area, error) {
	for i, name := range AreaNames {
		if name == s {
			return Area(i), nil
		}
	}
	return InvalidArea, fmt.Errorf("unknown area: %q", s)
}
