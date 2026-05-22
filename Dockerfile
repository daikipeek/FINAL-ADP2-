FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
ARG SERVICE
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app ./cmd/${SERVICE}

FROM alpine:3.21
WORKDIR /app
COPY --from=build /out/app /app/app
EXPOSE 8080 50051 50052 50053 50054 50055
ENTRYPOINT ["/app/app"]
