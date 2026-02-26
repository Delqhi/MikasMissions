FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/billing-service ./apps/billing-service/cmd

FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/billing-service /app
EXPOSE 8089
USER nonroot:nonroot
ENTRYPOINT ["/app"]
