// main.go
package main

import (
	"fmt"
	"time"
)

type Observer func(msg string)

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
	Name      string
	Slots     []Slot
	Observers []Observer
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
func (pl *ParkingLot) UnparkCar(carNumber string) (int, error) {
	for i := range pl.Slots {
		if !pl.Slots[i].IsEmpty && pl.Slots[i].Car.Number == carNumber {
			pl.Slots[i].Car = nil
			pl.Slots[i].IsEmpty = true
			return pl.Slots[i].Number, nil
		}
	}
	return -1, fmt.Errorf("car not found")
}
func (pl *ParkingLot) IsFull() bool {
	for _, slot := range pl.Slots {
		if slot.IsEmpty {
			return false
		}
	}
	return true
}

func (pl *ParkingLot) NotifyObservers(message string) {
	for _, observer := range pl.Observers {
		observer(message)
	}
}

func (pl *ParkingLot) ParkCarWithNotification(car *Car) (int, error) {
	slot, err := pl.ParkCar(car)
	if err != nil {
		pl.NotifyObservers("FULL")
	}
	return slot, err
}

func (pl *ParkingLot) UnparkCarWithNotification(carNumber string) (int, error) {
	for i := range pl.Slots {
		if !pl.Slots[i].IsEmpty && pl.Slots[i].Car.Number == carNumber {
			pl.Slots[i].Car = nil
			pl.Slots[i].IsEmpty = true
			pl.NotifyObservers("AVAILABLE")
			return pl.Slots[i].Number, nil
		}
	}
	return -1, fmt.Errorf("car not found")
}

type Attendant struct {
	Name string
	Lot  *ParkingLot
}

func (a *Attendant) ParkCarForDriver(car *Car) (int, error) {
	fmt.Printf("Attendant %s is parking car %s\n", a.Name, car.Number)
	return a.Lot.ParkCar(car)
}
