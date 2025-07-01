ARG ALPINE_VERSION=3.20
ARG GO_VERSION=1.24.0
ARG APP_NAME=godyno
ARG APP_PATH=./cmd/dyno
ARG GOCACHE=/root/.cache/go-build
ARG ASM_FLAGS="-trimpath"
ARG GC_FLAGS="-trimpath"
ARG LD_FLAGS_BASE="-w -s -extldflags '-static'"
ARG VERSION=dev

# amd64
FROM --platform=linux/amd64 golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder-amd64
ARG ASM_FLAGS
ARG GC_FLAGS
ARG LD_FLAGS_BASE
ARG GOCACHE
ARG APP_NAME
ARG APP_PATH
ARG VERSION

WORKDIR /go/src/${APP_NAME}
RUN apk add --no-cache upx 

COPY go.mod go.sum ./
RUN --mount=type=cache,target=${GOCACHE} go mod download
COPY . ./

ENV CGO_ENABLED=0
ENV GOARCH=amd64
ENV GOOS=linux

RUN --mount=type=cache,target=${GOCACHE} \
    go build -asmflags="${ASM_FLAGS}" \
             -ldflags="${LD_FLAGS_BASE} -X 'github.com/Mad-Pixels/go-dyno.Version=${VERSION}'" \
             -gcflags="${GC_FLAGS}" \
             -o /bin/${APP_NAME} \
             ${APP_PATH}
RUN upx --best --lzma /bin/${APP_NAME}

# arm64
FROM --platform=linux/arm64/v8 golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder-arm64
ARG ASM_FLAGS
ARG GC_FLAGS
ARG LD_FLAGS_BASE
ARG GOCACHE
ARG APP_NAME
ARG APP_PATH
ARG VERSION

WORKDIR /go/src/${APP_NAME}
RUN apk add --no-cache upx

COPY go.mod go.sum ./
RUN --mount=type=cache,target=${GOCACHE} go mod download
COPY . ./

ENV CGO_ENABLED=0
ENV GOARCH=arm64
ENV GOOS=linux

RUN --mount=type=cache,target=${GOCACHE} \
    go build -asmflags="${ASM_FLAGS}" \
             -ldflags="${LD_FLAGS_BASE} -X 'github.com/Mad-Pixels/go-dyno.Version=${VERSION}'" \
             -gcflags="${GC_FLAGS}" \
             -o /bin/${APP_NAME} \
             ${APP_PATH} 
RUN upx --best --lzma /bin/${APP_NAME}

# Final amd64 image
FROM scratch AS amd64
ARG APP_NAME=godyno
COPY --from=builder-amd64 /bin/${APP_NAME} /${APP_NAME}
ENTRYPOINT ["/godyno"]

# Final arm64 image
FROM scratch AS arm64
ARG APP_NAME=godyno
COPY --from=builder-arm64 /bin/${APP_NAME} /${APP_NAME}
ENTRYPOINT ["/godyno"]