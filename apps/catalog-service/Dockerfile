FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/catalog-service ./apps/catalog-service/cmd

FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/catalog-service /app
EXPOSE 8083
USER nonroot:nonroot
ENTRYPOINT ["/app"]
