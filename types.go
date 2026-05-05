package fred

import (
	"encoding/json"
	"fmt"
	"time"
)

// Observation represents a single FRED data observation.
type Observation struct {
	Date          time.Time
	Value         float64
	RealtimeStart time.Time
	RealtimeEnd   time.Time
	IsNA          bool
}

// rawObservation is the JSON wire format (all fields are strings from FRED).
type rawObservation struct {
	Date          string `json:"date"`
	Value         string `json:"value"`
	RealtimeStart string `json:"realtime_start"`
	RealtimeEnd   string `json:"realtime_end"`
}

// observationResponse wraps the FRED observations JSON response.
type observationResponse struct {
	RealtimeStart    string           `json:"realtime_start"`
	RealtimeEnd      string           `json:"realtime_end"`
	ObservationStart string           `json:"observation_start"`
	ObservationEnd   string           `json:"observation_end"`
	Units            string           `json:"units"`
	OutputType       int              `json:"output_type"`
	FileType         string           `json:"file_type"`
	OrderBy          string           `json:"order_by"`
	SortOrder        string           `json:"sort_order"`
	Count            int              `json:"count"`
	Offset           int              `json:"offset"`
	Limit            int              `json:"limit"`
	Observations     []rawObservation `json:"observations"`
}

// seriesInfoResponse wraps the FRED series metadata JSON response.
// FRED uses "seriess" (double 's') — not a typo.
type seriesInfoResponse struct {
	Series []Series `json:"seriess"`
}

// Series represents metadata for a FRED series.
type Series struct {
	ID                      string `json:"id"`
	RealtimeStart           string `json:"realtime_start"`
	RealtimeEnd             string `json:"realtime_end"`
	Title                   string `json:"title"`
	ObservationStart        string `json:"observation_start"`
	ObservationEnd          string `json:"observation_end"`
	Frequency               string `json:"frequency"`
	FrequencyShort          string `json:"frequency_short"`
	Units                   string `json:"units"`
	UnitsShort              string `json:"units_short"`
	SeasonalAdjustment      string `json:"seasonal_adjustment"`
	SeasonalAdjustmentShort string `json:"seasonal_adjustment_short"`
	LastUpdated             string `json:"last_updated"`
	Popularity              int    `json:"popularity"`
	Notes                   string `json:"notes,omitempty"`
	Tags                    []Tag  `json:"tags,omitempty"`
}

// SearchResult represents a series search result.
type SearchResult struct {
	ID                      string `json:"id"`
	RealtimeStart           string `json:"realtime_start"`
	RealtimeEnd             string `json:"realtime_end"`
	Title                   string `json:"title"`
	ObservationStart        string `json:"observation_start"`
	ObservationEnd          string `json:"observation_end"`
	Frequency               string `json:"frequency"`
	FrequencyShort          string `json:"frequency_short"`
	Units                   string `json:"units"`
	UnitsShort              string `json:"units_short"`
	SeasonalAdjustment      string `json:"seasonal_adjustment"`
	SeasonalAdjustmentShort string `json:"seasonal_adjustment_short"`
	LastUpdated             string `json:"last_updated"`
	Popularity              int    `json:"popularity"`
	Notes                   string `json:"notes,omitempty"`
}

// searchResponse wraps the FRED search JSON response.
type searchResponse struct {
	Count  int            `json:"count"`
	Series []SearchResult `json:"seriess"`
}

// Category represents a FRED category.
type Category struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ParentID int    `json:"parent_id"`
	Notes    string `json:"notes,omitempty"`
}

// categoriesResponse wraps the FRED categories JSON response.
type categoriesResponse struct {
	Categories []Category `json:"categories"`
}

// Release represents a FRED release.
type Release struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	PressRelease bool   `json:"press_release"`
	Link         string `json:"link,omitempty"`
	Notes        string `json:"notes,omitempty"`
}

// releasesResponse wraps the FRED releases JSON response.
type releasesResponse struct {
	Releases []Release `json:"releases"`
}

// ReleaseDate represents a release date entry.
type ReleaseDate struct {
	ReleaseID   int    `json:"release_id"`
	ReleaseName string `json:"release_name,omitempty"`
	Date        string `json:"date"`
}

// releasesDatesResponse wraps the FRED release dates JSON response.
type releasesDatesResponse struct {
	Dates []ReleaseDate `json:"release_dates"`
}

// Source represents a FRED data source.
type Source struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Link  string `json:"link,omitempty"`
	Notes string `json:"notes,omitempty"`
}

// sourcesResponse wraps the FRED sources JSON response.
type sourcesResponse struct {
	Sources []Source `json:"sources"`
}

// Tag represents a FRED tag.
type Tag struct {
	Name        string `json:"name"`
	GroupID     string `json:"group_id"`
	Notes       string `json:"notes,omitempty"`
	Created     string `json:"created"`
	Popularity  int    `json:"popularity"`
	SeriesCount int    `json:"series_count"`
}

// tagsResponse wraps the FRED tags JSON response.
type tagsResponse struct {
	Tags []Tag `json:"tags"`
}

// vintageDatesResponse wraps the FRED vintage dates JSON response.
type vintageDatesResponse struct {
	Dates []string `json:"vintage_dates"`
}

// ReleaseTableElement represents an element in a release table.
type ReleaseTableElement struct {
	ElementID int                   `json:"element_id"`
	ReleaseID int                   `json:"release_id"`
	SeriesID  string                `json:"series_id,omitempty"`
	ParentID  int                   `json:"parent_id"`
	Line      string                `json:"line"`
	Type      string                `json:"type"`
	Name      string                `json:"name"`
	Level     int                   `json:"-"` // computed from string
	LevelRaw  string                `json:"level"`
	Children  []ReleaseTableElement `json:"children,omitempty"`
}

// releaseTableResponse wraps the FRED release table JSON response.
type releaseTableResponse struct {
	Elements map[string]ReleaseTableElement `json:"elements"`
}

// SeriesGroup represents metadata for a GeoFRED map series.
type SeriesGroup struct {
	Title       string `json:"title"`
	RegionType  string `json:"region_type"`
	SeriesGroup string `json:"series_group"`
	Season      string `json:"season"`
	Units       string `json:"units"`
	Frequency   string `json:"frequency"`
	MinDate     string `json:"min_date"`
	MaxDate     string `json:"max_date"`
}

// MapDataEntry represents a single region's data point in GeoFRED.
type MapDataEntry struct {
	Region   string `json:"region"`
	Code     string `json:"code"`
	Value    string `json:"value"`
	SeriesID string `json:"series_id"`
}

// UnmarshalJSON handles both numeric and string value representations from FRED.
func (m *MapDataEntry) UnmarshalJSON(data []byte) error {
	type alias MapDataEntry
	aux := &struct {
		Value interface{} `json:"value"`
		*alias
	}{
		alias: (*alias)(m),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	switch v := aux.Value.(type) {
	case float64:
		if aux.Value == nil {
			m.Value = "."
		} else {
			m.Value = fmt.Sprintf("%v", v)
		}
	case string:
		m.Value = v
	}
	return nil
}

// MapData wraps GeoFRED series/regional data responses.
type MapData struct {
	Title       string `json:"title"`
	Region      string `json:"region"`
	Seasonality string `json:"seasonality"`
	Units       string `json:"units"`
	Frequency   string `json:"frequency"`
	Date        string `json:"date"`
	Data        []MapDataEntry
}

// mapDataResponse is the raw JSON form where data is map[string][]MapDataEntry.
type mapDataResponse struct {
	Title       string                    `json:"title"`
	Region      string                    `json:"region"`
	Seasonality string                    `json:"seasonality"`
	Units       string                    `json:"units"`
	Frequency   string                    `json:"frequency"`
	Date        string                    `json:"date"`
	Data        map[string][]MapDataEntry `json:"data"`
}
