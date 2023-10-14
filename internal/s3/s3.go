package s3

import (
	"encoding/csv"
	"fmt"
	"io"
	"lambda-s3-ses-go/internal/ses"
	"lambda-s3-ses-go/types"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func GetFile() {
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
		if *item.Key == "txns.csv" {
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

func csvFile(r io.Reader) (*types.Process, error) {
	reader := csv.NewReader(r)

	_, err := reader.Read()
	if err != nil {
		return nil, err
	}

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	processResult := &types.Process{
		TotalBalance:        0.0,
		AverageDebitAmount:  0.0,
		AverageCreditAmount: 0.0,
		TotalByMonth:        make(map[string]int),
	}

	var numDebitTransactions int
	var numCreditTransactions int

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
			numCreditTransactions++
		} else {
			processResult.AverageDebitAmount += transaction
			numDebitTransactions++
		}

		processResult.TotalByMonth[month] += 1
		processResult.TotalBalance += transaction
	}

	if numCreditTransactions > 0 {
		processResult.AverageCreditAmount /= float64(numCreditTransactions)
	}
	if numDebitTransactions > 0 {
		processResult.AverageDebitAmount /= float64(numDebitTransactions)
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

	ses.SendSummaryEmail(processResult)
}
