package repository

import (
	"fmt"
	"github.com/alhaos-measurement/api/internal/model"
	"github.com/jackc/pgx"
	"testing"
	"time"
)

func TestRepository_GetLastMeasure(t *testing.T) {

	data := []struct {
		sensorID int
	}{
		{
			sensorID: 1,
		},
	}

	// Init connection pool
	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "alhaos.online",
			Port:     5432,
			Database: "measurements",
			User:     "measurements_owner",
			Password: "kmsX*c5MBK1d",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	r := New(pool)

	if err != nil {
		panic(err)
	}

	for _, datum := range data {
		measure, err := r.GetLastMeasure(datum.sensorID)
		if err != nil {
			t.Fatal(err)
		}
		PrintMeasure(measure)
	}
}

func PrintMeasure(measure *model.Measure) {
	fmt.Printf("Sensor ID: %d\n", measure.SensorID)
	fmt.Printf("Value: %f\n", measure.Value)
	fmt.Printf("MeasuredAt: %s\n", measure.MeasuredAt)
}

func TestRepository_AddMeasure(t *testing.T) {

	data := []struct {
		measure *model.Measure
	}{
		{
			measure: &model.Measure{
				SensorID:      1,
				MeasureTypeID: 1,
				UnitID:        1,
				Value:         15.66,
				MeasuredAt:    time.Now().UTC(),
			},
		},
	}

	// Init connection pool
	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "alhaos.online",
			Port:     5432,
			Database: "measurements",
			User:     "measurements_owner",
			Password: "kmsX*c5MBK1d",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	r := New(pool)

	for _, datum := range data {
		err = r.AddMeasure(datum.measure)
		if err != nil {
			t.Fatal(err)
		}
	}
}
