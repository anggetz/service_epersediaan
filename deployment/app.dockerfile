FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

RUN go build -o ./executable -v ./

FROM alpine:3.16

COPY --from=builder /app/.env .
COPY --from=builder /app/executable .

CMD ./executable
