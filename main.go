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
	Number        int
	Row           string
	IsEmpty       bool
	Car           *Car
	AttendantName string
}

type ParkingLot struct {
	Name      string
	Slots     []Slot
	Observers []Observer
}

type Attendant struct {
	Name string
	Lot  *ParkingLot
}

type ParkingManager struct {
	Lots []*ParkingLot
}

type CarFilter struct {
	Color      string
	Make       string
	Size       string
	IsHandicap *bool // use pointer to distinguish unset vs false
}

type CarWithAttendant struct {
	Car
	Attendant string
	Row       string
}

func NewParkingLot(name string, capacity int) *ParkingLot {
	slots := make([]Slot, capacity)
	rowLetters := []string{"A", "B", "C", "D", "E"}

	for i := range slots {
		row := rowLetters[i%len(rowLetters)]
		slots[i] = Slot{
			Number:  i + 1,
			Row:     row,
			IsEmpty: true,
		}
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

func (a *Attendant) ParkCarForDriver(car *Car) (int, error) {
	fmt.Printf("Attendant %s is parking car %s\n", a.Name, car.Number)
	return a.Lot.ParkCarWithAttendant(car, a.Name)
}

func (pl *ParkingLot) ParkCarWithAttendant(car *Car, attendantName string) (int, error) {
	for i := range pl.Slots {
		if pl.Slots[i].IsEmpty {
			pl.Slots[i].Car = car
			pl.Slots[i].IsEmpty = false
			pl.Slots[i].AttendantName = attendantName
			car.ParkedAt = time.Now()
			return pl.Slots[i].Number, nil
		}
	}
	return -1, fmt.Errorf("lot is full")
}

func (pl *ParkingLot) FindCar(carNumber string) (*Slot, error) {
	for i := range pl.Slots {
		if !pl.Slots[i].IsEmpty && pl.Slots[i].Car.Number == carNumber {
			return &pl.Slots[i], nil
		}
	}
	return nil, fmt.Errorf("car %s not found in lot", carNumber)
}
func (pl *ParkingLot) UnparkCarAndCharge(carNumber string) (int, int, error) {
	for i := range pl.Slots {
		slot := &pl.Slots[i]
		if !slot.IsEmpty && slot.Car.Number == carNumber {
			duration := int(time.Since(slot.Car.ParkedAt).Minutes())
			if duration == 0 {
				duration = 1 // minimum charge for <1 minute
			}
			fee := duration * 2 // ₹2 per minute

			slot.Car = nil
			slot.IsEmpty = true
			pl.NotifyObservers("AVAILABLE")
			return slot.Number, fee, nil
		}
	}
	return -1, 0, fmt.Errorf("car not found")
}

func (pm *ParkingManager) ParkEvenly(car *Car) (string, int, error) {
	var targetLot *ParkingLot
	maxFree := -1

	for _, lot := range pm.Lots {
		freeSlots := 0
		for _, slot := range lot.Slots {
			if slot.IsEmpty {
				freeSlots++
			}
		}
		if freeSlots > maxFree {
			maxFree = freeSlots
			targetLot = lot
		}
	}

	if targetLot == nil {
		return "", -1, fmt.Errorf("all lots are full")
	}

	slotNum, err := targetLot.ParkCar(car)
	if err != nil {
		return "", -1, err
	}
	return targetLot.Name, slotNum, nil
}

func (a *Attendant) ParkCarWithStrategy(car *Car) (int, error) {
	if car.IsHandicap {
		// Handicap: nearest available slot (lowest slot number)
		for i := range a.Lot.Slots {
			if a.Lot.Slots[i].IsEmpty {
				return a.Lot.ParkCar(car) // default behavior already parks in lowest first
			}
		}
		return -1, fmt.Errorf("no available slot for handicap driver")
	}

	// Default strategy for non-handicap
	return a.ParkCarForDriver(car)
}

func (pm *ParkingManager) ParkLargeVehicle(car *Car) (string, int, error) {
	if car.Size != "large" {
		return "", -1, fmt.Errorf("not a large vehicle")
	}

	var targetLot *ParkingLot
	maxFree := -1

	for _, lot := range pm.Lots {
		free := 0
		for _, slot := range lot.Slots {
			if slot.IsEmpty {
				free++
			}
		}
		if free > maxFree {
			maxFree = free
			targetLot = lot
		}
	}

	if targetLot == nil {
		return "", -1, fmt.Errorf("no lot has space for large vehicle")
	}

	slot, err := targetLot.ParkCar(car)
	if err != nil {
		return "", -1, err
	}
	return targetLot.Name, slot, nil
}

func (pm *ParkingManager) FindCarsByColor(color string) []Car {
	var result []Car
	for _, lot := range pm.Lots {
		for _, slot := range lot.Slots {
			if !slot.IsEmpty && slot.Car.Color == color {
				result = append(result, *slot.Car)
			}
		}
	}
	return result
}

func (pm *ParkingManager) FindCars(filter CarFilter) []CarWithAttendant {
	var result []CarWithAttendant

	for _, lot := range pm.Lots {
		for _, slot := range lot.Slots {
			if slot.IsEmpty {
				continue
			}
			car := slot.Car
			if filter.Color != "" && car.Color != filter.Color {
				continue
			}
			if filter.Make != "" && car.Make != filter.Make {
				continue
			}
			if filter.Size != "" && car.Size != filter.Size {
				continue
			}
			if filter.IsHandicap != nil && car.IsHandicap != *filter.IsHandicap {
				continue
			}
			result = append(result, CarWithAttendant{
				Car:       *car,
				Attendant: slot.AttendantName,
			})
		}
	}

	return result
}

func (pm *ParkingManager) FindCarsParkedWithin(duration time.Duration) []CarWithAttendant {
	var result []CarWithAttendant
	cutoff := time.Now().Add(-duration)

	for _, lot := range pm.Lots {
		for _, slot := range lot.Slots {
			if !slot.IsEmpty && slot.Car.ParkedAt.After(cutoff) {
				result = append(result, CarWithAttendant{
					Car:       *slot.Car,
					Attendant: slot.AttendantName,
				})
			}
		}
	}
	return result
}

func (pm *ParkingManager) FindSmallHandicapInRowBOrD() []CarWithAttendant {
	var result []CarWithAttendant
	for _, lot := range pm.Lots {
		for _, slot := range lot.Slots {
			if !slot.IsEmpty && slot.Car.Size == "small" && slot.Car.IsHandicap &&
				(slot.Row == "B" || slot.Row == "D") {
				result = append(result, CarWithAttendant{
					Car:       *slot.Car,
					Attendant: slot.AttendantName,
					Row:       slot.Row,
				})
			}
		}
	}
	return result
}

func (pl *ParkingLot) GetAllParkedCars() []CarWithAttendant {
	var result []CarWithAttendant
	for _, slot := range pl.Slots {
		if !slot.IsEmpty {
			result = append(result, CarWithAttendant{
				Car:       *slot.Car,
				Attendant: slot.AttendantName,
				Row:       slot.Row,
			})
		}
	}
	return result
}
