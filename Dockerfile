# syntax=docker/dockerfile:1

## Build
FROM golang:1.22-bookworm AS build
LABEL org.opencontainers.image.source https://github.com/hadret/forwardly-go

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY main.go ./

RUN go build -o /forwardly-go


## Deploy
FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=build /forwardly-go /forwardly-go

EXPOSE 8000

USER nonroot:nonroot

CMD [ "/forwardly-go" ]
