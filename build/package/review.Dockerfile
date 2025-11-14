FROM golang:1.25.3-alpine AS builder

WORKDIR /usr/local/src

RUN apk add --no-cache bash git make gcc gettext musl-dev

COPY ./go.mod ./go.sum .

RUN go mod download

COPY ./ ./

ARG TAG=unknown

RUN go build -o ./build/bin/review -ldflags="-X main.version=$TAG"  ./cmd/review/main.go


FROM alpine:3.22 AS runner

RUN apk add --no-cache ca-certificates

COPY ./configs/config.yaml /etc/review/config.yaml

ENV POSTGRES_HOST=
ENV POSTGRES_PORT=
ENV POSTGRES_USER=
ENV POSTGRES_DB=
ENV POSTGRES_PASSWORD=
ENV POSTGRES_SSL_MODE=disable
ENV CONFIG_PATH /etc/review/config.yaml

ENV GIN_MODE=

ENV APP_VERSION=$TAG

COPY --from=builder /usr/local/src/build/bin/review /usr/bin/review

CMD ["/usr/bin/review"]
