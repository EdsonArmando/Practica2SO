# Manual Tecnico
En el siguiente manual se llevo a cabo la implementacion de un cluste de kubernetes en donde se expone el servicio por medio de un load balancer. A continuacion se muestran los pasos realizados

## Archivos de Configuracion Kubernetes

load-balancer-gRPC.yaml
Replicas = 1

Una parte importante del archivo .yaml es los containers del pod

seran 5
1. Kafka
	puerto: 9092
2. zookeeper
	puerto: 2181
3. clientnode
	puerto: 5000
4. producer
	puerto: 50051
5. subscriber
	puerto: 8080

### Comandos de Kubernetes

## Comando para el archivo generar el deployment

kubectl apply -f load-balancer-gRPC.yaml --namespace=practica2-201701029

## Comando para exponer el servicio

kubectl expose deployment practica2-201701029 --type=LoadBalancer --name=so-practica2 --load-balancer-ip=34.121.46.185

## Comando para obtener los servicios
kubectl get services

## Comando para crear nuevo namespace
kubectl config set-context --current --namespace=practica2-201701029

## Comando para cambiarse de namespace
kubectl config set-context --current --namespace=practica2-201701029

# Manual de Usuario
Al no poder contar con un frontend el usuario puede acceder al servicio utilizando la direccion ip del load balancer. En la siguiente imagen se ve como se puede mandar datos al client de node.

![postmand](https://user-images.githubusercontent.com/14056462/166612575-6ea1aa80-5f55-462b-aeb7-f011cbcc955b.png)
