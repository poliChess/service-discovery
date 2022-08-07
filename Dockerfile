# Build Step
FROM golang:1.18-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build ./cmd/main.go

RUN chmod +x /app/main


# Run Step
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/main /app

CMD [ "/app/main" ]
