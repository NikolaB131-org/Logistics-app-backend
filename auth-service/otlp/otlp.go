package otlp

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

var Tracer trace.Tracer

type Otlp struct {
	Tp *sdktrace.TracerProvider
}

func New() *Otlp {
	return &Otlp{}
}

func (o *Otlp) Init() {
	ctx := context.Background()
	insecureOpt := otlptracegrpc.WithInsecure()
	endpointOpt := otlptracegrpc.WithEndpoint(os.Getenv("OTLP_URL"))
	exp, err := otlptracegrpc.New(ctx, insecureOpt, endpointOpt)
	if err != nil {
		panic(err)
	}

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("authapp"),
		),
	)
	if err != nil {
		panic(err)
	}

	o.Tp = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)

	otel.SetTracerProvider(o.Tp)

	Tracer = o.Tp.Tracer("authapp")
}

func (o *Otlp) Shutdown() {
	o.Tp.Shutdown(context.Background())
}
