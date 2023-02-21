
# xcaddy is used to build caddy, we use this to allow access to build custom add-ons
# see:
#  - https://github.com/caddyserver/xcaddy
#  - https://hub.docker.com/_/caddy/tags?page=1&name=builder-alpine
FROM docker.io/caddy:builder-alpine AS builder

# Keep this version current!!
# see:
# - https://github.com/caddyserver/caddy/releases
RUN xcaddy build v2.6.4

FROM gruebel/upx:latest as upx
COPY --from=builder /usr/bin/caddy /caddy.orig

# Compress the binary ( https://upx.github.io )
RUN upx --best --lzma \ 
      /caddy.orig \
      -o /caddy


### final image
FROM alpine

RUN apk add --no-cache \
      ca-certificates

# add new user, lets not run as root
RUN adduser -D app
COPY --from=upx /caddy /caddy

# Default config + content
WORKDIR /app/content
COPY Caddyfile /app/Caddyfile
RUN chown -R app /app

# tell caddy where to listen
ENV SITE_ADDRESS=0.0.0.0:8080

# change user to "app"
USER app

CMD ["/caddy", "run", "--config", "/app/Caddyfile", "--adapter", "caddyfile"]

