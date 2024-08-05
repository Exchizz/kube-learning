# syntax=docker/dockerfile:1

FROM golang:1.21-alpine

RUN adduser -D runas
USER runas


WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /app/kube-learning

EXPOSE 8080
USER 1000:1000

CMD [ "/app/kube-learning" ]
