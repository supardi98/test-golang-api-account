FROM golang:1.23 AS build

WORKDIR /app
COPY . .
RUN go clean --modcache
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build main.go

FROM alpine:latest

RUN apk add --no-cache curl

WORKDIR /root
COPY --from=build /app/main .

# Set environment variables
ENV DB_HOST=db
ENV DB_USER=postgres
ENV DB_PASSWORD=postgres
ENV DB_NAME=account_db
ENV DB_PORT=5432

EXPOSE 8080
CMD ["./main"]