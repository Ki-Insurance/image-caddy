
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

# Compress the binary
RUN upx --best --lzma \ 
      /caddy.orig \
      -o /caddy





### final image
FROM alpine

RUN apk add --no-cache \
      ca-certificates

# add new user
RUN adduser -D app

COPY --from=upx /caddy /caddy

RUN ls -lash /

WORKDIR /app

CMD [ "/caddy" ]
