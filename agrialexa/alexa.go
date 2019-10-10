package agrialexa

import (
	"context"
	"errors"
	"fmt"

	"github.com/ericdaugherty/alexa-skills-kit-golang"
	"github.com/hatobus/UKEMOCHI/logging"
	"github.com/hatobus/UKEMOCHI/outbound"
	"go.uber.org/zap"
)

type SmartAgri struct{}

func (s *SmartAgri) OnSessionStarted(ctx context.Context, request *alexa.Request, session *alexa.Session, ctxPtr *alexa.Context, response *alexa.Response) error {

	logging.Log().Info("OnSessinon Started", zap.Strings("requestId=%s, sessionId=%s", []string{request.RequestID, session.SessionID}))
	return nil
}

// OnLaunch called with a reqeust is received of type LaunchRequest
func (s *SmartAgri) OnLaunch(ctx context.Context, request *alexa.Request, session *alexa.Session, ctxPtr *alexa.Context, response *alexa.Response) error {
	speechText := "これはスマートアグリの情報を取得できます。取得できる情報は 温度、湿度、二酸化炭素濃度、水分量に照度です。一から三号機までの情報にそれぞれアクセスできます。"

	logging.Log().Info("OnLaunch started", zap.Strings("OnLaunch requestId=%s, sessionId=%s", []string{request.RequestID, session.SessionID}))

	// response.SetSimpleCard(cardTitle, speechText)
	response.SetOutputText(speechText)
	response.SetRepromptText(speechText)

	response.ShouldSessionEnd = true

	return nil
}

// OnIntent called with a reqeust is received of type IntentRequest
func (s *SmartAgri) OnIntent(ctx context.Context, request *alexa.Request, session *alexa.Session, ctxPtr *alexa.Context, response *alexa.Response) error {

	logging.Log().Info(
		"OnIntent started",
		zap.Strings(
			"OnIntent requestId=%s, sessionId=%s, intent=%s",
			[]string{request.RequestID, session.SessionID, request.Intent.Name},
		),
	)

	switch request.Intent.Name {
	case "getParamIntent":
		speechText, err := outbound.Getsmartagriinfo(request.Intent.Slots)
		if err != nil {
			fmt.Println(err)
			speechText = "すみません、情報を取得できませんでした。"
		}

		response.SetOutputText(speechText)

		logging.Log().Info("Set Output speech", zap.Strings("now: %s", []string{response.OutputSpeech.Text}))
	case "AMAZON.HelpIntent":
		speechText := "何か助けが必要ですか"

		response.SetOutputText(speechText)
		response.SetRepromptText(speechText)
	default:
		return errors.New("Invalid Intent")
	}

	return nil
}

// OnSessionEnded called with a reqeust is received of type SessionEndedRequest
func (s *SmartAgri) OnSessionEnded(ctx context.Context, request *alexa.Request, session *alexa.Session, ctxPtr *alexa.Context, response *alexa.Response) error {

	logging.Log().Info("OnSessinonEnded", zap.Strings("OnSessionEnded requestId=%s, sessionId=%s", []string{request.RequestID, session.SessionID}))

	return nil
}
