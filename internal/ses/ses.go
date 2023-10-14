package ses

import (
	"fmt"
	"lambda-s3-ses-go/internal/dynamo"
	"lambda-s3-ses-go/types"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

func SendSummaryEmail(processResult *types.Process) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))

	svc := ses.New(sess)
	emailTo := os.Getenv("EMAIL_TO")
	emailFrom := os.Getenv("EMAIL_FROM")

	var htmlBody strings.Builder
	htmlBody.WriteString("<h1>Transaction summary</h1>")
	htmlBody.WriteString(fmt.Sprintf("<p><strong>Total balance:</strong> %.2f</p>", processResult.TotalBalance))
	htmlBody.WriteString(fmt.Sprintf("<p><strong>Average credit amount:</strong> %.2f</p>", processResult.AverageCreditAmount))
	htmlBody.WriteString(fmt.Sprintf("<p><strong>Average debit amount:</strong> %.2f</p>", processResult.AverageDebitAmount))

	htmlBody.WriteString("<p><strong>Number of transactions:</strong></p><ul>")
	for month, totalCredit := range processResult.TotalByMonth {
		htmlBody.WriteString(fmt.Sprintf("<li>In %s: %d</li>", month, totalCredit))
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
				Data: aws.String("Transaction summary"),
			},
		},
		Source: aws.String(emailFrom),
	}

	_, err := svc.SendEmail(input)
	if err != nil {
		fmt.Println("Error al enviar correo:", err)
	}
	dynamo.SaveToDynamoDB(processResult)
}
