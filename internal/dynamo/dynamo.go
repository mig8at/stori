package dynamo

import (
	"lambda-s3-ses-go/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func SaveToDynamoDB(prs *types.Process) error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))

	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(prs)
	if err != nil {
		return err
	}

	av["YEAR"] = &dynamodb.AttributeValue{S: aws.String("2023")}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Summary"),
	}

	_, err = svc.PutItem(input)

	return err
}
