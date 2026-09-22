FROM golang:1.26-alpine AS my-build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /order-api ./cmd

FROM alpine
WORKDIR /app
COPY --from=my-build /order-api .
EXPOSE 8081
CMD ["./order-api"]
