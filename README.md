# soa-grpc

## Сборка
```
# Build client
docker build --tag mafia-client . -f ./build/client.dockerfile

# Build server
docker build --tag mafia-server . -f ./build/server.dockerfile
```

## Запуск

### Запуск через докер

```
# Run server
docker run -p 8080:8080 mafia-server -openvoting

# Run client
docker run -it mafia-client -addr="server address"
```

Сервер в докере работает нормально, но с клиентом в докере не смог подключиться к серверу в докере/без докера

### Запуск без докера

```
# Run server
go run ./cmd/server

# Run client
go run ./cmd/client
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
