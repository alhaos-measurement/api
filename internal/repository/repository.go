package repository

import (
	"github.com/alhaos-measurement/api/internal/logger"
	"github.com/alhaos-measurement/api/internal/model"
	"github.com/jackc/pgx"
	"time"
)

type Repository struct {
	pool   *pgx.ConnPool
	logger *logger.Logger
}

func New(db *pgx.ConnPool, logger *logger.Logger) *Repository {
	return &Repository{pool: db, logger: logger}
}

func (r *Repository) AddMeasure(measure *model.Measure) error {
	tx, err := r.pool.Begin()
	if err != nil {
		return err
	}

	const query = `
INSERT INTO measurements (sensor_id, measure_type_id, unit_id, value, measured_at)
VALUES (
  $1,
  $2,
  $3,
  $4,
  NOW()
)`

	_, err = tx.Exec(
		query,
		measure.SensorID,
		measure.MeasureTypeID,
		measure.UnitID,
		measure.Value,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	const queryUpdate = `
update measurements_current
   set value = $1,
       measured_at = now(),
       unit_id = $2 
 where sensor_id = $3
   and measure_type_id = $4`

	result, err := tx.Exec(queryUpdate, measure.Value, measure.UnitID, measure.SensorID, measure.MeasureTypeID)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			panic(err)
		}
		return err
	}

	if result.RowsAffected() == 0 {
		const queryInsert = `
INSERT INTO measurements_current (sensor_id, measure_type_id, unit_id, value, measured_at)
VALUES ( $1, $2, $3, $4, NOW())`

		_, err = tx.Exec(queryInsert, measure.SensorID, measure.MeasureTypeID, measure.UnitID, measure.Value)
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				panic(err)
			}
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetLastMeasure(sensorID int, measureTypeID int) (*model.MeasureView, error) {

	const query = `
SELECT s.name,
       mt.name,
       u.name,
       m.value,
       m.measured_at
  FROM (
  	SELECT * 
  	  FROM measurements_current
  	 WHERE sensor_id = $1
       AND unit_id = $2
 	) m
  JOIN measure_type mt 
    ON m.measure_type_id = mt.measure_type_id
  JOIN units u 
    ON m.unit_id = u.unit_id
  JOIN sensors s
    ON m.sensor_id = s.sensor_id
`
	row := r.pool.QueryRow(query, sensorID, measureTypeID)

	var m model.MeasureView

	err := row.Scan(&m.SensorName, &m.MeasureName, &m.UnitName, &m.Value, &m.MeasuredAt)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

// Units return all units from database
func (r *Repository) Units() ([]model.Unit, error) {

	const query = "SELECT unit_id, name FROM units"

	rows, err := r.pool.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []model.Unit
	var u model.Unit

	for rows.Next() {
		err := rows.Scan(&u.ID, &u.Name)
		if err != nil {
			return nil, err
		}
		units = append(units, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return units, nil
}

// AvgPressureHourly
// returns average pressure hourly last 12 hours
func (r *Repository) AvgPressureHourly() ([]model.AvgMeasure, error) {

	const query = "SELECT hour_bucket, avg FROM pressure_hourly_avg_12h"

	rows, err := r.pool.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var hour_backet time.Time
	var avg float64
	var avgMeasureValues []model.AvgMeasure

	for rows.Next() {
		err := rows.Scan(&hour_backet, &avg)
		if err != nil {
			return nil, err
		}
		avgPressure := model.AvgMeasure{
			Start: hour_backet,
			End:   hour_backet.Add(time.Hour),
			Value: avg,
		}
		avgMeasureValues = append(avgMeasureValues, avgPressure)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return avgMeasureValues, nil
}

// PurgeOldMeasurements purge table measurements
func (r *Repository) PurgeOldMeasurements() {

	const query = "DELETE FROM measurements WHERE measured_at < date_trunc('day', NOW()) - INTERVAL '7 days';"

	_, _ = r.pool.Exec(query)

}
