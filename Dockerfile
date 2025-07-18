FROM golang:1.24.4-alpine as base
WORKDIR /root/openslides-autoupdate-service

RUN apk add git

COPY go.mod go.sum ./
RUN go mod download

COPY main.go main.go
COPY internal internal

# Build service in seperate stage.
FROM base as builder
RUN go build


# Test build.
FROM base as testing

RUN apk add build-base

CMD go vet ./... && go test -test.short ./...


# Development build.
FROM base as development

RUN ["go", "install", "github.com/githubnemo/CompileDaemon@latest"]
EXPOSE 9012

WORKDIR /root
CMD CompileDaemon -log-prefix=false -build="go build -o autoupdate-service ./openslides-autoupdate-service" -command="./autoupdate-service"


# Productive build
FROM scratch

LABEL org.opencontainers.image.title="OpenSlides Autoupdate Service"
LABEL org.opencontainers.image.description="The Autoupdate Service is a http endpoint where the clients can connect to get the current data and also updates."
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.source="https://github.com/OpenSlides/openslides-autoupdate-service"

COPY --from=builder /root/openslides-autoupdate-service/openslides-autoupdate-service .
EXPOSE 9012
ENTRYPOINT ["/openslides-autoupdate-service"]
HEALTHCHECK CMD ["/openslides-autoupdate-service", "health"]
