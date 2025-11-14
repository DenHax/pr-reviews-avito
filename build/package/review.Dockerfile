FROM golang:1.25.3-alpine AS builder

WORKDIR /usr/local/src

RUN apk add --no-cache bash git make gcc gettext musl-dev

COPY ./go.mod ./go.sum .

RUN go mod download

COPY ./ ./

ARG VERSION=unknown

RUN go build -o ./build/bin/review -ldflags="-X main.version=$VERSION"  ./cmd/review/main.go


FROM alpine:3.22 AS runner

RUN apk add --no-cache ca-certificates postgresql-client

COPY ./scripts/wait-storage.sh /wait.sh

RUN chmod 744 /wait.sh

COPY ./configs/config.yaml /config.yaml

ENV CONFIG_PATH /config.yaml

COPY --from=builder /usr/local/src/bin/apw /usr/bin/review

ENV APP_VERSION=$VERSION

CMD ["/usr/bin/review"]
