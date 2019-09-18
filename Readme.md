Для работы на компьютере должен быть установлен nats-streaming-server.
Nats server необходимо запустить командой nats-streaming-server.

Порт приложения - :8110
GET /api/news/1 -  получить новость по id (1- id)
POST /api/news/ - создать новость 
body (json)
```json
{
	"title": "Name",
	"date": "2019-01-01"
}
```

Есть возможность развернуть докер образ, но в этом случае потребуется иметь nats-streaming-server и postgres не локально. Если условие выполняется то адреса nats-streaming-server и postgres необходимо прописать в config/config.toml

[NatsServer]
  Address = "АДРЕСС NATS"

Для запуска докер образа необходимо выполнить в директории программы:

1) CGO_ENABLED=0 GOOS=linux go build -o main .

2) docker build -t api_test:0.0.1 .

3) docker run api_test:0.0.1
