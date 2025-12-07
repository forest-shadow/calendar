# Calendar

## Назначение
- Сервис календаря для CRUD по событиям, отправки уведомлений и периодической очистки; используется локально через `docker-compose` (Postgres, Kafka) и запускается из `cmd/main.go`.
- Репозиторий служит учебным примером: демонстрирует HTTP API на chi, фоновые джобы и работу с миграциями goose.

## Архитектура
- Architecture Approach (Layered/гибрид):
  - Сборка зависимостей, транспорты и джобы находятся в одном модуле композиции без портов (`internal/application/app.go:24-74`), что типично для Layered.
  - Входной адаптер обращается к сервису напрямую, без портов/DTO уровня домена (`internal/controllers/http/router.go:44-51`).
  - Доменный сервис тянет инфраструктуру (логгер) и конкретное подключение к БД через репозиторий (`internal/events/events_service.go:23-60`, `internal/events/events_repo.go:12-38`), что закрепляет зависимость вниз по слоям.
- Application Core (Transaction Script / Anemic):
  - Методы сервиса — прямые прокси к репозиторию без инвариантов и бизнес-правил (`internal/events/events_service.go:32-62`).
  - Джобы используют те же CRUD/поиск/удаление без обогащения доменной логикой (`internal/jobs/event_notifier.go:56-116`, `internal/jobs/event_cleaner.go:24-54`).
  - CQRS отсутствует: единый интерфейс `eventService` обслуживает и команды, и запросы (`internal/controllers/http/router.go:16-22`).
- Матрица Application Core × Architecture Approach:
  - Anemic + Layered — ✔ для учебного CRUD/джобов: минимальный бойлерплейт, быстрая сборка окружения (`internal/application/app.go:24-74`).
  - Anemic + Layered — △ по DIP/тестируемости: домен зависит от infra и логгера, контракты сервиса размножены (`internal/application/domain_events.go:16-35`, `internal/jobs/event_notifier.go:17-36`, `test/mocks/events_service_mock.go`), усложняя моки и эволюцию API.
  - Anemic + Layered — ✖ для усложнения домена: отсутствие use-case слоя и инвариантов, нет Result/ошибок домена и политики устойчивости (`internal/resilience` пуст), поэтому рост правил потребует перехода к Clean/Hex.

## Слои/модули
- `cmd` — точка входа, создаёт корневой контекст с отменой по сигналу и делегирует запуск `internal/application` (`cmd/main.go`).
- `internal/application` — композиция зависимостей, запуск HTTP и джобов, управление жизненным циклом (`internal/application/app.go`, `internal/application/domain_events.go`).
- `internal/controllers/http` — REST-ручки chi, DTO, базовая валидация `notify_before`, маппинг ошибок `apperrs` в статусы (`internal/controllers/http/router.go`, `internal/controllers/http/events.go`).
- `internal/transport/http` — обёртка над `net/http.Server`, listener и shutdown с таймаутом (`internal/transport/http/server.go`).
- `internal/events` — доменная модель события и репозиторий Postgres, тонкий сервис без инвариантов (`internal/events/event.go`, `internal/events/events_repo.go`, `internal/events/events_service.go`).
- `internal/jobs` — фоновые задачи уведомлений (Kafka producer/consumer) и очистки по тикеру (`internal/jobs/event_notifier.go`, `internal/jobs/event_cleaner.go`, `internal/jobs/jobs_runner.go`).
- `internal/database` — создание и закрытие `*sql.DB` через pgx (`internal/database/db.go`).
- `internal/config` — загрузка YAML/ENV через viper (`internal/config/config.go`).
- `internal/logger` — zap-логгеры, уровень по окружению, отдельный логгер с ротацией для нотификатора (`internal/logger/logger.go`, `internal/logger/event_notifier_logger.go`).
- `internal/apperrs` — sentinel-ошибки для маппинга статусов (`internal/apperrs/errors.go`).
- `internal/resilience` — резерв под устойчивость (пока пуст).
- `docs` — текст требований и заданий по сервису (`docs/calendar.md`).
- `migrations` — схема таблицы `events` и сиды для примеров (`migrations/20241016191005_create_events.sql`).
- `scripts` — миграционный Dockerfile, hurl-сценарии, pre-commit сетап (`scripts/README.md`).
- `test` — интеграционные тесты (testcontainers) и моки для HTTP (`test/integration/events_service_and_repository_test.go`, `test/mocks/events_service_mock.go`).

## Правила зависимостей и практики
- Зависимости сверху вниз: `internal/application` собирает конфиг/логгер/БД и передаёт сервис событий в HTTP и джобы; прямые обращения к репозиторию допускаются только из домена (`internal/application/app.go` → `buildEventsDomain`).
- Входные адаптеры должны оставаться тонкими: парсят ввод, валидируют `notify_before`, маппят ошибки (`internal/controllers/http/events.go`, `internal/controllers/http/utils.go`); бизнес-логики в них нет.
- CQRS не выделен: команды/запросы объединены в интерфейсе `eventService`; риск роста связанности при усложнении запросов.
- Валидация ограничена (формат длительности), обязательные поля/границы дат не проверяются до вызова сервиса.
- Use-case слой отсутствует: сервис — прокси; новая бизнес-логика должна добавляться в домен и оформляться контрактами, а не в контроллерах/джобах.
- Единственный источник схемы — миграции `migrations/`; конфиг — `env.*.yml` с перекрытием `APP_*`.
- Ошибки — sentinel из `internal/apperrs/errors.go`, маппятся в HTTP-коды в `internal/controllers/http/utils.go`; Result-типов/доменных кодов нет.

## Достоинства
- Простая цепочка Layered упрощает локальный запуск и чтение кода; минимальный порог для студентов/джунов.
- Композиция вынесена в один модуль (`internal/application/app.go`), что облегчает замену инфраструктуры и конфигурацию окружений.
- Наличие интеграционных тестов на реальной БД и моков для контроллеров повышает уверенность в CRUD и базовой валидации.
- Логирование/конфиг централизованы и легко расширяются; нотификатор имеет отдельный лог с ротацией.
- Инфраструктура автоматизирована: `docker-compose`, миграционный контейнер, Makefile-цели и предустановленные тулзы в `tools/bin`.

## Недостатки/отклонения
- Домен зависит от инфраструктуры (`internal/events/events_service.go` использует `logger.Logger`, `internal/events/events_repo.go` требует `database.DBConnection`), порты не выделены.
- Контракты сервиса размножены (`internal/application/domain_events.go`, `internal/controllers/http/router.go`, `internal/jobs/event_notifier.go`, `test/mocks/events_service_mock.go`), что усложняет эволюцию API и регенерацию моков.
- HTTP-слой нарушает семантику: DELETE без `{id}` (`internal/controllers/http/router.go`), GET списка читает тело (`internal/controllers/http/events.go`), запись тела до `WriteHeader` (`getEvent`, `getEventsList`).
- Джобы блокируют старт приложения и завершают его ошибкой при штатном `context.Done` (`internal/jobs/jobs_runner.go` → `internal/application/app.go`), нет ретраев/таймаутов для Kafka/БД.
- `notify_before` строкой без конверсии к `INTERVAL` и без типовой валидации БД; `DeleteOldEvents` удаляет записи старше ~30 секунд вместо года (`internal/events/events_repo.go` vs `docs/calendar.md`).
- Миграция использует `gen_random_uuid()` без `pgcrypto`, сиды не применятся на чистой БД (`migrations/20241016191005_create_events.sql`).
- Нет слоя устойчивости (`internal/resilience` пуст), нет единой модели ошибок/Result, отсутствуют доменные инварианты (пересечение дат, обязательные поля).

## Соответствие назначению
- Учебный CRUD/джобы: выбранная пара Anemic + Layered (✔) оправдана простотой, быстрой сборкой окружения и низкой стоимостью входа.
- Эволюция к сложным правилам/новым адаптерам: текущая комбинация даёт △/✖ — придётся выделять порты и use-case слой, вводить инварианты, типизированные ошибки/Result и устойчивость (ретраи, таймауты), иначе регрессии по API/данным будут частыми.
- Нагрузочные/безопасные сценарии: отсутствуют таймауты/лимиты сервера (`internal/transport/http/server.go`) и ретраи, поэтому под нагрузкой или при сбоях Kafka/БД сервис будет хрупким.

## Ключевые артефакты
- `cmd/main.go` — точка входа, корневой контекст.
- `internal/application/app.go` — сборка зависимостей, запуск HTTP и джобов.
- `internal/controllers/http/router.go` — маршруты и middleware REST API.
- `internal/controllers/http/events.go` — CRUD-хендлеры, DTO, базовая валидация.
- `internal/events/events_service.go`, `internal/events/events_repo.go` — сервис и репозиторий событий.
- `internal/jobs/event_notifier.go`, `internal/jobs/event_cleaner.go` — фоновые задачи уведомлений и очистки.
- `internal/transport/http/server.go` — HTTP-сервер и graceful shutdown.
- `migrations/20241016191005_create_events.sql` — схема таблицы `events`.
- `Makefile`, `docker-compose.yaml`, `env.local.yml` — запуск и локальное окружение.

## Нарушения/риски
- DELETE-ручка зарегистрирована без `{id}` (`internal/controllers/http/router.go:44-50`), клиент не сможет удалить событие — критично.
- GET списка читает тело вместо query-параметров (`internal/controllers/http/events.go:118-147`), большинство клиентов/прокси его отбросят — критично.
- Статусы ответа пишутся после тела (`internal/controllers/http/events.go:101-147`), из-за чего HTTP-коды могут быть неверными — критично.
- Очистка удаляет события старше ~30 секунд, а не года (`internal/events/events_repo.go:155-173` против требований `docs/calendar.md`) — значимо.
- `RunJobs` возвращает ошибку при штатном завершении контекста (`internal/jobs/jobs_runner.go:7-16` → `internal/application/app.go:66-100`), ломая graceful shutdown — значимо.
- `notify_before` передаётся строкой без приведения к `INTERVAL` (`internal/events/events_repo.go:23-44`, `internal/events/event.go:11-25`), возможны ошибки записи и рассинхрон с БД — значимо.
- Сиды миграции используют `gen_random_uuid()` без включения `pgcrypto` (`migrations/20241016191005_create_events.sql:16-20`), «чистая» БД не применит миграцию — значимо.
- Продюсер/консьюмер Kafka и тикеры работают без ретраев/таймаутов и падают по первой ошибке (`internal/jobs/event_notifier.go:85-115`, `internal/jobs/event_cleaner.go:24-54`); каталог `internal/resilience` пуст — минорно, но снижает отказоустойчивость.
