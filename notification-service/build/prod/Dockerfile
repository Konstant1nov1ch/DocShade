FROM golang:1.22.1-bullseye AS base

ARG CONFIG_FILE
ARG GITLAB_TOKEN
ARG GITLAB_USER


# Настройка GOPRIVATE и добавление токена доступа
ENV GOPRIVATE=gitlab.com/docshade/*

WORKDIR /src

# Настройка Git для использования токена доступа
RUN git config --global https.emptyAuth true
RUN git config --global url."https://${GITLAB_USER}:${GITLAB_TOKEN}@gitlab.com/".insteadOf "https://gitlab.com/"

COPY . .
COPY build/prod/${CONFIG_FILE} ./config.yaml

# Скачивание зависимостей
RUN --mount=type=cache,target=/go/pkg/mod go mod download

FROM base AS build

ENV CGO_ENABLED=0
RUN go mod tidy
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    go build \
        -o /go/bin/ \
        .

FROM gcr.io/distroless/static
COPY --from=build /go/bin/notification-service /bin/notification-service
ENTRYPOINT [ "/bin/notification-service" ]
EXPOSE 8080 8000