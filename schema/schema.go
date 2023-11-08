package schema

import "time"

type AnyMap map[string]interface{}

type ManticoreResponseHitSource struct {
	ID                     int64   `json:"id"`
	Name                   string  `json:"name"`
	Slug                   string  `json:"slug"`
	Summary                string  `json:"summary,omitempty"`
	Storyline              string  `json:"storyline,omitempty"`
	Checksum               string  `json:"checksum,omitempty"`
	AggregatedRating       float64 `json:"aggregated_rating,omitempty"`
	AggregatedRatingCount  int64   `json:"aggregated_rating_count,omitempty"`
	Category               uint8   `json:"category"`
	Status                 uint8   `json:"status"`
	FirstReleaseDate       uint64  `json:"first_release_date,omitempty"`
	CreatedAt              uint64  `json:"created_at"`
	UpdatedAt              uint64  `json:"updated_at"`
	CoverUrl               string  `json:"cover_url,omitempty"`
	NumViews               uint64  `json:"num_views,omitempty"`
	NumLikes               uint64  `json:"num_likes,omitempty"`
	GenresNames            string  `json:"genres_names,omitempty"`
	PlatformsNames         string  `json:"platforms_names,omitempty"`
	PlatformsAbbreviations string  `json:"platforms_abbreviations,omitempty"`
	KeywordsNames          string  `json:"keywords_names,omitempty"`
	Source                 string  `json:"source"`
}

type SearchGame struct {
	ID                     int64     `json:"id"`
	Name                   string    `json:"name"`
	Slug                   string    `json:"slug"`
	Summary                string    `json:"summary,omitempty"`
	Storyline              string    `json:"storyline,omitempty"`
	Checksum               string    `json:"checksum,omitempty"`
	AggregatedRating       float64   `json:"aggregatedRating,omitempty"`
	AggregatedRatingCount  int64     `json:"aggregatedRatingCount,omitempty"`
	Category               uint8     `json:"category"`
	Status                 uint8     `json:"status"`
	FirstReleaseDate       time.Time `json:"firstReleaseDate,omitempty"`
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
	CoverUrl               string    `json:"coverUrl,omitempty"`
	NumViews               uint64    `json:"numViews,omitempty"`
	NumLikes               uint64    `json:"numLikes,omitempty"`
	GenresNames            string    `json:"genresNames,omitempty"`
	PlatformsNames         string    `json:"platformsNames,omitempty"`
	PlatformsAbbreviations string    `json:"platformsAbbreviations,omitempty"`
	KeywordsNames          string    `json:"keywordsNames,omitempty"`
	Source                 string    `json:"source"`
}

type ManticoreResponseHit struct {
	ID    string `json:"_id"`
	Score int    `json:"_score"`
	// To avoid re-mapping here, we only map the final SearchGame which is used in GameSearchResponseHit
	// This is basically SearchGame as returned by ManticoreSearch, in snake_case.
	Source ManticoreResponseHitSource `json:"_source"`
}

type ManticoreSearchHits struct {
	MaxScore      *float64               `json:"max_score,omitempty"`
	Total         *int                   `json:"total,omitempty"`
	TotalRelation *string                `json:"total_relation,omitempty"`
	Hits          []ManticoreResponseHit `json:"hits,omitempty"`
}

// ManticoreSearchResponse
// This struct defines the shape of the object returned by ManticoreSearch
// The Manticore naming convention is snake_case, which doesn't align with ours (camelCase) (or even Go's).
// So, before returning this to the client, make sure to convert it to a GameSearchResponseDto.
type ManticoreSearchResponse struct {
	Took         *int                    `json:"took,omitempty"`
	TimedOut     *bool                   `json:"timed_out,omitempty"`
	Aggregations map[string]interface{}  `json:"aggregations,omitempty"`
	Hits         *ManticoreSearchHits    `json:"hits,omitempty"`
	Profile      *map[string]interface{} `json:"profile,omitempty"`
	Warning      map[string]interface{}  `json:"warning,omitempty"`
}

type ManticoreErrorResponse struct {
	Total   int    `json:"total"`
	Warning string `json:"warning"`
	Error   string `json:"error"`
}

type GameSearchResponseHit struct {
	ID     uint       `json:"id"`
	Score  int        `json:"score"`
	Source SearchGame `json:"source"`
}

type GameSearchResponseHits struct {
	MaxScore      *float64                `json:"maxScore,omitempty"`
	Total         *int                    `json:"total,omitempty"`
	TotalRelation *string                 `json:"totalRelation,omitempty"`
	Hits          []GameSearchResponseHit `json:"hits,omitempty"`
}

type GameSearchResponseDto struct {
	Took         *int                    `json:"took,omitempty"`
	TimedOut     *bool                   `json:"timedOut,omitempty"`
	Aggregations map[string]interface{}  `json:"aggregations,omitempty"`
	Hits         *GameSearchResponseHits `json:"hits,omitempty"`
	Profile      *map[string]interface{} `json:"profile,omitempty"`
	Warning      map[string]interface{}  `json:"warning,omitempty"`
}

type GameSearchRequestDto struct {
	Index          string   `json:"index"`
	Query          AnyMap   `json:"query"`
	FulltextFilter AnyMap   `json:"fulltextFilter,omitempty"`
	AttrFilter     AnyMap   `json:"attrFilter,omitempty"`
	Limit          *int     `json:"limit,omitempty"`
	Offset         *int     `json:"offset,omitempty"`
	MaxMatches     *int     `json:"maxMatches,omitempty"`
	Sort           []AnyMap `json:"sort,omitempty"`
	Aggs           []AnyMap `json:"aggs,omitempty"`
	Expressions    []AnyMap `json:"expressions,omitempty"`
	Highlight      AnyMap   `json:"highlight,omitempty"`
	Source         AnyMap   `json:"source,omitempty"`
	Profile        *bool    `json:"profile,omitempty"`
	TrackScores    *bool    `json:"trackScores,omitempty"`
}
