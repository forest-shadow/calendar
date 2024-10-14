
# default value for build dir
FROM golang:1.23.0 as builder

ENV CGO_ENABLED=0
ADD . /app/
WORKDIR /app

# Можно использовать кеширование зависимостей для ускорения сборок
# --mount=type=cache,target=/go/pkg/mod \
# --mount=type=cache,target=/root/.cache/go-build \
RUN go build -o calendar cmd/main.go

FROM  alpine:3.19.1
ENV APP_ENV=prod
COPY --from=builder /app/calendar /app/calendar

WORKDIR /app
RUN cat <<EOF > env.prod.yml
http:
  port: 8080
db:
  uri: "postgresql://postgres:password@postgres:5432/calendar"
EOF

EXPOSE 8080/tcp
CMD ["/app/calendar"]