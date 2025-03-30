FROM golang:1.24.0-alpine3.20 AS builder

ENV GOCACHE=/root/.cache/go-build

WORKDIR /app

COPY ./consumer .
COPY ./.env .
COPY ./storage ./storage
COPY ./go.mod .

RUN go mod tidy

RUN go build -o consumer main.go

FROM alpine:3.20

WORKDIR /myapp

COPY --from=builder /app/consumer ./
COPY --from=builder /app/.env ./

ARG CONSUMER_NAME_DEFAULT  
RUN echo -e "\nCONSUMER_NAME=${CONSUMER_NAME_DEFAULT}" >> .env

CMD [ "/myapp/consumer" ]