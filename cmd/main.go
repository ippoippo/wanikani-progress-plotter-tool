package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strconv"

	"github.com/ippoippo/wanikani-progress-plotter-tool/internal/wanikani"
	slogg "github.com/ippoippo/wanikani-progress-plotter-tool/pkg/slog"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const (
	apiKeyEnvVarName = "WANIKANI_API_KEY"
)

// Run me via
// `go run ./cmd/main.go `
func main() {
	ctx := context.Background()
	slog.InfoContext(ctx, "starting data extraction and plotting")

	apiKey := readKey(ctx) // This can cause program exit

	// Fetch data
	levels, err := wanikani.GetLevelProgressionData(ctx, apiKey)
	if err != nil {
		slogg.ErrorContextWithOSExit(ctx, "failed to fetch level progression data", err)
	}

	p := plot.New()
	hLabels := make([]string, 0, len(levels))
	values := make(plotter.Values, 0, len(levels))
	keys := make([]int, 0, len(levels))
	for k := range levels {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for key := range keys {
		hLabels = append(hLabels, strconv.Itoa(key+1))
		values = append(values, levels[key+1])
	}

	barChart, err := plotter.NewBarChart(values, 0.5*vg.Centimeter)
	if err != nil {
		slogg.ErrorContextWithOSExit(ctx, "error creating chart", err)
	}
	barChart.Horizontal = false
	p.Add(barChart)
	p.NominalX(hLabels...)

	// Save the plot to a PNG file
	if err := p.Save(10*vg.Inch, 6*vg.Inch, "wanikani_progress.png"); err != nil {
		slogg.ErrorContextWithOSExit(ctx, "error saving plot", err)
	}

	slog.InfoContext(ctx, "graph saved as wanikani_progress.png")
}

func readKey(ctx context.Context) string {
	val, ok := os.LookupEnv(apiKeyEnvVarName)
	if !ok {
		slogg.ErrorContextFromMsgWithOSExit(ctx, fmt.Sprintf("env var %s not set\n", apiKeyEnvVarName))

	}
	return val
}
