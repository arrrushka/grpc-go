# Build executable stage
FROM golang
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o server calculator/calc_server/*.go
ENTRYPOINT /app/server
# Build final image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=0 /app/server .
EXPOSE 50052
ENTRYPOINT ["./server"]
