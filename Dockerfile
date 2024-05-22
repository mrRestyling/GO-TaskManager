FROM golang:1.21
WORKDIR /app
COPY . .
RUN go mod download
ENV TODO_PORT=7540
RUN GOARCH=amd64 go build -o /todofinal
EXPOSE 7540
CMD ["/todofinal"] 
