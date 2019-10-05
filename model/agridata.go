package model

type AgriData struct {
	No               string `json:"no"`
	Date             string `json:"date"`
	Time             string `json:"time"`
	Temperature      string `json:"temperature"`
	Humidity         string `json:"humidity"`
	SoilHumidity     string `json:"soil_humidity"`
	Co2Concentration string `json:"co2_concentration"`
	Wavelength       string `json:"wavelength"`
	Illuminance      string `json:"illuminance"`
}
