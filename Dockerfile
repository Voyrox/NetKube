FROM golang:1.26.1-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/netkube .

FROM alpine:3.21

WORKDIR /app

RUN addgroup -S netkube && adduser -S netkube -G netkube

COPY --from=builder /out/netkube /app/netkube
COPY views /app/views
COPY public /app/public
COPY reference /app/reference

RUN mkdir -p /app/config/uploaded-sources && chown -R netkube:netkube /app

USER netkube

EXPOSE 3000

CMD ["/app/netkube"]
