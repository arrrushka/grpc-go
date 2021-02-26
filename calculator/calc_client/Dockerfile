# Build executable stage
FROM golang
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o client calculator/calc_client/*.go
ENTRYPOINT /app/client
# Build final image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=0 /app/client .
ENTRYPOINT ["./client"]
