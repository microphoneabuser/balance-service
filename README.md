# Микросервис для работы с балансом пользователей

## Список используемых технологий:

* [Gin](https://github.com/gin-gonic/gin) - HTTP Go Framework
* [Postgres](https://github.com/lib/pq) - СУБД PostgreSQL
* [Redis](https://github.com/go-redis/redis) - Реализация кэширования курса валют (ответа [freecurrencyapi.net](https://freecurrencyapi.net))
* [RabbitMQ](https://github.com/streadway/amqp) - RabbitMQ для реализации очереди SMS-уведомлений для последующей отправки, выполняемой другим микросервисом
* [Docker](https://www.docker.com/) - Docker
* [viper](https://github.com/spf13/viper) - Работа с файлами конфигурации
* [sqlx](https://github.com/jmoiron/sqlx) - Работа с БД
* [migrate](https://github.com/golang-migrate/migrate) - Миграции в БД

## Для запуска приложения:

``` bash
make build && make run
```

Если приложение запускается впервые, необходимо применить миграции к базе данных:

``` bash
make migrate
```

# Примеры запросов/ответов

[![Run in Postman](https://run.pstmn.io/button.svg)](https://www.postman.com/altimetry-cosmonaut-24747535/workspace/9b7c651f-7961-43ec-99a6-3af777ee7f1e/documentation/17406947-515440b7-8466-465d-bc3e-a49c8e40cc53)

### Получить баланс пользователя
``` bash
GET http://localhost:8080/balance
# Body(json)
{
    "id": 1
}
```

### Получить баланс пользователя в валюте, отличной от рубля (пример - USD)
``` bash
GET http://localhost:8080/balance?currency=USD
# Body(json)
{
    "id": 1
}
```

### Зачислить деньги на счет
``` bash
POST http://localhost:8080/accrual
# Body(json)
{
    "id": 1,
    "amount": 1000.0
}
```

### Списать деньги со счета
``` bash
POST http://localhost:8080/debiting
# Body(json)
{
    "id": 1,
    "amount": 200.0
}
```

### Перевести деньги с одного счета на другой
``` bash
POST http://localhost:8080/transfer
# Body(json)
{
    "sender_id": 1,
    "recipient_id": 2,
    "amount": 1000.0,
    "description": "Перевод на покупку чего-то"
}
```

### Просмотреть транзакции связанные с заданным счетом (с сортировкой по убыванию даты и пагинацией)
``` bash
GET http://localhost:8080/transactions?limit=10&offset=0&sort=timestamp:desc
# Body(json)
{
    "id": 2
}
```

### Просмотреть транзакции связанные с заданным счетом (с сортировкой по возрастанию суммы и пагинацией)
``` bash
GET http://localhost:8080/transactions?limit=10&offset=0&sort=amount:asc
# Body(json)
{
    "id": 2
}
```

## RabbitMQ

Чтобы просмотреть содержимое очереди для отправки SMS-уведомлений нужно перейти по ссылке: http://localhost:15672/#/queues/%2F/sms-queue
* Username: guest
* Password: guest

Данная очередь предназначена для уведомления пользователей о всех действиях произведенных с его счетом (зачисления, списания, переводы). В коде данного микросервиса формируется и публикуется в очередь полное сообщение с id пользователя, которому нужно отправить SMS. 
Пример сообщения:
``` bash
{
    "account_id":2,
    "message":"Счет-2 Зачисление 1000.00р Баланс: 1000.00р"
}
```
Задумка в том, что данные сообщения должен читать другой микросервис (consumer) и непосредственно производить отправку SMS-сообщения. 