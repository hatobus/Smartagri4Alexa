package outbound

import (
	"fmt"

	"github.com/ericdaugherty/alexa-skills-kit-golang"
	db "github.com/hatobus/UKEMOCHI/presenterdb"
)

func GetSmartAgriInfoFromMachineNO(slot map[string]alexa.IntentSlot) (string, error) {
	var speech string

	// 取得する機器が指定されなくてはいけないので確認する
	n := slot["machineNO"].Value
	if n == "" {
		return "取得したい機器のナンバーを一から三号機で指定してください", nil
	}

	farmInfoMachineNO, err := db.GetLatestDataFromDynamoDB(n)
	if err != nil {
		return "", err
	}

	k := slot["parameter"].Value
	if k == "" {
		speech = fmt.Sprintf(
			"%sからの情報は、温度は%f度、湿度は%fパーセント、水分量は%fパーセント、二酸化炭素濃度は%fppm、照度は%fルクスです、この情報は%Tに取得された情報です。",
			n, farmInfoMachineNO.Temperature, farmInfoMachineNO.Humidity, (farmInfoMachineNO.SoilHumidity/1024)*100,
			farmInfoMachineNO.Co2Concentration, farmInfoMachineNO.Illuminance, farmInfoMachineNO.Gettime,
		)

		return speech, nil
	}

	var resval string

	switch k {
	case "温度":
		resval = fmt.Sprintf("%s度", farmInfoMachineNO.Temperature)
	case "湿度":
		resval = fmt.Sprintf("%sパーセント", farmInfoMachineNO.Humidity)
	case "水分量":
		soilHumid := (farmInfoMachineNO.SoilHumidity / 1024) * 100
		resval = fmt.Sprintf("%fパーセント", soilHumid)
	case "二酸化炭素濃度":
		resval = fmt.Sprintf("%sppm", farmInfoMachineNO.Co2Concentration)
	case "照度":
		resval = fmt.Sprintf("%sルクス", farmInfoMachineNO.Illuminance)
	}

	speech = fmt.Sprintf("%sの%sは%sです。この情報は%sに取得された情報です。", n, k, resval, farmInfoMachineNO.Gettime)

	return speech, nil

	return "", nil
}
