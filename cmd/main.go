package main

import (
	"crate_demo/internal"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	db := internal.NewDatabase()
	e.GET("/vehicles", internal.GetVehicles(db))
	e.GET("/vehicles/search", internal.GetVehiclesForSearch(db))

	e.POST("/vehicles/:id/measurements", internal.PostMeasurements(db))
	e.GET("/vehicles/:id/state_of_charge_chart", internal.GetStateOfChargeChart(db))

	e.Logger.Fatal(e.Start(":8080"))
}
