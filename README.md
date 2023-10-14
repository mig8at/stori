# AWS Lambda Function - Summary
Esta función Lambda de AWS, denominada `summary`

1. **Lectura de archivos CSV desde S3**: La función inicialmente lee un archivo CSV almacenado en un bucket de Amazon S3. Este archivo contiene una lista de transacciones financieras que necesitan ser procesadas.

2. **Procesamiento de Datos**: Una vez obtenido el archivo, la función lo procesa para calcular varios parámetros estadísticos como el balance total, la cantidad media de débito y crédito, y la distribución de transacciones por mes.

3. **Envío de correo electrónico**: Tras el procesamiento de los datos, la función utiliza Amazon Simple Email Service (SES) para enviar un resumen de las estadísticas calculadas a una dirección de correo electrónico específica.

4. **Almacenamiento en DynamoDB**: Finalmente, todos los datos procesados y calculados se almacenan en una tabla de Amazon DynamoDB para futuras consultas y análisis.


## Requisitos previos

- AWS CLI instalado y configurado.
- Go instalado (versión 1.x).
- Git instalado.
- Cuenta en AWS.

## Pasos

### Crear una función Lambda en AWS

1. Ve al panel de AWS Lambda y haz clic en "Crear función".
2. Elige "Crear desde cero" e introduce los detalles necesarios.
   - Nombre de la función: `summary`
   - Lenguaje de ejecución: `Go 1.x`
3. Haz clic en "Crear función".

### Configurar permisos IAM

1. Ve al panel de IAM en AWS.
2. Crea un nuevo rol con permisos para ejecutar Lambda y acceder a otros recursos necesarios (como DynamoDB, S3, SES, etc.).

**Ejemplo de política IAM:**

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
```


## Clonar el repositorio

Para obtener una copia del código en tu máquina local, sigue estos pasos:

```bash
# Clonar el repositorio
git clone https://github.com/mig8at/stori.git

# Cambiar al directorio del proyecto
cd stori

# Dar permisos de ejecución al script
chmod +x deploy.sh

# Ejecutar el script de despliegue
./deploy.sh