# build
FROM golang:1.24-alpine AS build
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
ARG SERVICE
RUN go build -o /bin/app ./cmd/${SERVICE}

# runtime
FROM alpine:3.20
WORKDIR /app

COPY --from=build /bin/app /app/app
COPY --from=build /app/ui /app/ui

EXPOSE 8080
ENTRYPOINT ["/app/app"]
