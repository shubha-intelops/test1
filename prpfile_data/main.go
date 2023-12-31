package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	restcontrollers "github.com/shubha-intelops/test1/prpfile_data/pkg/rest/server/controllers"
	"github.com/sinhashubham95/go-actuator"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
	"os"

	"github.com/shubha-intelops/test1/prpfile_data/pkg/rest/client"
)

func main() {

	router := gin.Default()
	if len(serviceName) > 0 && len(collectorURL) > 0 {
		// add opentel
		cleanup := initTracer()
		defer func(func(context.Context) error) {
			_ = cleanup(context.Background())
		}(cleanup)
		router.Use(otelgin.Middleware(serviceName))
	}

	// add actuator
	addActuator(router)
	// add prometheus
	addPrometheus(router)

	profileController, err := restcontrollers.NewProfileController()
	if err != nil {
		log.Errorf("error occurred: %s", err)
		return
	}

	v1 := router.Group("/v1")
	{

		v1.GET("/profiles/:id", profileController.FetchProfile)
		v1.POST("/profiles", profileController.CreateProfile)
		v1.PUT("/profiles/:id", profileController.UpdateProfile)
		v1.DELETE("/profiles/:id", profileController.DeleteProfile)
		v1.GET("/profiles", profileController.ListProfiles)
		v1.PATCH("/profiles/:id", profileController.PatchProfile)
		v1.HEAD("/profiles", profileController.HeadProfile)
		v1.OPTIONS("/profiles", profileController.OptionsProfile)

	}

	Port := ":8086"
	log.Println("Server started")
	if err = router.Run(Port); err != nil {
		log.Errorf("error occurred: %s", err)
		return
	}

	// call external services here if the HasRestClients value is true
	// (that means this repo is a client to external service(s)
	var err0 error

	bNode07, err0 := client.ExecuteNode07()
	if err0 != nil {
		log.Printf("error occurred: %s", err0)
		return
	}
	log.Printf("response received: %s", string(bNode07))

}

var (
	serviceName  = os.Getenv("SERVICE_NAME")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	insecure     = os.Getenv("INSECURE_MODE")
)

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func addPrometheus(router *gin.Engine) {
	router.GET("/metrics", prometheusHandler())
}

func addActuator(router *gin.Engine) {
	actuatorHandler := actuator.GetActuatorHandler(&actuator.Config{Endpoints: []int{
		actuator.Env,
		actuator.Info,
		actuator.Metrics,
		actuator.Ping,
		// actuator.Shutdown,
		actuator.ThreadDump,
	},
		Env:     "dev",
		Name:    "prpfile_data",
		Port:    8086,
		Version: "0.0.1",
	})
	ginActuatorHandler := func(ctx *gin.Context) {
		actuatorHandler(ctx.Writer, ctx.Request)
	}
	router.GET("/actuator/*endpoint", ginActuatorHandler)
}

func initTracer() func(context.Context) error {
	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if len(insecure) > 0 {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(collectorURL),
		),
	)

	if err != nil {
		log.Errorf("error occurred: %s", err)
		return nil
	}
	restResources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("services.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Printf("could not set restResources: %s", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(restResources),
		),
	)
	return exporter.Shutdown
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}
