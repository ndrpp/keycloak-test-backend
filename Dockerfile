#Build the application from source
FROM golang:latest AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /go-server

# Run tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

#Deploy the app binary into a lean image
FROM gcr.io/distroless/base-debian12 AS build-release-stage

WORKDIR /

COPY --from=build-stage /go-server /go-server

EXPOSE 8081

USER nonroot:nonroot

ENTRYPOINT ["/go-server"]
