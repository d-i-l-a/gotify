FROM golang:1.22-alpine as build

ENV CGO_ENABLED=1
ENV GOCACHE=/root/.cache/go-build

WORKDIR /work
COPY . .

RUN apk add \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev
RUN apk add make npm nodejs

RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate
RUN mkdir -p app/internal/static/css
RUN mkdir -p app/internal/static/js
RUN npx tailwindcss -o ./app/internal/static/css/output.css
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o app/main main.go
COPY ./internal/static/js/player.js ./app/internal/static/js/player.js


FROM alpine:latest as run
RUN apk --no-cache add ca-certificates
WORKDIR /run
COPY --from=build /work/app .
EXPOSE 80

#ENTRYPOINT ["."]
CMD ["./main"]