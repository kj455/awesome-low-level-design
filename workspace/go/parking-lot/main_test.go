package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParkingLot(t *testing.T) {
	idGen := &SlotIDGenerator{}
	id1 := idGen.generate()
	id2 := idGen.generate()
	id3 := idGen.generate()
	id4 := idGen.generate()
	id5 := idGen.generate()
	id6 := idGen.generate()
	id7 := idGen.generate()
	id8 := idGen.generate()
	id9 := idGen.generate()
	parkingLot := ParkingLot{
		levels: []ParkingLevel{
			{
				level: 0,
				slots: []*ParkingSlot{
					{
						id:    id1,
						types: []VehicleType{VehicleCar},
					},
					{
						id:    id2,
						types: []VehicleType{VehicleCar},
					},
					{
						id:    id3,
						types: []VehicleType{VehicleCar},
					},
				},
			},
			{
				level: 1,
				slots: []*ParkingSlot{
					{
						id:    id4,
						types: []VehicleType{VehicleMotorCycle},
					},
					{
						id:    id5,
						types: []VehicleType{VehicleMotorCycle},
					},
					{
						id:    id6,
						types: []VehicleType{VehicleMotorCycle},
					},
				},
			},
			{
				level: 2,
				slots: []*ParkingSlot{
					{
						id:    id7,
						types: []VehicleType{VehicleTruck},
					},
					{
						id:    id8,
						types: []VehicleType{VehicleTruck},
					},
					{
						id:    id9,
						types: []VehicleType{VehicleTruck},
					},
				},
			},
		},
	}

	v1 := &Car{id: VehicleID(1)}
	v2 := &MotorCycle{id: VehicleID(2)}
	v3 := &Truck{id: VehicleID(3)}

	v1ID, err := parkingLot.entry(v1)
	assert.NoError(t, err)
	assert.Equal(t, id1, v1ID)
	v2ID, err := parkingLot.entry(v2)
	assert.NoError(t, err)
	assert.Equal(t, id4, v2ID)
	v3ID, err := parkingLot.entry(v3)
	assert.NoError(t, err)
	assert.Equal(t, id7, v3ID)

	assert.Equal(t, `X..
X..
X..
`, parkingLot.formatSlotStatus())

	assert.NoError(t, parkingLot.exit(v1ID))
	assert.NoError(t, parkingLot.exit(v2ID))
	assert.NoError(t, parkingLot.exit(v3ID))
}
