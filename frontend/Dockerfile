FROM golang:1.11.11-alpine3.8
WORKDIR /app/client
RUN apk update && apk add git
COPY . .
CMD ["./main"]