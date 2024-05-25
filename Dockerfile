# FROM golang:1.21

# WORKDIR /app

# COPY . .

# RUN go mod download

# ENV TODO_PORT=7540

# RUN GOARCH=amd64 go build -o /todofinal

# EXPOSE 7540

# CMD ["/todofinal"] 


FROM golang:1.21 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN GOARCH=amd64 go build -o /app/todofinal

FROM golang:1.21 

WORKDIR /app

COPY --from=builder /app/todofinal /app/todofinal

COPY --from=builder /app/scheduler.db /app/scheduler.db

ENV TODO_PORT=7540

EXPOSE 7540

CMD ["/app/todofinal"]

