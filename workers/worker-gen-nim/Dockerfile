FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/worker-gen-nim ./workers/worker-gen-nim/cmd

FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/worker-gen-nim /app
USER nonroot:nonroot
ENTRYPOINT ["/app"]
