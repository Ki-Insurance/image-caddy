# What is this

This repo builds an opinionated docker image with a focus on:
- [x] improved security (not running as root)
- [x] prometheus metrics
- [x] opentelemetry tracing
- [x] structured (json) logs
- [x] log obfuscation abilities
- [x] the webserver is specifialised for serving PWA/SPA (static) content
- [x] support for (optional) reverse proxying


# Usage

### 1. Add contents
Add static assets such as your `index.html`, CSS and images  into the `/app/content` folder

### 2. (Optionally) adjust Caddy's configuration
Out of the box there is a very simple [Caddyfile](https://github.com/Ki-Insurance/image-caddy/blob/main/Caddyfile) already installed to `/app/Caddyfile`.
If you need to adjust settings, override this file.
> NOTE: The env var `$SITE_ADDRESS` controls which `hostname:port` caddy will listen to.
>> Its default is: `:8080`
>
> This means that it will respond to any hostname on TCP port `8080`


Example Dockerfile:

```
FROM ghcr.io/ki-insurance/image-caddy

# Override the default Caddyfile (optional)
COPY Caddyfile /app/Caddyfile

# copy in some content (where "build" is a folder containing your assets and "index.html")
COPY build/index.html /app/content
```

> NOTE: this docker image does not run as root. So do not attempt to listen to TCP ports under `1024`

# OpenTelemetry Demo

There is a sample backend and `docker-compose.yaml` you can use to see how things work in OpenTelemetry.

To test it out (in the `demo` directory):

`docker-compose up --build`

After its launched, you'll want to trigger some traffic. You can do this by visiting: http://localhost:8080

> alternatively you can trigger something via some curl
>
> eg: `curl localhost:8080/socks/echo`

1 out of 10 of the curls should result in a non-200 (and they should be slower also)

Check out the following links

Link | Service | Usage
---|---|---
http://localhost:8080 | [sample SPA](https://github.com/Ki-Insurance/image-caddy/tree/main/demo/content) | Simple UI to trigger sample traffic
http://localhost:9411 | zipkin | Visualise OpenTelemetry traces
http://localhost:16686 | Jaegar | Visualise OpenTelemetry traces

# TODO
- [x] integrate and test opentelemetry tracing with caddy
- [x] integrate and test opentelemetry tracing with a demo backend
- [ ] test `http/json` + `b3` `baggage` against OTEL trace collector
- [ ] integrate and test opentelemetry metrics (prometheus) with caddy
- [ ] improved documentation (Sockshop etc)