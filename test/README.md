# test

## Назначение каталога
- Тестовый слой (Layered support): интеграционные сценарии и вспомогательные утилиты для проверки CRUD/поиска событий на реальной БД (`test/integration/events_service_and_repository_test.go`).
- Application Core стиля у тестов нет: они проверяют Anemic/Layered прод-код через реальные инфраструктурные зависимости.
- Содержит генерацию моков для сервисного интерфейса HTTP-слоя (`test/mocks/events_service_mock.go`).

## Правила зависимостей
- Интеграционные тесты используют `testcontainers-go` и реальные миграции (`test/integration/events_service_and_repository_test.go:23-32`); не должны зависеть от HTTP/джобов.
- Утилиты поднимают Postgres и гоняют goose миграции (`test/util/db.go:24-98`); возвращают `*database.DB` из прод-кода.
- Моки генерируются по интерфейсу `eventService` и применяются только в тестах (`test/mocks/events_service_mock.go`).

## Конвенции и практики
- Каждый suite поднимает контейнер, применяет миграции и закрывает ресурсы (`test/integration/events_service_and_repository_test.go:23-33`).
- CRUD/поиск проверяются через реальный репозиторий событий (`test/integration/events_service_and_repository_test.go:34-153`).
- Хелперы из `test/util` изолируют инфраструктурный сетап и делают тесты самодостаточными (`test/util/db.go:24-98`).
- CQRS не выделен: тесты используют единый сервис/репозиторий; для читателей это отражает фактический Anemic+Layered.

## Ключевые артефакты
- `integration/events_service_and_repository_test.go` — CRUD/поиск ближайших событий на реальной БД (`test/integration/events_service_and_repository_test.go:34-153`).
- `util/db.go` — создание контейнера Postgres, запуск миграций и закрытие (`test/util/db.go:24-98`).
- `mocks/events_service_mock.go` — gomock-генерация интерфейса сервиса событий.

## Нарушения/риски
- Покрытие ограничено: нет тестов фоновых джобов или HTTP-роутов (покрыт только репозиторий `events`), риск регрессий в API/джобах — значимо.
- Маршруты update/delete в HTTP тестируются только частично через моки, интеграционный слой их не проверяет (`test/integration/events_service_and_repository_test.go:34-153`) — значимо.
- Запуск зависит от Docker/сети; в ограниченных CI тесты упадут без подготовленной среды (`test/util/db.go:24-84`) — минорно, требует настройки.
