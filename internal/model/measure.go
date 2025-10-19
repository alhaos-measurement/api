package model

import "time"

type Measure struct {
	SensorID      int       `json:"sensorID"`
	MeasureTypeID int       `json:"measureTypeID"`
	UnitID        int       `json:"unitID"`
	Value         float64   `json:"value"`
	MeasuredAt    time.Time `json:"measuredAt,omitempty"`
}
