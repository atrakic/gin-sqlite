FROM golang:1.19 AS builder
WORKDIR /src
COPY . /src
RUN go get -d -v
# Statically compile our app for use in a distroless container
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -v -o /bin/server .

# A distroless container image with some basics like SSL certificates
# https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/static AS final
LABEL maintainer="Admir Trakic <atrakic@users.noreply.github.com>"
WORKDIR /app
COPY --from=builder /bin/server ./
#HEALTHCHECK --start-period=1s --interval=10s --timeout=5s CMD [ "executable" ]
ENTRYPOINT ["./server"]
