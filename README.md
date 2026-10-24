# auth

gRPC-сервис аутентификации на Go.

## Структура

```text
cmd/
  auth/          Запуск сервера
  migrator/      Применение миграций
config/          YAML-конфигурация
internal/
  app/           Сборка приложения и запуск gRPC
  config/        Загрузка конфигурации
  domain/model/  Модель пользователя
  grpc/auth/     Обработчики RPC и валидация запросов
  service/auth/  Логика регистрации и входа
  storage/       Работа с SQLite
  lib/           JWT и логирование
proto/           Protobuf-контракт
gen/go/          Сгенерированные сообщения и gRPC-код
migration/       SQL-миграции Goose
test/            Интеграционные тесты
```

## Запуск

```powershell
$env:AUTH_JWT_SECRET = "your-secret-key"

# Подготовка базы данных
New-Item -ItemType Directory -Force storage | Out-Null
go run ./cmd/migrator -storagepath storage/auth.db -migrationpath migration

# Запуск gRPC-сервера
go run ./cmd/auth -config config/default_config.yaml
```

Порт, путь к базе и срок действия JWT задаются в [конфигурации](config/default_config.yaml). Секрет также можно указать в `.env` через `AUTH_JWT_SECRET`.

## API

| Метод | Запрос | Ответ |
| --- | --- | --- |
| `Register` | `username`, `email`, `password` | `user_id` нового пользователя |
| `Login` | `username`, `password` | JWT в поле `token` |
| `IsAdmin` | `user_id` | Признак администратора `is_admin` |

Контракт сервиса `auth.Auth`: [auth.proto](proto/auth/v1/auth.proto).

## Разработка

```sh
make build    # Сборка сервера
make migrate  # Применение миграций
make proto    # Генерация protobuf-кода
make test     # Интеграционные тесты
```

Без Make тесты можно запустить в отдельном терминале:

```sh
go test -v -count=1 ./test/...
```

Они проверяют регистрацию, дубликаты имени и email, вход и выдачу JWT. Созданные тестами пользователи остаются в базе.

## Конфигурация

```yaml
env: "production"
storage_path: "storage/auth.db"
token_ttl: 1h
grpc:
  port: 50051
  timeout: 5s
```

| Параметр | Описание |
| --- | --- |
| `env` | Название окружения (Сейчас не влияет на поведение сервиса). |
| `storage_path` | Путь к файлу SQLite относительно рабочей директории. |
| `token_ttl` | Срок действия JWT. |
| `grpc.port` | Порт gRPC-сервера. |
| `grpc.timeout` | Таймаут запросов. |
