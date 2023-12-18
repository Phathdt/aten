FROM golang:1.21.3-alpine as builder
WORKDIR /app
RUN apk update && apk add --no-cache gcc musl-dev git
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
WORKDIR /app
RUN go build -ldflags '-w -s' -a -o aten ./main.go

# ----------------------
FROM golang:1.21.3-alpine as migrate
WORKDIR /app
RUN apk update && apk add --no-cache gcc musl-dev git
COPY migrate/go.mod migrate/go.sum ./
RUN go mod download
COPY migrate/ .
RUN go build -ldflags '-w -s' -a -o migrate main.go

# Deployment environment
# ----------------------
FROM alpine:3.18.4
WORKDIR /app
RUN chown nobody:nobody /app
USER nobody:nobody
COPY --from=builder --chown=nobody:nobody ./app/aten .
COPY --from=builder --chown=nobody:nobody ./app/run.sh .
COPY --from=migrate --chown=nobody:nobody ./app/migrate .
COPY --from=migrate --chown=nobody:nobody ./app/migrations ./migrations

ENTRYPOINT sh run.sh
