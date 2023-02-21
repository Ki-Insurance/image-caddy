package main

// FROM: https://echo.labstack.com/middleware/jaegertracing/#createchildspan

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit/v6"
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

// Info
type Info struct {
	Time  time.Duration `json:"time"`
	Count int           `json:"count"`
}

// Sock
type Sock struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Time    time.Duration `json:"time"`
}

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

	e.GET("/socks/info", getInfo)
	e.GET("/socks/*", getSocks)
	e.Logger.Info(e.Start(":1323"))
}

// getInfo gets some info to the UI
func getInfo(c echo.Context) error {
	start := time.Now()
	ctx := context.TODO()
	defer ctx.Done()

	// custom Tracer
	_, span := tracer.Start(
		ctx,
		"Get-Info", // set the name, this will be searchable later
		trace.WithAttributes(commonAttrs...))
	defer span.End()

	if span.IsRecording() {
		span.SetAttributes(
			attribute.String("http.method", "GET"),
			attribute.String("http.route", "/socks/info"),
		)
	}

	// success scenario
	randomSleep(20, 10)
	info := &Info{
		Time: time.Since(start),
	}
	span.SetStatus(codes.Ok, fmt.Sprintln(info))
	return c.JSON(http.StatusOK, info)
}

// getSocks is an example handler with a custom span
func getSocks(c echo.Context) error {
	start := time.Now()
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

	sock := &Sock{}
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(10) // generate int between 0 and 9
	if i > 8 {         // 1 out of 10 times we fail
		sock.Success = false
		msg := "no socks found"
		sock.Message = msg
		randomSleep(100, 15)
		err := errors.New(msg)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		sock.Time = time.Since(start)
		return c.JSON(http.StatusBadGateway, sock)
	}

	// success scenario
	randomSleep(20, 10)
	msg := fmt.Sprintf("%s socks", gofakeit.Color())
	span.SetStatus(codes.Ok, msg)
	sock.Message = msg
	sock.Success = true
	sock.Time = time.Since(start)
	return c.JSON(http.StatusOK, sock)
}

// randomSleep will take some input numbers and sleep
func randomSleep(maxMS, minMS int) {
	// the minMS is just because APIs as fast as 1MS are rare in reality
	rand.Seed(time.Now().Local().UnixNano())
	n := rand.Intn(maxMS-minMS) + minMS
	time.Sleep(time.Duration(n) * time.Millisecond)
}
