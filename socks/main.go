package main

// FROM: https://echo.labstack.com/middleware/jaegertracing/#createchildspan

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var serviceName = "socks" // NOTE: typically use $OTEL_SERVICE_NAME for this

var tracer = otel.Tracer(serviceName)

func main() {
	log.Printf("Waiting for connection...")
	ctx := context.TODO()
	shutdown, err := initProvider()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	e := echo.New()
	e.Use(middleware.Logger())

	e.Use(otelecho.Middleware("unknown")) // TODO: not sure if this MW is working

	e.GET("/socks/*", getSocks)
	e.Logger.Info(e.Start(":1323"))
}

// getSocks is an example handler with a custom span
func getSocks(c echo.Context) error {
	ctx := context.TODO()
	defer ctx.Done()

	// custom Tracer
	_, span := tracer.Start(
		ctx,
		"Get-Socks", // set the name, this will be searchable later
		trace.WithAttributes(commonAttrs...))

	defer span.End()

	// if span.IsRecording() {}
	span.SetAttributes(
		attribute.String("http.method", "GET"),
		attribute.String("http.route", "/socks/:id"),
	)

	// randomise failure
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(10) // generate int between 0 and 9
	if i > 8 {         // 1 out of 10 times we fail
		randomSleep(100, 15)
		msg := "no socks for you"
		err := errors.New(msg)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return c.String(http.StatusBadGateway, msg)
	}

	// success scenario
	randomSleep(20, 10)
	msg := "you got socks"
	span.SetStatus(codes.Ok, msg)
	return c.String(http.StatusOK, msg)
}

// randomSleep will take some input numbers and sleep
func randomSleep(maxMS, minMS int) {
	// the minMS is just because APIs as fast as 1MS are rare in reality
	rand.Seed(time.Now().Local().UnixNano())
	n := rand.Intn(maxMS-minMS) + minMS
	time.Sleep(time.Duration(n) * time.Millisecond)
}
