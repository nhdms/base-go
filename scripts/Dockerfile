# Stage 1: Build Stage
FROM --platform=$BUILDPLATFORM golang:1.23-alpine3.20 AS builder

WORKDIR /app
ARG BIN
ARG TARGETOS
ARG TARGETARCH
ENV GO111MODULE=on

RUN apk add -qU openssh git

# Pre-cache dependencies
COPY go.mod go.sum ./

RUN go mod download && go mod verify

# Copy source code and build
COPY . .
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /app/main ${BIN}/main.go
COPY ${BIN}/config /app/config

# Stage 2: Runtime Stage
FROM --platform=$TARGETPLATFORM alpine:3.20 AS runtime
WORKDIR /app
ARG BIN
ENV CONFIG_PATH=/app/config

# Copy the built binary and config file from the build stage
COPY --from=builder /app /app

CMD ["./main"]