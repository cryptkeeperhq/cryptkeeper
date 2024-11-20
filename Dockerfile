# syntax=docker/dockerfile:1

FROM golang:1.22.3-alpine as builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o /cryptkeeper-app

FROM alpine:latest

WORKDIR /root/

COPY config.yaml /config.yaml
COPY master.key /master.key

# Install wait-for-it script
COPY scripts/wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

COPY --from=builder /cryptkeeper-app .

EXPOSE 8000

CMD ["./cryptkeeper-app"]
# CMD ["/wait-for-it.sh", "postgres:5432", "--", "/cryptkeeper-app"]
