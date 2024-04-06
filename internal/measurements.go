package internal

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type Measurement struct {
	Value float64 `json:"value"`
	Type  string  `json:"type"`
}

func PostMeasurements(db *Database) func(c echo.Context) error {
	return func(c echo.Context) error {

		var measurement Measurement
		err := c.Bind(&measurement)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}
		_, err = db.conn.Exec(context.Background(), `
			INSERT INTO vehicle_measurements 
			(vehicle_id, measurement_type, measurement_time, measurement_value) 
			VALUES (?, ?, ?, ?)
		`, c.Param("id"), measurement.Type, time.Now().UTC(), measurement.Value)
		if err != nil {
			log.Errorf("Failed to insert measurement: %v", err)
			return c.String(http.StatusInternalServerError, "internal server error")
		} else {
			log.Infof(
				"Inserted new measurement: %s=%f for vehicle %s",
				measurement.Type,
				measurement.Value,
				c.Param("id"),
			)
		}

		return c.JSON(http.StatusOK, measurement)
	}
}
