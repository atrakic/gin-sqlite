FROM golang:1.19 AS builder
WORKDIR /src

RUN useradd -u 1001 nonroot

COPY go.mod ./

ENV GIN_MODE=release

# Use cache mounts to speed up the installation of existing dependencies
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . .

# Statically compile our app for use in a distroless container
#RUN CGO_ENABLED=0 go build -ldflags="-w -s" -v -o /bin/server .
RUN go build \
    #-ldflags="-w -s" \
    -ldflags="-linkmode external -extldflags -static" \
    -v -o /bin/server


# A distroless container image with some basics like SSL certificates
# https://github.com/GoogleContainerTools/distroless
#FROM gcr.io/distroless/static-debian12 AS final
FROM scratch
LABEL maintainer="Admir Trakic <atrakic@users.noreply.github.com>"
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /bin/server /
USER nonroot
EXPOSE 8080
CMD ["/server"]
