package main

import (
	"errors"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
)

type ParkingLot struct {
	levels []ParkingLevel
}

func (p *ParkingLot) entry(v Vehicle) (SlotID, error) {
	for _, level := range p.levels {
		for _, slot := range level.slots {
			if sl, err := slot.tryAccommodate(v); err == nil {
				return sl, nil
			}
		}
	}
	return 0, errors.New("not available now")
}

func (p *ParkingLot) exit(id SlotID) error {
	for _, level := range p.levels {
		for _, slot := range level.slots {
			if slot.id == id {
				return slot.exit()
			}
		}
	}
	return errors.New("slot not found")
}

func (p *ParkingLot) getAvailableSlots(typ VehicleType) []SlotID {
	available := make([]SlotID, 0)
	for _, level := range p.levels {
		for _, slot := range level.slots {
			if slot.canAccommodate(typ) {
				available = append(available, slot.id)
			}
		}
	}
	return available
}

func (p *ParkingLot) formatSlotStatus() string {
	const (
		empty = "."
		full  = "X"
	)
	var status strings.Builder
	for _, level := range p.levels {
		for _, slot := range level.slots {
			slot.mu.Lock()
			if slot.vehicleID == 0 {
				status.WriteString(empty)
			} else {
				status.WriteString(full)
			}
			slot.mu.Unlock()
		}
		status.WriteString("\n")
	}
	return status.String()
}

type ParkingLevel struct {
	level int
	slots []*ParkingSlot
}

type SlotID int

type ParkingSlot struct {
	id        SlotID
	types     []VehicleType
	vehicleID VehicleID
	mu        sync.Mutex
}

func (s *ParkingSlot) canAccommodate(typ VehicleType) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.vehicleID == 0 && slices.Contains(s.types, typ)
}

func (s *ParkingSlot) tryAccommodate(v Vehicle) (SlotID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.vehicleID != 0 || !slices.Contains(s.types, v.Type()) {
		return 0, errors.New("not available")
	}
	s.vehicleID = v.ID()
	return s.id, nil
}

func (s *ParkingSlot) accommodate(v Vehicle) (SlotID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.vehicleID != 0 {
		return 0, errors.New("conflict")
	}
	s.vehicleID = v.ID()
	return s.id, nil
}

func (s *ParkingSlot) exit() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.vehicleID == 0 {
		return errors.New("already exit")
	}
	s.vehicleID = 0
	return nil
}

type VehicleID int

type Car struct {
	id VehicleID
}

func (c *Car) ID() VehicleID {
	return c.id
}

func (c *Car) Type() VehicleType {
	return VehicleCar
}

type MotorCycle struct {
	id VehicleID
}

func (m *MotorCycle) ID() VehicleID {
	return m.id
}

func (m *MotorCycle) Type() VehicleType {
	return VehicleMotorCycle
}

type Truck struct {
	id VehicleID
}

func (t *Truck) ID() VehicleID {
	return t.id
}

func (t *Truck) Type() VehicleType {
	return VehicleTruck
}

type Vehicle interface {
	ID() VehicleID
	Type() VehicleType
}

type VehicleType int

const (
	VehicleUnknown VehicleType = iota
	VehicleCar
	VehicleMotorCycle
	VehicleTruck
)

type SlotIDGenerator struct {
	lastID int64
}

func (g *SlotIDGenerator) generate() SlotID {
	return SlotID(atomic.AddInt64(&g.lastID, 1))
}

func main() {

}
