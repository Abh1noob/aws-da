FROM golang:alpine

RUN apk update && apk upgrade && \
    apk add --no-cache git

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]