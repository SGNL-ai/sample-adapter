ARG GOLANG_IMAGE=golang:1.23.4-bookworm

# STAGE 1: build...
FROM ${GOLANG_IMAGE} as build

WORKDIR /app
COPY . ./

RUN go mod download

ARG GOPS_VERSION=v0.3.27
RUN CGO_ENABLED=0 go install -ldflags "-s -w" github.com/google/gops@${GOPS_VERSION}
RUN CGO_ENABLED=0 GOOS=linux go build -C /app/cmd/adapter -o /sgnl/adapter

# STAGE 2: run...
FROM gcr.io/distroless/static AS run

# Fixture files are loaded from `pkg/mock/...`, but we copy the files to `/sgnl/pkg/mock/...`
# so we change the working directory to `/sgnl` to make sure the files are found.
WORKDIR /sgnl

COPY --from=build --chown=nonroot:nonroot /go/bin/gops /sgnl/gops
COPY --from=build --chown=nonroot:nonroot /sgnl/adapter /sgnl/adapter

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT [ "/sgnl/adapter" ]
