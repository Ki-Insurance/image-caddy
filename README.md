# What is this

This repo builds an opinionated docker image with a focus on:
- improved security (not running as root)
- prometheus metrics
- opentelemetry tracing
- structured (json) logs
- log obfuscation abilities
- the webserver is specifialised for serving PWA/SPA (static) content


# Usage

One simply needs to:
1. Override the default `/Caddyfile` for caddy configuration
2. Add contents (`index.html`, assets etc..) into the `/content` folder

Example Dockerfile:

```
FROM ghcr.io/ki-insurance/image-caddy:667f7c9

# Override the default Caddyfile (optional)
COPY Caddyfile /Caddyfile

## copy in some content
COPY build/index.html .
```


# OpenTelemetry Demo

There is a sample backend and `docker-compose.yaml` you can use to see how things work in OpenTelemetry.

To test it out:

`docker-compose up --build`

After its launched, trigger some traffic

```sh
while true; do
  curl localhost:8080/socks/foo -s
  sleep 0.5
  echo
done
```

1 out of 10 of the curls should result in a non-200 (and they should be slower also)

Now open Zipkin or Jaegar and you should see some traces
see:
- http://127.0.0.1:9411 (zipkin)
- http://127.0.0.1:16686/search
