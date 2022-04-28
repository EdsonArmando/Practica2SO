# Pasos para ejecutar
Pasos para poder ejecutar el gRPC Node

## Cliente Node

El api de Node que leera los datos que son enviados desde locust y hacer la conexion al server de Go

### `node api_server.js`

## Server  Go

Es donde se ejecutaran los juegos .

### `go run gRPC-Server.go`

## Locust
Ejecuta el endpoitn en el cliente de node

### `locust -f generador.py`
