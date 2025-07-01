ARG ALPINE_VERSION=3.20
ARG GO_VERSION=1.24.0
ARG APP_NAME=godyno
ARG APP_PATH=./cmd/dyno
ARG GOCACHE=/root/.cache/go-build
ARG ASM_FLAGS="-trimpath"
ARG GC_FLAGS="-trimpath"
ARG LD_FLAGS_BASE="-w -s -extldflags '-static'"
ARG VERSION=dev

# Builder
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder-base
ARG ASM_FLAGS
ARG GC_FLAGS
ARG LD_FLAGS_BASE
ARG GOCACHE
ARG APP_NAME
ARG APP_PATH
ARG VERSION

WORKDIR /go/src/${APP_NAME}
COPY go.mod go.sum ./
COPY . ./
RUN apk add --no-cache upx
RUN --mount=type=cache,target=${GOCACHE} \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN --mount=type=cache,target=${GOCACHE} \
    --mount=type=cache,target=/go/pkg/mod \
    go install golang.org/x/tools/cmd/goimports@latest && \
    go install mvdan.cc/gofumpt@latest

# amd64
FROM builder-base AS builder-amd64
ENV GOARCH=amd64

RUN --mount=type=cache,target=${GOCACHE} \
    --mount=type=cache,target=/go/pkg/mod \
    go build -asmflags="${ASM_FLAGS}" \
             -ldflags="${LD_FLAGS_BASE} -X 'github.com/Mad-Pixels/go-dyno.Version=${VERSION}'" \
             -gcflags="${GC_FLAGS}" \
             -o /bin/${APP_NAME} \
             ${APP_PATH}
RUN upx --best --lzma /bin/${APP_NAME}

# arm64
FROM builder-base AS builder-arm64
ENV GOARCH=arm64

RUN --mount=type=cache,target=${GOCACHE} \
    --mount=type=cache,target=/go/pkg/mod \
    go build -asmflags="${ASM_FLAGS}" \
             -ldflags="${LD_FLAGS_BASE} -X 'github.com/Mad-Pixels/go-dyno.Version=${VERSION}'" \
             -gcflags="${GC_FLAGS}" \
             -o /bin/${APP_NAME} \
             ${APP_PATH}
RUN upx --best --lzma /bin/${APP_NAME}

FROM alpine:${ALPINE_VERSION} AS runtime-base
ARG APP_NAME=godyno

RUN apk add --no-cache go && rm -rf /var/cache/apk/*
COPY --from=builder-base /go/bin/goimports /usr/local/bin/
COPY --from=builder-base /go/bin/gofumpt /usr/local/bin/

RUN adduser -D -s /bin/sh ${APP_NAME}
USER ${APP_NAME}

# Final amd64
FROM runtime-base AS amd64
ARG APP_NAME=godyno
COPY --from=builder-amd64 /bin/${APP_NAME} /usr/local/bin/${APP_NAME}
ENTRYPOINT ["godyno"]

# Final arm64
FROM runtime-base AS arm64
ARG APP_NAME=godyno
COPY --from=builder-arm64 /bin/${APP_NAME} /usr/local/bin/${APP_NAME}
ENTRYPOINT ["godyno"]