# Тестовое в Пилотные проекты VK

## Запуск приложения

```commandline
docker-compose up
```

## Описание ручек

### Логин

#### `POST /api/login`

Авторизация пользователя по юзернейму и паролю, выдает токен авторизации, сформированный с помощью JWT и приватного ключа RSA

`Запрос`

```json
Host: 0.0.0.0:9241/api/login
Content-Type: application/json

{
    "username": "admin",
    "password": "presale"
}
```

`Ответ`

```json
status = 200
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiU29tZSBSaWNrcm9sbCJ9.TnMTNkbyt12KHJ55fQFX1Cz-SA5V4UqJkZop5Ufp2SQ"
}

status = 500
```

### Запись

#### `POST /api/write`

Записывает данные по ключам в базу данных

`Запрос`

```json
Host: 0.0.0.0:9241/api/login
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
    "data": {
        "key1": "value1",
        "key2": "value2",
        "key3": "value3"
    }
}
```

`Ответ`

```json
status = 200
{
    "status": "success"
}

status = 500
```

### Чтение

#### `POST /api/read`

Читает данные по ключам в базу данных

`Запрос`

```json
Host: 0.0.0.0:9241/api/login
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
    "keys": ["key1", "key2", "key3"]
}
```

`Ответ`

```json
status = 200
{
    "data": {
        "key1": "value1",
        "key2": "value2",
        "key3": 1
    }
}

status = 500
```
