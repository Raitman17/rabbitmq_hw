FROM golang:1.24.0-alpine3.20 AS builder

ENV GOCACHE=/root/.cache/go-build

WORKDIR /app

COPY ./producer .
COPY ./.env .
COPY ./storage ./storage
COPY ./go.mod .

RUN go mod tidy

RUN go build -o producer main.go

FROM alpine:3.20

WORKDIR /myapp

COPY --from=builder /app/producer ./
COPY --from=builder /app/.env ./
COPY --from=builder /app/text.txt ./

CMD [ "/myapp/producer" ]