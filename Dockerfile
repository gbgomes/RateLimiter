FROM golang:1.21.0 as builder

WORKDIR /app

COPY . .

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o ratelimiter ./cmd/server


FROM scratch
COPY --from=builder /app/ratelimiter .
COPY --from=builder /app/cmd/server/.env .
COPY --from=builder /app/cmd/server/tokens.json .

CMD ["./ratelimiter"]
