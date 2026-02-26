FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/api-gateway-service ./apps/api-gateway-service/cmd

FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/api-gateway-service /app
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app"]
