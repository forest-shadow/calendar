# Проект Calendar

## Пример тестовой API

|Method|URL|Descrition|
|:--|:--|:--|
|POST|/api/v1/signup| Создать пользователя и получить для него токен|
|POST|/api/v1/signin| Получить токен для существующего пользователя |
|POST|/api/v1/verify| Получить информацию о пользователе по токену |
|PATCH|/api/v1/users/{id}| Обновить данные пользователя |

## Структура проекта

```sh
├── cmd -- точки входа в приложение
├── internal
│   ├── apperrs -- основные типы ошибки приложения
│   │   └── errors.go
│   ├── application -- сборка компонентов приложения / ручной DI
│   ├── controllers -- обработчики входящих запросов
│   ├── auth -- домен авторизации
│   ├── users -- домен с пользователями
│   │   └── repository
│   ├── config -- конфигурация
│   ├── databases -- обертка над пакетов sql
│   ├── logger
│   └── transport -- http-server / клиенты http, gRPC, если есть
├── migrations -- миграции БД
├── scripts -- вспомогательные скрипты

```

## Запуск приложения

Используются следующие переменные окружения:

```sh
docker-compose up -d
# дождаться завершения миграций
go run cmd/main.go

```

## Разработка проекта

Для разработки могут понадобиться дополнитльные инструменты. Их можно установить через:
```sh
# инструменты будут скачены в папку ./tools/bin
make install-tools
```

## Вопросы для ревьюера:

### HW1
1. Когда следует использовать оборачивание ошибок при помощи директивы %w, а когда нет? Правильно ли я использовал эту директиву в коде? 
2. Нужно/можно ли использовать кастомные ошибки на текущем этапе разработки? В каких случаях это будет уместно?
3. Почему в бойлерплейте проекта контекст для сервера передается в виде параметра closer функции stopHTTPServer в пакете application? Выглядит так, что он вряд ли будет часто меняться и можно его переместить в пакет server.
4. Контекст - пункт задания: Используется логгер и он настраивается из конфига. Вопрос: какие настройки логгера следует выносить в конфиг? По какому принципу выносить?
5. Правильно ли реализован и использован логгер? Рекомендации по улучшению?
6. Как лучше реализовать конфиг приложения?
В описании задания написано:
```markdown
Обычно для конфигурирования приложения используется несколько источников:
- переменные окружения
- файл конфигурации
- внешние системы хранения конфигураций (consul, vault)
```

Вопрос: Какой оптимальный способ реализации конфига в реальном проекте?

godotenv - загрузка переменных окружения из файла .env.
"github.com/caarlos0/env/v11" - парсинг переменных окружения в структуру.

Как я понял, в боевом проекте нужно использовать оба пакета (.yml иди .env) + в зависимости от переменной окружения, подтягивать тот или иной .env. Типа .env.local, .env.dev, .env.prod.

7. На текущий момент запуск приложения golang не докеризируется и приложение запускается напрямую на хосте. Это распространенный подход для локальной среды разработки или лучше добавить использование Dockerfile в docker-compose.yml? Или же Dockerfile здесь нужен только для запуска в продакшен? Хотелось бы узнать про это подробнее. Какие есть практики касательно локальной разработки и докера? Можно ли разрабатывать полностью  в докере? Плюсы и минусы?

8. Имеет ли смысл на данном этапе разработки добавлять тесты? И если да, то для чего их конкретно уже можно написать? Для чего обычно пишутся тесты в проектах на голанг: только для доменной логики или также для транспортного уровня и других пакетов?

#### Доработки после ревью:
1. добавлен линтер golangci-lint с файлом конфигурации
2. добавлен pre-commit git hook для исправления ошибок линтинга и добавления изменений в коммит
3. добавлен makefile target для инициализации pre-commit hook
4. добавлен makefile targets для линтинга и исправления ошибок линтинга
