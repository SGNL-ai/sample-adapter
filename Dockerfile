ARG GOLANG_IMAGE=golang:1.23.4-bookworm

# STAGE 1: build...
FROM ${GOLANG_IMAGE} as build

WORKDIR /app
COPY . ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -C /app/cmd/adapter -o /sgnl/adapter

# STAGE 2: run...
FROM gcr.io/distroless/static AS run

COPY --from=build --chown=nonroot:nonroot /sgnl/adapter /sgnl/adapter

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT [ "/sgnl/adapter" ]