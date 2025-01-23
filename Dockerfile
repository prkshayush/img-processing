FROM golang:1.23.5-alpine

WORKDIR /app

# first copy only mod to track dependency change
COPY  go.mod go.sum ./

RUN go mod download

# then copy all
COPY . .

RUN go build -o main

EXPOSE 8080

CMD ["./main"]