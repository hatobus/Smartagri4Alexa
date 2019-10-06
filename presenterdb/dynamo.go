package presenterdb

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/hatobus/UKEMOCHI/model"
)

var c = credentials.NewStaticCredentials(os.Getenv(("AWS_ACCESS_KEY")), os.Getenv("AWS_SECRET_ACCESS_KEY"), "") // 最後の引数は[セッショントークン]今回はなしで

var db = dynamo.New(session.New(), &aws.Config{
	Credentials: c,
	// 特に理由がない限りlambdaと同じリージョンでDynamoDBを動かす
	Region: aws.String(os.Getenv("AWS_REGION")),
})
var table = db.Table(os.Getenv("DYNAMO_TABLE"))

func GetLatestDataFromDynamoDB(machineNO string) (*model.IoTTable, error) {
	record := &model.IoTTable{}

	if err := table.Get("no", machineNO).All(record); err != nil {
		return nil, err
	}

	return record, nil
}
