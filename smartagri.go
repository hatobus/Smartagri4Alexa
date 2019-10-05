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
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	alexa "github.com/ericdaugherty/alexa-skills-kit-golang"
	"github.com/hatobus/UKEMOCHI/agrialexa"
	"github.com/joho/godotenv"
)

var err = godotenv.Load()
var URL = os.Getenv("APIURL")
var AppID = os.Getenv("APPID")

var ALX = &alexa.Alexa{ApplicationID: AppID, RequestHandler: &agrialexa.SmartAgri{}, IgnoreApplicationID: true, IgnoreTimestamp: true}

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

// Handle processes calls from Lambda
func Handle(ctx context.Context, requestEnv *alexa.RequestEnvelope) (interface{}, error) {
	return ALX.ProcessRequest(ctx, requestEnv)
}

func main() {
	lambda.Start(Handle)
}
