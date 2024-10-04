FROM golang:1.22.6 AS build-stage

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/music_service/ ./cmd/music_service/
COPY ./internal/ ./internal/
COPY ./config/ ./config/
COPY ./pkg/ ./pkg/
COPY ./.env ./
COPY ./Makefile ./

RUN OUTPUT_DIR=/app make buildgo

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go generate ./...
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/music_service /music_service

EXPOSE 8080

# USER nonroot:nonroot

ENTRYPOINT ["/music_service"]