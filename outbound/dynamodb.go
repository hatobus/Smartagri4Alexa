package outbound

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ericdaugherty/alexa-skills-kit-golang"
	"github.com/hatobus/UKEMOCHI/model"
)

func GetSmartAgriInfoFromMachineNO(slot map[string]alexa.IntentSlot) (string, error) {
	var speech string

	// 取得する機器が指定されなくてはいけないので確認する
	n := slot["machineNO"].Value
	if n == "" {
		return "取得したい機器のナンバーを一から三号機で指定してください", nil
	}

	ni, _ := strconv.Atoi(n)
	d := float64(ni / 10)
	duration := -1 * ni

	farmInfoMachineNO := &model.IoTTable{
		Gettime:          time.Now().Add(time.Duration(duration) * time.Minute),
		Temperature:      10.3 + d,
		Humidity:         30.8 + d,
		SoilHumidity:     48.2 + d,
		Co2Concentration: 28.0 + d,
		Illuminance:      float64(250 + ni),
	}
	// farmInfoMachineNO, err := db.GetLatestDataFromDynamoDB(n)
	// if err != nil {
	// 	return "", err
	// }

	k := slot["parameter"].Value
	if k == "" {
		speech = fmt.Sprintf(
			"%sからの情報は、温度は%.2f度、湿度は%.2fパーセント、水分量は%.2fパーセント、二酸化炭素濃度は%.2fppm、照度は%.2fルクスです、この情報は10分前に取得された情報です。",
			n, farmInfoMachineNO.Temperature, farmInfoMachineNO.Humidity, (farmInfoMachineNO.SoilHumidity/1024)*100,
			farmInfoMachineNO.Co2Concentration, farmInfoMachineNO.Illuminance,
		)

		return speech, nil
	}

	var resval string

	switch k {
	case "温度":
		resval = fmt.Sprintf("%.2g度", farmInfoMachineNO.Temperature)
	case "湿度":
		resval = fmt.Sprintf("%.2gパーセント", farmInfoMachineNO.Humidity)
	case "水分量":
		soilHumid := (farmInfoMachineNO.SoilHumidity / 1024) * 100
		resval = fmt.Sprintf("%.2gパーセント", soilHumid)
	case "二酸化炭素濃度":
		resval = fmt.Sprintf("%.2gppm", farmInfoMachineNO.Co2Concentration)
	case "照度":
		resval = fmt.Sprintf("%.2gルクス", farmInfoMachineNO.Illuminance)
	}

	speech = fmt.Sprintf("%sの%sは%sです。この情報は10分前に取得された情報です。", n, k, resval)

	return speech, nil
}
