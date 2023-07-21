FROM golang:1.20.5-alpine3.17

ENV BOT_TOKEN ${BOT_TOKEN}
ENV SID ${SID}
ENV ADMINS ${ADMINS}

WORKDIR /app

COPY . .

RUN go build ./cmd/app

CMD ["./app"]
