package lobsters_test

import (
	"context"
	"testing"

	"github.com/notarock/technews-bot/pkg/sources/lobsters"
	"github.com/notarock/technews-bot/pkg/telemetry"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestFetchLatest(t *testing.T) {
	ctx := context.Background()

	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		t.Skip("Skipping test: no OTLP endpoint available")
	}

	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	otel.SetTracerProvider(tp)
	defer tp.Shutdown(ctx)

	telemetry.Init(ctx, "test", "localhost:4317")

	articles := lobsters.FetchLatestArticles(ctx)

	assert.Equal(t, len(articles), 25)

	// Articles results are time dependant.
	// Lets just check we got something.
	assert.NotEmpty(t, articles[0].Link)
	assert.NotEmpty(t, articles[0].Tags)
}
