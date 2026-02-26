FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/progress-service ./apps/progress-service/cmd

FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/progress-service /app
EXPOSE 8086
USER nonroot:nonroot
ENTRYPOINT ["/app"]
