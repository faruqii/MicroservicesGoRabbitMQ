FROM golang:1.22.3-alpine

RUN apk update && apk add --no-cache git

WORKDIR /go/src/app

COPY . .


RUN go mod download && go mod verify

RUN go mod tidy && go mod vendor

RUN chmod +x /go/src/app

ENV WAIT_VERSION 2.7.3
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait

RUN go build -o /server main.go

EXPOSE 3001

CMD ["/server"]