package model

// LastSensorMeasure REST api request
type LastSensorMeasure struct {
	SensorID      int `json:"sensorID"`
	MeasureTypeID int `json:"measureTypeID"`
}
