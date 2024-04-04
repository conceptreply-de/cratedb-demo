package internal

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type Vehicle struct {
	Id          string `json:"id"`
	Vin         string `json:"vin"`
	Description string `json:"description"`
	Uses        string `json:"use"`
}

func GetVehicles(db *Database) func(c echo.Context) error {
	return func(c echo.Context) error {

		rows, err := db.conn.Query(context.Background(), "select * from vehicles;")
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			panic(err)
		}
		defer rows.Close()
		var vehicles []Vehicle = []Vehicle{}
		for rows.Next() {
			var id, vin, description, use string
			err = rows.Scan(&id, &vin, &description, &use)
			vehicles = append(vehicles, Vehicle{Id: id, Vin: vin, Description: description, Uses: use})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Scan failed: %v\n", err)
				panic(err)
			}
		}
		fmt.Fprintf(os.Stdout, "Vehicles: %s\n", vehicles)
		return c.JSON(http.StatusOK, vehicles)
	}
}

func GetVehiclesForSearch(db *Database) func(c echo.Context) error {
	return func(c echo.Context) error {
		var search_term = c.QueryParam("q")
		var sql_command = `
			SELECT id, vin, description, use, _score 
			FROM vehicles 
			WHERE match((text_fields, vin_grams), $1) 
			USING best_fields WITH (fuzziness=1) 
			ORDER BY _score desc;
		`

		fmt.Fprintf(os.Stdout, "Query: %s SQL Command: %s\n", search_term, sql_command)
		rows, err := db.conn.Query(context.Background(), sql_command, search_term)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			panic(err)
		}
		defer rows.Close()
		var vehicles []Vehicle = []Vehicle{}
		for rows.Next() {
			var id, vin, description, use, _score string
			err = rows.Scan(&id, &vin, &description, &use, &_score)
			vehicles = append(vehicles, Vehicle{Id: id, Vin: vin, Description: description, Uses: use})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Scan failed: %v\n", err)
				panic(err)
			}
		}
		fmt.Fprintf(os.Stdout, "Vehicles: %s\n", vehicles)
		return c.JSON(http.StatusOK, vehicles)
	}
}
