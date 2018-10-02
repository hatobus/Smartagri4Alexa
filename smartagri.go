package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

type AgreData struct {
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

func getHouseinfoJSON(APIURL string) error {

	res, err := http.Get(APIURL)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	// 返ってきたjsonじゃないjsonを}でsplitする
	bodystring := strings.Split(string(body), "}")

	// spew.Dump(bodystring)

	var resdataNO1, resdataNO2, resdataNO3 []AgreData
	var tmpdata AgreData

	// } で splitしたので消えているから } をくっつけてjson unmarshall
	for _, jsondata := range bodystring {
		jsondata = jsondata + "}"
		json.Unmarshal([]byte(jsondata), &tmpdata)

		// 機械Noで分ける
		switch tmpdata.No {
		case "1":
			resdataNO1 = append(resdataNO1, tmpdata)
		case "2":
			resdataNO2 = append(resdataNO2, tmpdata)
		case "3":
			resdataNO3 = append(resdataNO3, tmpdata)
		default:
			log.Fatal("Invalid value")
		}

	}

	spew.Dump(resdataNO1)
	spew.Dump(resdataNO2)
	spew.Dump(resdataNO3)

	return err
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	url := os.Getenv("APIURL")

	getHouseinfoJSON(url)

}
