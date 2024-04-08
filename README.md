# effective-mobile-test
Тестовое задание от Effective Mobile

## Используемые технологии
- Swagger
- Redis
- JaegerTracing
- Mockery & testify

## О проекте
Проект был выполнен в чистой архитектуре. Имеются тесты и валидация.
Внешнее апи было замокано.

## Инструкция к запуску
Для запуска проекта требуется docker-compose.

`docker-compose up -d`

Далее запустите миграции.

`make migrate-up`

## API
http://localhost:3010/api/v1
## Jaeger
http://localhost:16686/search
## Swagger
http://localhost:3010/swagger/index.html

