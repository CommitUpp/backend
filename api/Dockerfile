# syntax=docker/dockerfile:1
# ==============================================================================
# ▼ developステージ
# ==============================================================================
FROM --platform=$TARGETPLATFORM golang:1.26-alpine AS develop 

RUN --mount=type=cache,target=/var/cache/apk \
    apk update && apk add --no-cache git curl ca-certificates openssh

WORKDIR /app

RUN go install github.com/air-verse/air@v1.61.7

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

CMD ["air", "-c", ".air.toml"]

# ==============================================================================
# ▼ builderステージ
# ==============================================================================
FROM develop AS builder

WORKDIR /app

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -trimpath \
    -o /sheater-api \
    ./cmd/main.go

# ==============================================================================
# ▼ runnerステージ
# ==============================================================================
FROM gcr.io/distroless/static-debian13:nonroot AS runner

WORKDIR /app

COPY --from=builder /sheater-api ./sheater-api

USER nonroot:nonroot

ENTRYPOINT [ "./sheater-api" ]