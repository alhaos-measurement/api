package model

import (
	"time"
)

type Measure struct {
	SensorID      int       `json:"sensorID"`
	MeasureTypeID int       `json:"measureTypeID"`
	UnitID        int       `json:"unitID"`
	Value         float64   `json:"value"`
	MeasuredAt    time.Time `json:"measuredAt,omitempty"`
}

type MeasureView struct {
	SensorName  string    `json:"sensorName"`
	MeasureName string    `json:"measureName"`
	UnitName    string    `json:"unitName"`
	Value       float64   `json:"value"`
	MeasuredAt  time.Time `json:"measuredAt,omitempty"`
}

type AvgMeasure struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Value float64   `json:"value"`
}
