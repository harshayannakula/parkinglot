// main_test.go
package main

import (
	"testing"
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
