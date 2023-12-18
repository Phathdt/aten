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
FROM node:20.10-alpine as migrator
WORKDIR /app
COPY migrator/package.json .
COPY migrator/yarn.lock .
RUN yarn install
COPY . .

# Deployment environment
# ----------------------
FROM node:20.10-alpine
WORKDIR /app
RUN chown nobody:nobody /app
USER nobody:nobody
COPY --from=builder --chown=nobody:nobody ./app/aten .
COPY --from=builder --chown=nobody:nobody ./app/run.sh .
COPY --from=migrator --chown=nobody:nobody ./app ./migrator
RUN ls -lha

ENTRYPOINT sh run.sh
