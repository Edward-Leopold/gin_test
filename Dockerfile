FROM golang:1.24-alpine AS builder
WORKDIR /gin_test
COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /gin_test
COPY --from=builder /gin_test/main .
COPY --from=builder /gin_test/frontend ./frontend
EXPOSE 8080
CMD ["./main"]