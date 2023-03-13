ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

# xcaddy is used to build caddy, we use this to allow access to build custom add-ons
# see:
#  - https://github.com/caddyserver/xcaddy
#  - https://hub.docker.com/_/caddy/tags?page=1&name=builder-alpine
FROM --platform=${BUILDPLATFORM:-linux/amd64} docker.io/caddy:builder-alpine AS builder

# Keep this version current!!
# see:
# - https://github.com/caddyserver/caddy/releases
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
      xcaddy build 9e943319b4ba2120e1f862f6615afb3c99e3a81e

### final image
FROM --platform=${BUILDPLATFORM:-linux/amd64} docker.io/alpine

RUN apk add --no-cache \
      ca-certificates

# add new user, lets not run as root
RUN adduser -D app
COPY --from=builder /usr/bin/caddy /caddy

# Default config + content
WORKDIR /app/content
COPY Caddyfile /app/Caddyfile
RUN chown -R app /app

# tell caddy where to listen
ENV SITE_ADDRESS=0.0.0.0:8080

# change user to "app"
USER app

CMD ["/caddy", "run", "--config", "/app/Caddyfile", "--adapter", "caddyfile"]

