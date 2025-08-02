# Stage 1: Build Stage
FROM --platform=$BUILDPLATFORM golang:1.23-alpine3.20 AS builder

WORKDIR /app
ARG BIN
ARG TARGETOS
ARG TARGETARCH
ENV GOPRIVATE=gitlab.com/a7923/
ENV GONOPROXY=gitlab.com/a7923/
ENV GONOSUMDB=gitlab.com/a7923/
ENV GO111MODULE=on

RUN apk add -qU openssh git

# Set up SSH for GitLab access
RUN mkdir -p -m 0600 ~/.ssh && ssh-keyscan gitlab.com >> ~/.ssh/known_hosts
RUN git config --global url.git@gitlab.com:.insteadOf https://gitlab.com/

# Pre-cache dependencies
COPY go.mod go.sum ./

RUN --mount=type=ssh go mod download && go mod verify

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