# -- COMPILE -- #
FROM golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/backend/main.go


# -- RUN -- #
FROM alpine:3.20

WORKDIR /app

COPY --from=build /app/app /app/app

EXPOSE 5000
ENTRYPOINT ["/app/app"]
