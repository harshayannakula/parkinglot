// main.go
package main

import (
	"fmt"
	"time"
)

type Car struct {
	Number     string
	Color      string
	Make       string
	Size       string // "small" or "large"
	IsHandicap bool
	ParkedAt   time.Time
}

type Slot struct {
	Number  int
	IsEmpty bool
	Car     *Car
}

type ParkingLot struct {
	Name  string
	Slots []Slot
}

func NewParkingLot(name string, capacity int) *ParkingLot {
	slots := make([]Slot, capacity)
	for i := range slots {
		slots[i] = Slot{Number: i + 1, IsEmpty: true}
	}
	return &ParkingLot{Name: name, Slots: slots}
}

func (pl *ParkingLot) ParkCar(car *Car) (int, error) {
	for i := range pl.Slots {
		if pl.Slots[i].IsEmpty {
			pl.Slots[i].Car = car
			pl.Slots[i].IsEmpty = false
			car.ParkedAt = time.Now()
			return pl.Slots[i].Number, nil
		}
	}
	return -1, fmt.Errorf("parking lot is full")
}
