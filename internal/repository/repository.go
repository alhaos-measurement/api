package repository

import (
	"github.com/alhaos-measurement/api/internal/model"
	"github.com/jackc/pgx"
)

type Repository struct {
	pool *pgx.ConnPool
}

func New(db *pgx.ConnPool) *Repository {
	return &Repository{pool: db}
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

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetLastMeasure(sensorID int) (*model.Measure, error) {

	tx, err := r.pool.Begin()
	if err != nil {
		return nil, err
	}

	const query = `
SELECT s.name,
       mt.name,
       u.name,
       m.value,
       m.measured_at
  FROM (
    SELECT sensor_id,
           measure_type_id,
           unit_id,
           value,
           measured_at,
           ROW_NUMBER() OVER (ORDER BY measured_at DESC) AS rn
      FROM measurements
     WHERE sensor_id = $1
  ) m
  JOIN measure_type mt 
    ON m.measure_type_id = mt.measure_type_id
   AND m.rn = 1
  JOIN units u 
    ON m.unit_id = u.unit_id
  JOIN sensors s
    ON m.sensor_id = s.sensor_id`
	row := tx.QueryRow(query, sensorID)

	var m model.Measure

	err = row.Scan(&m.SensorID, &m.MeasureTypeID, &m.UnitID, &m.Value, &m.MeasuredAt)
	if err != nil {
		return nil, err
	}

	return &m, nil
}
