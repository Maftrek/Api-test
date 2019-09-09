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