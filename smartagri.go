/*
  Copyright [2018] [Haga Fumito]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	alexa "github.com/ericdaugherty/alexa-skills-kit-golang"
	"github.com/joho/godotenv"
)

var err = godotenv.Load()
var URL = os.Getenv("APIURL")
var AppID = os.Getenv("APPID")

var ALX = &alexa.Alexa{ApplicationID: AppID, RequestHandler: &SmartAgri{}, IgnoreApplicationID: true, IgnoreTimestamp: true}

type SmartAgri struct{}

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

func (s *SmartAgri) OnSessionStarted(ctx context.Context, request *alexa.Request, session *alexa.Session, ctxPtr *alexa.Context, response *alexa.Response) error {

	log.Printf("OnSessionStarted requestId=%s, sessionId=%s", request.RequestID, session.SessionID)

	return nil
}

// Handle processes calls from Lambda
func Handle(ctx context.Context, requestEnv *alexa.RequestEnvelope) (interface{}, error) {
	return ALX.ProcessRequest(ctx, requestEnv)
}

// OnLaunch called with a reqeust is received of type LaunchRequest
func (s *SmartAgri) OnLaunch(ctx context.Context, request *alexa.Request, session *alexa.Session, ctxPtr *alexa.Context, response *alexa.Response) error {
	speechText := "これはスマートアグリの情報を取得できます。"
	speechText = speechText + "取得できる情報は 温度、湿度、二酸化炭素濃度、水分量に照度です。"
	speechText = speechText + "一から三号機までの情報にそれぞれアクセスできます。"

	log.Printf("OnLaunch requestId=%s, sessionId=%s", request.RequestID, session.SessionID)

	// response.SetSimpleCard(cardTitle, speechText)
	response.SetOutputText(speechText)
	response.SetRepromptText(speechText)

	response.ShouldSessionEnd = true

	return nil
}

// OnIntent called with a reqeust is received of type IntentRequest
func (s *SmartAgri) OnIntent(ctx context.Context, request *alexa.Request, session *alexa.Session, ctxPtr *alexa.Context, response *alexa.Response) error {

	log.Printf("OnIntent requestId=%s, sessionId=%s, intent=%s", request.RequestID, session.SessionID, request.Intent.Name)

	switch request.Intent.Name {
	case "getParamIntent":
		log.Println("getParamIntent triggered")
		speechText, err := Getsmartagriinfo(request.Intent.Slots)
		if err != nil {
			fmt.Println(err)
			speechText = "すみません、情報を取得できませんでした。"
		}

		// speechText := "コレが聞こえているということはAPIとかがおかしいよ"

		// response.SetSimpleCard(cardTitle, speechText)
		response.SetOutputText(speechText)

		log.Printf("Set Output speech, value now: %s", response.OutputSpeech.Text)
	case "AMAZON.HelpIntent":
		log.Println("何か助けが必要ですか")
		speechText := "何か助けが必要ですか"

		// response.SetSimpleCard("SmartAgri", speechText)
		response.SetOutputText(speechText)
		response.SetRepromptText(speechText)
	default:
		return errors.New("Invalid Intent")
	}

	return nil
}

// OnSessionEnded called with a reqeust is received of type SessionEndedRequest
func (s *SmartAgri) OnSessionEnded(ctx context.Context, request *alexa.Request, session *alexa.Session, ctxPtr *alexa.Context, response *alexa.Response) error {

	log.Printf("OnSessionEnded requestId=%s, sessionId=%s", request.RequestID, session.SessionID)

	return nil
}

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
		resval = farmInfoMachineNO.Temperature + "度"
	case "湿度":
		resval = farmInfoMachineNO.Humidity + "パーセント"
	case "水分量":
		soilHumid, _ := strconv.ParseFloat(farmInfoMachineNO.SoilHumidity, 32)
		soilHumid = (soilHumid / 1024) * 100
		resval = strconv.FormatFloat(soilHumid, 'f', 2, 64) + "パーセント"
	case "二酸化炭素濃度":
		resval = farmInfoMachineNO.Co2Concentration + "ppm"
	case "照度":
		resval = farmInfoMachineNO.Illuminance + "ルクス"
	}

	resval = resval + "です。この情報は" + farmInfoMachineNO.Time + "に取得された情報です。"

	speech = n + "の" + k + "は、" + resval
	return speech, nil
}

func GetFarmInfoMachineNO(machine string) (AgriData, error) {
	var machineNO string

	switch machine {
	case "一号機":
		machineNO = "1"
	case "二号機":
		machineNO = "2"
	case "三号機":
		machineNO = "3"
	}

	farmstruct, err := getFarmInfoFromAPI(machineNO)

	return farmstruct, err
}

// APIを叩いて情報を取得する
func getFarmInfoFromAPI(machineNO string) (AgriData, error) {
	var resdata []AgriData
	var tmpdata AgriData

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

func main() {
	lambda.Start(Handle)
}
