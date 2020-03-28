FROM golang:1.12.5-alpine

RUN apk add --no-cache ca-certificates
ADD . /app/
WORKDIR /app
EXPOSE 3000

ENTRYPOINT SERVICE_HOST=borges.ai SERVICE_URL=https://borges.ai TEMPLATES_DIR=templates STATICS_DIR=static PORT=3000 ./app

