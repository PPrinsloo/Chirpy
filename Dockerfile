FROM golang:1.24-alpine AS builder
WORKDIR /build
COPY go.mod ./
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o chirpy

FROM alpine:latest
# Set working directory to match our execution environment
WORKDIR /srv
# Copy binary to working directory
COPY --from=builder /build/chirpy ./
# Copy static files maintaining the same structure as local
COPY ./app/ ./app/
EXPOSE 8080
CMD ["./chirpy"]