package outbound

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ericdaugherty/alexa-skills-kit-golang"
	"github.com/hatobus/UKEMOCHI/model"
)

func Getsmartagriinfo(slot map[string]alexa.IntentSlot) (string, error) {
	var speech string

	// 取得する機器が指定されなくてはいけないので確認する
	n := slot["machineNO"].Value
	if n == "" {
		return "取得したい機器のナンバーを一から三号機で指定してください", nil
	}

	// 指定した機器の情報を構造体で取得してくる
	farmInfoMachineNO, err := GetFarmInfoMachineNO(n)
	if err != nil {
		return "", err
	}

	k := slot["parameter"].Value
	if k == "" {
		soilHumid, _ := strconv.ParseFloat(farmInfoMachineNO.SoilHumidity, 32)
		soilHumid = (soilHumid / 1024) * 100
		speech = n + "からの情報は、" +
			"温度は" + farmInfoMachineNO.Temperature + "度、" +
			"湿度は" + farmInfoMachineNO.Humidity + "パーセント、" +
			"水分量は" + strconv.FormatFloat(soilHumid, 'f', 2, 64) + "パーセント、" +
			"二酸化炭素濃度は" + farmInfoMachineNO.Co2Concentration + "ppm、" +
			"照度は" + farmInfoMachineNO.Illuminance + "ルクスです。" +
			"この情報は" + farmInfoMachineNO.Time + "に取得された情報です。"

		return speech, nil
	}

	var resval string

	switch k {
	case "温度":
		resval = fmt.Sprintf("%s度", farmInfoMachineNO.Temperature)
	case "湿度":
		resval = fmt.Sprintf("%sパーセント", farmInfoMachineNO.Humidity)
	case "水分量":
		soilHumid, _ := strconv.ParseFloat(farmInfoMachineNO.SoilHumidity, 32)
		soilHumid = (soilHumid / 1024) * 100
		resval = fmt.Sprintf("%fパーセント", soilHumid)
	case "二酸化炭素濃度":
		resval = fmt.Sprintf("%sppm", farmInfoMachineNO.Co2Concentration)
	case "照度":
		resval = fmt.Sprintf("%sルクス", farmInfoMachineNO.Illuminance)
	}

	speech = fmt.Sprintf("%sの%sは%sです。この情報は%sに取得された情報です。", n, k, resval, farmInfoMachineNO.Time)

	return speech, nil
}

func GetFarmInfoMachineNO(machine string) (model.AgriData, error) {
	var machineNO string

	switch machine {
	case "一号機":
		machineNO = "1"
	case "二号機":
		machineNO = "2"
	case "三号機":
		machineNO = "3"
	}

	// DynamoDBに置き換える
	farmstruct, err := getFarmInfoFromAPI(machineNO, os.Getenv("APIURL"))

	return farmstruct, err
}

// APIを叩いて情報を取得する
func getFarmInfoFromAPI(machineNO, APIURL string) (model.AgriData, error) {
	var resdata []model.AgriData
	var tmpdata model.AgriData

	res, err := http.Get(APIURL)
	if err != nil {
		log.Println("API呼び出せてません！！！！")
		return tmpdata, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	// 返ってきたjsonじゃないjsonを}でsplitする
	bodystring := strings.Split(string(body), "}")

	// } で splitしたので消えているから } をくっつけてjson unmarshall
	for _, jsondata := range bodystring {
		jsondata = jsondata + "}"
		json.Unmarshal([]byte(jsondata), &tmpdata)

		// 機械Noにあった番号でstructを追加
		if tmpdata.No == machineNO {
			resdata = append(resdata, tmpdata)
		}

	}

	log.Println(resdata[len(resdata)-1])
	return resdata[len(resdata)-1], nil
}
