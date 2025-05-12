FROM golang:1.22-alpine AS builder

RUN apk add --no-cache ca-certificates git gcc g++ libc-dev binutils

WORKDIR /opt

COPY go.mod go.sum ./
RUN go mod tidy && go mod verify

COPY . .

RUN ls -la /opt

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/application ./cmd/main.go

FROM alpine AS runner

RUN apk add --no-cache ca-certificates libc6-compat bash openssh

WORKDIR /opt

COPY --from=builder /opt/bin/application ./
COPY --from=builder /opt/migrations ./migrations
COPY wait-for-it.sh ./
RUN chmod +x wait-for-it.sh

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

CMD ["sh", "-c", "./wait-for-it.sh postgres:5432 --  && ./application"]
