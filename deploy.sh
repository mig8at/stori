#!/bin/bash

export GOOS=linux

go build -o main
zip deployment.zip main
aws lambda update-function-configuration --function-name summary --environment "Variables={EMAIL_TO=mig8at@gmail.com,EMAIL_FROM=mig8at@gmail.com}"
aws lambda update-function-code --function-name summary --zip-file fileb://deployment.zip
aws lambda invoke --function-name summary --payload '{}' output.txt
