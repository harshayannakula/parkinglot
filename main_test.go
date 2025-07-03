// main_test.go
package main

import (
	"testing"
	"time"
)

func TestParkCar_Success(t *testing.T) {
	lot := NewParkingLot("Lot A", 2)
	car := &Car{Number: "KA01AB1234", Color: "Red", Make: "Honda", Size: "small", IsHandicap: false}

	slot, err := lot.ParkCar(car)
	if err != nil {
		t.Fatalf("expected to park car successfully, got error: %v", err)
	}

	if slot != 1 {
		t.Errorf("expected car to be parked at slot 1, got slot %d", slot)
	}
}

func TestParkCar_FullLot(t *testing.T) {
	lot := NewParkingLot("Lot A", 1)
	car1 := &Car{Number: "KA01AB1234"}
	car2 := &Car{Number: "KA01AB5678"}

	_, _ = lot.ParkCar(car1)
	_, err := lot.ParkCar(car2)

	if err == nil {
		t.Errorf("expected error when parking in full lot, got nil")
	}
}
func TestUnparkCar_Success(t *testing.T) {
	lot := NewParkingLot("Lot A", 2)
	car := &Car{Number: "KA01AB9999"}
	_, _ = lot.ParkCar(car)

	slot, err := lot.UnparkCar("KA01AB9999")
	if err != nil {
		t.Fatalf("expected successful unpark, got error: %v", err)
	}
	if slot != 1 {
		t.Errorf("expected to unpark from slot 1, got %d", slot)
	}
}

func TestUnparkCar_NotFound(t *testing.T) {
	lot := NewParkingLot("Lot A", 1)
	_, err := lot.UnparkCar("NOTEXIST123")
	if err == nil {
		t.Error("expected error for car not found, got nil")
	}
}
func TestIsFull(t *testing.T) {
	lot := NewParkingLot("Lot A", 1)
	car := &Car{Number: "KA01F1111"}
	_, _ = lot.ParkCar(car)

	if !lot.IsFull() {
		t.Errorf("expected lot to be full")
	}
}
func TestNotifyObserversOnFull(t *testing.T) {
	called := false
	observer := func(msg string) {
		if msg == "FULL" {
			called = true
		}
	}

	lot := NewParkingLot("Lot A", 1)
	lot.Observers = []Observer{observer}

	car1 := &Car{Number: "A"}
	car2 := &Car{Number: "B"}
	_, _ = lot.ParkCarWithNotification(car1)
	_, _ = lot.ParkCarWithNotification(car2)

	if !called {
		t.Error("expected observer to be called on full lot")
	}
}
func TestNotifyObserversOnAvailable(t *testing.T) {
	called := false
	observer := func(msg string) {
		if msg == "AVAILABLE" {
			called = true
		}
	}

	lot := NewParkingLot("Lot A", 1)
	lot.Observers = []Observer{observer}

	car := &Car{Number: "KA01XX0001"}
	_, _ = lot.ParkCar(car)

	// Now unpark, should trigger "AVAILABLE"
	_, _ = lot.UnparkCarWithNotification("KA01XX0001")

	if !called {
		t.Error("expected observer to be called with AVAILABLE")
	}
}
func TestAttendantParksCar(t *testing.T) {
	lot := NewParkingLot("Lot A", 2)
	attendant := &Attendant{Name: "John", Lot: lot}
	car := &Car{Number: "KA09VV7777"}

	slot, err := attendant.ParkCarForDriver(car)
	if err != nil {
		t.Fatalf("attendant failed to park car: %v", err)
	}

	if slot != 1 {
		t.Errorf("expected slot 1, got %d", slot)
	}
}
func TestFindCar_Success(t *testing.T) {
	lot := NewParkingLot("Lot A", 2)
	car := &Car{Number: "KA01MM2323"}
	_, _ = lot.ParkCar(car)

	slot, err := lot.FindCar("KA01MM2323")
	if err != nil {
		t.Fatalf("expected to find car, got error: %v", err)
	}
	if slot.Number != 1 {
		t.Errorf("expected to find car at slot 1, got %d", slot.Number)
	}
}

func TestFindCar_NotFound(t *testing.T) {
	lot := NewParkingLot("Lot A", 1)
	_, err := lot.FindCar("UNKNOWN")
	if err == nil {
		t.Error("expected error for unknown car, got nil")
	}
}
func TestUnparkCarAndCharge(t *testing.T) {
	lot := NewParkingLot("Lot A", 1)
	car := &Car{Number: "KA01ZZ8888"}
	_, _ = lot.ParkCar(car)

	// Simulate some time passed
	car.ParkedAt = car.ParkedAt.Add(-3 * time.Minute)

	slot, fee, err := lot.UnparkCarAndCharge("KA01ZZ8888")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if slot != 1 {
		t.Errorf("expected to unpark from slot 1, got %d", slot)
	}
	if fee != 6 {
		t.Errorf("expected ₹6 charge (3 mins), got ₹%d", fee)
	}
}
