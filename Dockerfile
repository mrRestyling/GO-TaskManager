FROM golang:1.21 AS builder

WORKDIR /app

COPY . .

RUN go mod download

ENV CGO_ENABLED=1

RUN go build -o /app/todofinal

FROM golang:1.21 

WORKDIR /app

COPY --from=builder /app/todofinal /app/todofinal

COPY --from=builder /app/scheduler.db /app/scheduler.db

COPY --from=builder /app/web/. /app/web/.

ENV TODO_PORT=7540

EXPOSE 7540

CMD ["/app/todofinal"]

