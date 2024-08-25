FROM golang:1.23.0-alpine3.19 as builder

# Set necessary environment variables for the Go proxy
# ENV GO111MODULE=on
# ENV GOPROXY=https://goproxy.io,direct

WORKDIR /app

COPY . .

RUN go build -o wrapto .

EXPOSE 3000

CMD ["./wrapto"]
