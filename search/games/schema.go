package games

import (
	"game-node-search/schema"
	"time"
)

type ManticoreResponseHitSource struct {
	ID                     uint64  `json:"id"`
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
	GenresNames            string  `json:"genres_names,omitempty"`
	PlatformsNames         string  `json:"platforms_names,omitempty"`
	PlatformsAbbreviations string  `json:"platforms_abbreviations,omitempty"`
	KeywordsNames          string  `json:"keywords_names,omitempty"`
	NumLikes               uint64  `json:"num_likes,omitempty"`
	NumViews               uint64  `json:"num_views,omitempty"`
	Source                 string  `json:"source"`
	ThemesNames            string  `json:"themes_names,omitempty"`
}

type ManticoreResponseHit struct {
	ID    uint64 `json:"_id"`
	Score int    `json:"_score"`
	// To avoid re-mapping here, we only map the final SearchGame which is used in GameSearchResponseHit
	// This is basically SearchGame as returned by ManticoreSearch, in snake_case.
	Source ManticoreResponseHitSource `json:"_source"`
}

type SearchGame struct {
	ID                     uint64    `json:"id"`
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
	GenresNames            string    `json:"genresNames,omitempty"`
	ThemesNames            string    `json:"themesNames,omitempty"`
	PlatformsNames         string    `json:"platformsNames,omitempty"`
	PlatformsAbbreviations string    `json:"platformsAbbreviations,omitempty"`
	KeywordsNames          string    `json:"keywordsNames,omitempty"`
	NumLikes               uint64    `json:"numLikes,omitempty"`
	NumViews               uint64    `json:"numViews,omitempty"`
	Source                 string    `json:"source"`
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

type ResponseData struct {
	Took    *int32                 `json:"took,omitempty"`
	Items   *[]SearchGame          `json:"items,omitempty"`
	Profile map[string]interface{} `json:"profile,omitempty"`
}

type GameSearchResponseDto struct {
	Data       ResponseData          `json:"data"`
	Pagination schema.PaginationInfo `json:"pagination,omitempty"`
}

type GameSearchRequestDto struct {
	Query     string   `json:"query"`
	Category  []int    `json:"category,omitempty"`
	Status    []int    `json:"status,omitempty"`
	Genres    []string `json:"genres,omitempty"`
	Themes    []string `json:"themes,omitempty"`
	Platforms []string `json:"platforms,omitempty"`
	Limit     *int32   `json:"limit,omitempty"`
	Page      *int32   `json:"page,omitempty"`
	Profile   *bool    `json:"profile,omitempty"`
}

type GameAutocompleteResponseDto struct {
	Total uint     `json:"total"`
	Data  []string `json:"data"`
}
