package main

import (
	"errors"
	"slices"
	"sync"
	"sync/atomic"
)

type ParkingLot struct {
	levels []ParkingLevel
}

func (p *ParkingLot) entry(v *Vehicle) (SlotID, error) {
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

type ParkingLevel struct {
	level int
	slots []*ParkingSlot
}

type SlotID int

type ParkingSlot struct {
	id        SlotID
	types     []VehicleType
	vehicleID *VehicleID
	mu        sync.Mutex
}

func (s *ParkingSlot) canAccommodate(typ VehicleType) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.vehicleID == nil && slices.Contains(s.types, typ)
}

func (s *ParkingSlot) tryAccommodate(v *Vehicle) (SlotID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.vehicleID != nil || !slices.Contains(s.types, v.typ) {
		return 0, errors.New("not available")
	}
	s.vehicleID = &v.id
	return s.id, nil
}

func (s *ParkingSlot) accommodate(v *Vehicle) (SlotID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.vehicleID != nil {
		return 0, errors.New("conflict")
	}
	s.vehicleID = &v.id
	return s.id, nil
}

func (s *ParkingSlot) exit() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.vehicleID == nil {
		return errors.New("already exit")
	}
	s.vehicleID = nil
	return nil
}

type VehicleID int

type Vehicle struct {
	id  VehicleID
	typ VehicleType
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
