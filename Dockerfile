
# xcaddy is used to build caddy, we use this to allow access to build custom add-ons
# see:
#  - https://github.com/caddyserver/xcaddy
#  - https://hub.docker.com/_/caddy/tags?page=1&name=builder-alpine
FROM docker.io/caddy:builder-alpine AS builder

# Keep this version current!!
# see:
# - https://github.com/caddyserver/caddy/releases
RUN xcaddy build v2.6.4

RUN wget -O /index.html https://github.com/caddyserver/dist/raw/v2.6.4/welcome/index.html ;


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

WORKDIR /content
COPY --from=builder /index.html .

# tell caddy where to listen
ENV SITE_ADDRESS=0.0.0.0:8080

CMD ["/caddy", "run", "--config", "/Caddyfile", "--adapter", "caddyfile"]

