# Проектная работа. Ротация баннеров.

## Запуск проект в докере

### Сборка проекта

```
$ make build
```

### Сборка и запуск проекта

```
$ make run
```

## Запуск проекта локально

### Запуск внешних сервисов

```
$ make run-external
$ make ps
$ make log
```

### Запуск гошного приложения

```
$ make run-local
```

### Запуск интеграционных тестов

```
$ make run-local
$ make integration-test
```

## Примеры запросов

### Добавление баннера в ротацию

```
POST localhost:8080/api/bannerRotatorService/v1/banner/add
Content-Type: application/json

{
  "slotID": 1,
  "bannerID": 1
}
```

Пример ответа

```
{
  "success": true,
  "errors": []
}
```

### Удаление баннера из ротации

```
POST localhost:8080/api/bannerRotatorService/v1/banner/remove
Content-Type: application/json

{
  "slotID": 1,
  "bannerID": 1
}
```

Пример ответа

```
{
  "success": true,
  "errors": []
}
```

### Переход по баннеру

```
POST localhost:8080/api/bannerRotatorService/v1/banner/hit
Content-Type: application/json

{
  "slotID": 1,
  "bannerID": 1,
  "groupID": 2
}
```

Пример ответа

```
{
  "success": true,
  "errors": []
}
```

### Получение избранного баннера для показа

```
POST localhost:8080/api/bannerRotatorService/v1/banner/get
Content-Type: application/json

{
  "slotID": 1,
  "groupID": 2
}
```

Пример ответа

```
{
  "success": true,
  "errors": [],
  "bannerID": "2"
}
```
