package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ses"
)

type Process struct {
	TotalBalance        float64
	AverageDebitAmount  float64
	AverageCreditAmount float64
	TotalByMonth        map[string]int
}

func saveToDynamoDB(prs *Process) error {
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

func ListS3Files(fileName string) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))

	svc := s3.New(sess)

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String("stori-sumary"),
	}

	result, err := svc.ListObjectsV2(input)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, item := range result.Contents {
		if *item.Key == fileName {
			fmt.Println("Encontrado txns.csv, procediendo a procesar...")
			ProcessCSVFile("stori-sumary", *item.Key)
			break
		}
	}
}

func extractMonthFromDate(date string) (time.Time, error) {
	dateParts := strings.Split(date, "/")
	if len(dateParts) != 2 {
		return time.Time{}, fmt.Errorf("formato de fecha no válido: %s", date)
	}

	currentYear := time.Now().Year()
	dateWithYear := fmt.Sprintf("%s/%s/%d", dateParts[0], dateParts[1], currentYear)
	dateFormat := "1/2/2006"

	parsedDate, err := time.Parse(dateFormat, dateWithYear)
	if err != nil {
		return time.Time{}, err
	}

	return parsedDate, nil
}

func csvFile(r io.Reader) (*Process, error) {
	reader := csv.NewReader(r)

	_, err := reader.Read()
	if err != nil {
		return nil, err
	}

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	processResult := &Process{
		TotalBalance:        0.0,
		AverageDebitAmount:  0.0,
		AverageCreditAmount: 0.0,
		TotalByMonth:        make(map[string]int),
	}

	for _, row := range rows {

		if len(row) != 3 {
			continue
		}

		date := row[1]
		transactionStr := row[2]

		transaction, err := strconv.ParseFloat(transactionStr, 64)
		if err != nil {
			fmt.Println("Error al convertir transacción:", err)
			continue
		}

		monthTime, err := extractMonthFromDate(date)
		if err != nil {
			continue
		}

		month := monthTime.Month().String()

		if transaction > 0 {
			processResult.AverageCreditAmount += transaction
		} else {
			processResult.AverageDebitAmount += transaction
		}

		processResult.TotalByMonth[month] += 1
		processResult.TotalBalance += transaction
	}

	return processResult, nil
}

func ProcessCSVFile(bucket string, itemKey string) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))

	svc := s3.New(sess)
	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(itemKey),
	})

	if err != nil {
		fmt.Println("Error al descargar el archivo:", err)
		return
	}

	processResult, err := csvFile(result.Body)
	if err != nil {
		fmt.Println("Error al procesar el archivo:", err)
		return
	}

	SendSummaryEmail(processResult)
}

func SendSummaryEmail(processResult *Process) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))

	svc := ses.New(sess)
	emailTo := os.Getenv("EMAIL_TO")
	emailFrom := os.Getenv("EMAIL_FROM")

	var htmlBody strings.Builder
	htmlBody.WriteString("<h1>Resumen de transacciones</h1>")
	htmlBody.WriteString(fmt.Sprintf("<p><strong>Balance total:</strong> %.2f</p>", processResult.TotalBalance))

	htmlBody.WriteString("<p><strong>Average Credit:</strong></p><ul>")
	for month, totalCredit := range processResult.TotalByMonth {
		htmlBody.WriteString(fmt.Sprintf("<li>Number of transactions in %s: %d</li>", month, totalCredit))
	}
	htmlBody.WriteString("</ul>")

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(emailTo),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Data: aws.String(htmlBody.String()),
				},
			},
			Subject: &ses.Content{
				Data: aws.String("Resumen de Transacciones"),
			},
		},
		Source: aws.String(emailFrom),
	}

	_, err := svc.SendEmail(input)
	if err != nil {
		fmt.Println("Error al enviar correo:", err)
	}
	saveToDynamoDB(processResult)
}

func handler() {
	ListS3Files("txns.csv")
}

func main() {
	lambda.Start(handler)
}
