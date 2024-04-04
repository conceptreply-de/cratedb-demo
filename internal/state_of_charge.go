package internal

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type StateOfChargeAggValue struct {
	PeriodFrom       time.Time `json:"period_from"`
	PeriodTo         time.Time `json:"period___to"`
	VehicleId        string    `json:"vehicle_id"`
	AvgStateOfCharge float64   `json:"avg_state_of_charge"`
}

type StateOfChargeChart struct {
	Values []StateOfChargeAggValue `json:"values"`
}

func GetStateOfChargeChart(db *Database) func(c echo.Context) error {
	return func(c echo.Context) error {

		rows, err := db.conn.Query(c.Request().Context(), `
			WITH avg_measurements AS (
				SELECT vehicle_id,
				DATE_BIN('1 minutes'::INTERVAL, measurement_time, 0) AS period,
				AVG(measurement_value) AS avg_state_of_charge
				FROM vehicle_measurements
				WHERE measurement_type = 'STATE_OF_CHARGE'
				AND measurement_partition > NOW() - INTERVAL '20 minutes'
				GROUP BY 1, 2 
				ORDER BY 1, 2
			)
			SELECT period period_from,
			       period + INTERVAL '1 minutes' period_to,
				am.vehicle_id vehicle_id,
				avg_state_of_charge  
			FROM avg_measurements am, vehicles v
			WHERE am.vehicle_id = v.id
			AND v.id = $1
			LIMIT 100;
		`, c.Param("id"))

		if err != nil {
			log.Errorf("QueryRow failed: %v\n", err)
			return c.String(http.StatusInternalServerError, "internal server error")
		}

		defer rows.Close()
		var values []StateOfChargeAggValue = make([]StateOfChargeAggValue, 0)
		for rows.Next() {
			var periodFrom, periodTo time.Time
			var vehicleId string
			var avgStateOfCharge float64
			err = rows.Scan(&periodFrom, &periodTo, &vehicleId, &avgStateOfCharge)
			values = append(values, StateOfChargeAggValue{
				PeriodFrom:       periodFrom,
				PeriodTo:         periodTo,
				VehicleId:        vehicleId,
				AvgStateOfCharge: avgStateOfCharge,
			})
			if err != nil {
				log.Errorf("Scan failed: %v\n", err)
				return c.String(http.StatusInternalServerError, "internal server error")
			}
		}

		return c.JSON(http.StatusOK, StateOfChargeChart{
			Values: values,
		})
	}
}
