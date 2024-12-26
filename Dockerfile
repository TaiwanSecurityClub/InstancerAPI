FROM golang:1.23-alpine as builder

RUN apk add --no-cache make build-base

WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
WORKDIR /
RUN rm -r /src
COPY . /src
WORKDIR /src
RUN make clean && make

FROM docker:dind as release

COPY --from=builder /src/bin /app

WORKDIR /app

RUN touch .env

COPY ./docker-entrypoint.sh ./docker-entrypoint.sh

RUN chmod +x ./docker-entrypoint.sh

ENTRYPOINT ["./docker-entrypoint.sh"]
