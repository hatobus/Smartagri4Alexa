package model

import "time"

type IoTTable struct {
	No               string    `dynamo:"no"`
	MACAddr          string    `dynamo:"mac_addr"`
	Gettime          time.Time `dynamo:"gettime"`
	Temperature      float64   `dynamo:"temperature"`
	Humidity         float64   `dynamo:"humidity"`
	SoilHumidity     float64   `dynamo:"soil_humidity"`
	Co2Concentration float64   `dynamo:"co2_concentration"`
	Wavelength       float64   `dynamo:"wavelength"`
	Illuminance      float64   `dynamo:"illuminance"`
}
