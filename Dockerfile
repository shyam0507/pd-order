# syntax=docker/dockerfile:1

FROM golang:alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o /pd-order ./src

EXPOSE 3005

CMD [ "/pd-order" ]