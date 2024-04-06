package internal

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type VehicleCurrentLocation struct {
	Latitude   float64   `json:"lat"`
	Longitude  float64   `json:"lon"`
	VehicleId  string    `json:"vehicle_id"`
	ReceivedAt time.Time `json:"received_at"`
}

func SetVehicleLocation(db *Database) func(c echo.Context) error {
	return func(c echo.Context) error {

		var vehicleLocation VehicleCurrentLocation
		err := c.Bind(&vehicleLocation)
		if err != nil {
			log.Errorf("Failed to bind vehicle location: %v", err)
			return c.String(http.StatusBadRequest, "bad request")
		}

		vehicleLocation.ReceivedAt = time.Now().UTC()
		vehicleLocation.VehicleId = c.Param("id")

		_, err = db.conn.Exec(context.Background(), `
			INSERT INTO vehicle_current_locations 
			(vehicle_id, vehicle_location, received_at) 
			VALUES ($1, [$2::decimal, $3::decimal], $4)
			ON CONFLICT (vehicle_id) 
			DO UPDATE SET vehicle_location = [ $2::decimal, $3::decimal ], received_at = $4
		`,
			vehicleLocation.VehicleId,
			fmt.Sprintf("%f", vehicleLocation.Longitude),
			fmt.Sprintf("%f", vehicleLocation.Latitude),
			vehicleLocation.ReceivedAt,
		)

		if err != nil {
			log.Errorf("Failed to insert measurement: %v", err)
			return c.String(http.StatusInternalServerError, "internal server error")
		} else {
			log.Infof(
				"Upserted new vehicle location: %f=%f for vehicle %s",
				vehicleLocation.Longitude,
				vehicleLocation.Latitude,
				c.Param("id"),
			)
		}

		var countryName string

		err = db.conn.QueryRow(context.Background(), `
			SELECT name FROM countries 
			WHERE match (shape, concat('POINT(', $1  , ' '  , $2 ,  ' )') )
			LIMIT 1;`,
			fmt.Sprintf("%f", vehicleLocation.Longitude),
			fmt.Sprintf("%f", vehicleLocation.Latitude),
		).Scan(&countryName)

		if err != nil {
			if err != pgx.ErrNoRows {
				log.Errorf("Failed to find country: %v", err)
				return c.String(http.StatusInternalServerError, "internal server error")
			} else {
				log.Warn("No country found")

			}
		} else {
			log.Infof("Vehicle is in %s", countryName)
		}

		return c.JSON(http.StatusOK, vehicleLocation)
	}
}

// 16 48 - Vienna
// 11 47 - Munich
// 8.5 47 - Zurich (sort of)
