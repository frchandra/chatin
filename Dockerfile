# Step 1
FROM golang:latest AS builder
ENV GOPROXY=https://goproxy.io,direct

RUN apt update && apt install -y apt-utils

WORKDIR /go/src
COPY . .
RUN go mod tidy
RUN go mod download -x

RUN CGO_ENABLED=1 go build -o ./bin/app ./cmd/app/main.go
RUN CGO_ENABLED=1 go build -o ./bin/migrator ./cmd/migrator/main.go


# Step 2
FROM debian:stable-slim AS runner
RUN apt update && apt install -y apt-utils

WORKDIR /chatin

COPY --from=builder /go/src/bin /chatin

EXPOSE 5000
CMD ["/chatin/app"]