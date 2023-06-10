# soa-grpc

## Сборка

### С докером
```
# Build client
docker build --tag mafia-client . -f ./build/client.dockerfile

# Build server
docker build --tag mafia-server . -f ./build/server.dockerfile
```

### Без докера

```
# Build client
go build ./cmd/client

# Build server
go build ./cmd/server
```

## Запуск

### С докером

```
# Run server
docker run -p 8113:8113 mafia-server

# Run client
docker run -it mafia-client -addr="server address"
```

Сервер в докере работает нормально, но с клиентом в докере не смог подключиться к серверу в докере/без докера

### Через docker-compose

```
# Run server
cd build
docker-compose up
```

### Без докера

```
# Run server
./server

# Run client
./client
```

## Опции

### Сервер

```
addr - адрес сервера
openvoting - рассылает уведомления о голосовании
ssize - количество игроков в одной игровой сессии
```
### Клиент

```
addr - адрес сервера
auto - автоматический выбор действий
lag - задержка между автоматическим выбором
```

## Подробности для проверки

Реализовал все пункты, но не смог в докер в клиенте
