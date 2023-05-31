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
COPY storage /chatin/storage
#COPY .env /chatin/
#COPY .senpro-381803-cffcf39a0935.json /chatin/senpro-381803-cffcf39a0935.json

EXPOSE 5000
CMD ["/chatin/app"]