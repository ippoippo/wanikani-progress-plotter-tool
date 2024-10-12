package wanikani

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/ippoippo/wanikani-progress-plotter-tool/internal/wanikani/levelprogressions"
	slogg "github.com/ippoippo/wanikani-progress-plotter-tool/pkg/slog"
)

type Levels map[int]float64 // Key=level, Value=days taken

func GetLevelProgressionData(ctx context.Context, apiKey string) (Levels, error) {
	rawData, err := getRawDataFromApi(ctx, apiKey)
	if err != nil {
		return nil, err
	}
	return mapRawDataToLevels(rawData), nil
}

func getRawDataFromApi(ctx context.Context, apiKey string) (*levelprogressions.ResponseBody, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.wanikani.com/v2/level_progressions", nil)
	if err != nil {
		slog.ErrorContext(ctx, "http.NewRequestWithContext has returned error: ", slogg.ErrorAttr(err))
		return nil, err // We don't expect this to happen
	}

	// Set the Authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "httpClient.Do(req) has returned error: ", slogg.ErrorAttr(err), slogg.CtxErrorAttr(ctx))
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.ErrorContext(ctx, "bodyBytes could not be read",
			slogg.ErrorAttr(err), slog.Int("http-status-code", resp.StatusCode))
		return nil, err
	}

	result, err := levelprogressions.NewResponseFrom(bodyBytes)
	if err != nil {
		slog.ErrorContext(ctx, "failed to parse result", slogg.ErrorAttr(err),
			slog.Int("http-status-code", resp.StatusCode),
			slog.String("raw-body-bytes", string(bodyBytes)))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		slog.ErrorContext(ctx, "request failed with non-200 Status Code", slog.Int("status-code", resp.StatusCode),
			slog.String("raw-body-bytes", string(bodyBytes)), slog.Any("result", result))
		return nil, err
	}

	return result, nil
}

func mapRawDataToLevels(raw *levelprogressions.ResponseBody) Levels {
	levels := Levels{}

	for _, lp := range raw.Data {
		if lpd := lp.Data; lpd != nil {
			if lpd.PassedAt.IsZero() {
				levels[lpd.Level] = 0
			} else {
				levels[lpd.Level] = lpd.PassedAt.Sub(lpd.UnlockedAt).Hours() / 24
			}
		}
	}

	return levels
}
