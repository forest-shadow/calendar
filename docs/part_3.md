# (в доработке) Домашнее задание. Часть 3

Итоговый результат описан в [ТЗ](calendar.md). 
В этой части домашнего задания мы подключим очереди к работе нашего сервиса.

## Компоненты для работы с очередью

В качестве очереди мы будем использоваться kafka.

Есть 2 популярные библиотеки для работы с kafka:
- [sarama](https://github.com/IBM/sarama)
- [franz-go](https://github.com/twmb/franz-go) (рекомендуемая)

Необходимо разработать producer (компонентб отправляющий сообщения в кафку) и consumer (компонент читающий сообщения из кафки).

Для поднятия окружения добавьте в docker-compose:
```yaml
  zookeeper:
    image: confluentinc/cp-zookeeper:5.4.0
    hostname: zookeeper
    container_name: zookeeper
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000

  kafka:
    image: confluentinc/cp-kafka:5.3.0
    hostname: kafka
    container_name: kafka
    environment:
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka1:9092,PLAINTEXT_HOST://localhost:19092
      KAFKA_ZOOKEEPER_CONNECT: "zookeeper:2181"
      KAFKA_BROKER_ID: 1
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
      KAFKA_CONFLUENT_LICENSE_TOPIC_REPLICATION_FACTOR: 1
    ports:
        - 19092:19092
    depends_on:
      - zookeeper
    healthcheck:
      test: nc -z localhost 9092 || exit -1
      interval: 5s
      timeout: 10s
      retries: 10

```

## Отправка уведомлений через Kafka
В первой версии нашего приложения отравитель вычитывал данные из базы и сразу же отправлял уведомления.

Необходимо разделить отправить на 2 части:
1. вычитывает сообщения из БД и отправляет их в очередь
2. читает сообщения из очереди и пишет текст сообщения в лог также устанавливая дату отпаврки

## Критерии приемки задачи

* При запуске приложения подключается Kafka и она не крашится. В Docker-Compose прописано все для Zookeeper, Kafka **- 1 балл**
* Приложение вычитывает события из БД и отправляет их в Kafka (producer) **- 2 балла**
* Приложение читает сообщения из очереди (consumer) и пишет текст сообщения в лог, также устанавливая дату отправки в БД **- 2 балла**

Зачет от **3 баллов**