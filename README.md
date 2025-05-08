# Сервис подписок с использованием пакета `subpub`.


## 1. Пакет `subpub`
[pkg/subpub](https://github.com/andreyxaxa/PubSub_gRPC_Service/tree/main/pkg/subpub)

Шина событий, работающая по принципу Publisher-Subscribe.
- На один subject может подписываться (и отписываться) множество подписчиков - Мапа подписчиков на каждый `subject`
- Один медленный подписчик не тормозит остальных - У каждого подписчика своя горутина и буферизированный канал.
- Сохраняется порядок сообщений (FIFO очередь) - Канал(chan) FIFO.
- Метод `Close` учитывает переданный контекст.
- Горутины не текут - `WaitGroup`, `select`, `closeCh`.

`subpub` наружу предоставляет интерфейсы для внедрения (Dependency Injection) и конструктор `NewSubPub()`.

Протестировать пакет:
- `make test`

## 2. gRPC сервис подписок.
Клиенты могут:
- Подписаться на сообщение по ключу - `PubSub/Subscribe`
```
{
	"key": "WEATHER"
}
```
- Публиковать сообщения по ключу - `PubSub/Publish`
```
{
	"data": "Moscow - 08.05.25 - +10C",
	"key": "WEATHER"
}
```
- Подписчик определенной темы получает опубликованные сообщения через grpc-стрим.

## Детали 

- У сервиса есть конфиг - [config/config.go](https://github.com/andreyxaxa/PubSub_gRPC_Service/blob/main/config/config.go); Читается из `.env` файла. В рамках тестового задания .env прямо в репозитории, очевидно в проде он должен быть заигнорен.
- Есть логгер - [pkg/logger](https://github.com/andreyxaxa/PubSub_gRPC_Service/tree/main/pkg/logger); Интерфейс позволяет подменить логгер.
- В слое хэндлеров применяется версионирование - [internal/controller/grpc/v1](https://github.com/andreyxaxa/PubSub_gRPC_Service/tree/main/internal/controller/grpc/v1).
  Для версии v2 нужно будет просто добавить папку `grpc/v2` с таким же содержимым, в файле [internal/controller/grpc/router.go](https://github.com/andreyxaxa/PubSub_gRPC_Service/blob/main/internal/controller/grpc/router.go) добавить строку:
```go
{
    v1.NewPubSubRouter(app, sp, l)
}

{
    v2.NewPubSubRouter(app, sp, l)
}
```
- Используется dependency injection - [internal/controller/grpc/v1/controller.go](https://github.com/andreyxaxa/PubSub_gRPC_Service/blob/main/internal/controller/grpc/v1/controller.go).
- Реализован graceful shutdown - [internal/app/app.go](https://github.com/andreyxaxa/PubSub_gRPC_Service/blob/main/internal/app/app.go).
- Удобная и гибкая конфигурация gRPC сервера - [pkg/grpcserver/options.go](https://github.com/andreyxaxa/PubSub_gRPC_Service/blob/main/pkg/grpcserver/options.go).
  Позволяет конфигурировать сервер в конструкторе таким образом:
```go
grpcServer := grpcserver.New(grpcserver.Port(cfg.GRPC.Port))
```

## Запуск

### Local:
Устанавливает зависимости, генерирует исходники из `.proto` и запускает приложение. 
```
make run
```

### Docker:
```
make compose-up
```

## Прочие `make` команды
- `make deps`:
```
go mod tidy && go mod verify
```
- `make proto-v1`:
```
protoc --go_out=. \
--go_opt=paths=source_relative \
--go-grpc_out=. \
--go-grpc_opt=paths=source_relative \
docs/proto/pubsub/v1/*.proto
```
- `make compose-down`:
```
docker compose -f docker-compose.yml down
```
