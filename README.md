# Stori Transaction Summary Processor

Este proyecto procesa un archivo CSV que contiene una lista de transacciones de débito y crédito. Luego, envía un resumen de estas transacciones por correo electrónico.

## Requisitos

- Go 1.15 o superior
- AWS CLI configurado
- Cuenta de AWS con acceso a S3, Lambda y SES

## Configuración

1. Configura tus variables de entorno para el correo electrónico de destino y de origen en tu función AWS Lambda.

    ```shell
    EMAIL_TO=your.email@example.com
    EMAIL_FROM=sender.email@example.com
    ```

2. Crea un bucket S3 donde se guardarán los archivos CSV de transacciones.

3. Configura una función Lambda para ejecutar el código. Asegúrate de otorgar los permisos necesarios para acceder a S3 y SES.

## Instalación

1. Clona este repositorio.

    ```shell
    git clone https://github.com/mig8at/stori.git    
    ```

2. Cambia al directorio del proyecto.

    ```shell
    cd stori-transaction-summary
    ```

3. Ejecuta el siguiente comando para construir el proyecto.

    ```shell
    go build .
    ```

4. Sube el ejecutable a tu función Lambda en AWS.

## Uso

1. Sube tu archivo CSV de transacciones al bucket de S3.

2. Invoca la función Lambda manualmente o configura un disparador, como un evento de S3 cuando se sube un nuevo archivo.

3. Verifica tu correo electrónico para el resumen de transacciones.

## Configuración

### Permisos de IAM

Es necesario configurar un rol de IAM con los siguientes permisos para que la función Lambda pueda interactuar con S3, SES y DynamoDB:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "s3:GetObject",
                "s3:ListBucket"
            ],
            "Resource": [
                "arn:aws:s3:::stori-sumary/*",
                "arn:aws:s3:::stori-sumary"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "ses:SendEmail",
                "ses:SendRawEmail"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "dynamodb:PutItem",
                "dynamodb:GetItem",
                "dynamodb:Query",
                "dynamodb:Scan",
                "dynamodb:UpdateItem",
                "dynamodb:DeleteItem"
            ],
            "Resource": "*"
        }
    ]
}