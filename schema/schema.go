package schema

import "time"

type AnyMap map[string]interface{}

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

type ManticoreSearchRequestQueryBool struct {
	Match   *ManticoreSearchRequestQueryMatch `json:"match"`
	Must    *AnyMap                           `json:"must"`
	MustNot *AnyMap                           `json:"must_not"`
	Should  *AnyMap                           `json:"should"`
}

type ManticoreSearchRequestQueryMatch struct {
	All *AnyMap `json:"_all"`
}

type ManticoreSearchRequestQuery struct {
	Match ManticoreSearchRequestQueryMatch `json:"match"`
	Bool  *ManticoreSearchRequestQueryBool `json:"bool"`
}

type ManticoreSearchRequest struct {
	Index          string                      `json:"index"`
	Query          ManticoreSearchRequestQuery `json:"query"`
	FulltextFilter AnyMap                      `json:"fulltext_filter,omitempty"`
	AttrFilter     AnyMap                      `json:"attr_filter,omitempty"`
	Limit          *int                        `json:"limit,omitempty"`
	Offset         *int                        `json:"offset,omitempty"`
	MaxMatches     *int                        `json:"max_matches,omitempty"`
	Sort           []AnyMap                    `json:"sort,omitempty"`
	Aggs           []AnyMap                    `json:"aggs,omitempty"`
	Expressions    []AnyMap                    `json:"expressions,omitempty"`
	Highlight      AnyMap                      `json:"highlight,omitempty"`
	Source         AnyMap                      `json:"source,omitempty"`
	Profile        *bool                       `json:"profile,omitempty"`
	TrackScores    *bool                       `json:"trackScores,omitempty"`
}

type ManticoreErrorResponse struct {
	Total   *int    `json:"total"`
	Warning *string `json:"warning"`
	Error   *string `json:"error"`
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

type PaginationInfo struct {
	TotalItems  uint64 `json:"totalItems"`
	TotalPages  uint16 `json:"totalPages"`
	HasNextPage bool   `json:"hasNextPage"`
}
type ResponseData struct {
	Took    *int                    `json:"took,omitempty"`
	Items   *[]SearchGame           `json:"items,omitempty"`
	Profile *map[string]interface{} `json:"profile,omitempty"`
}

type GameSearchResponseDto struct {
	Data       ResponseData   `json:"data"`
	Pagination PaginationInfo `json:"pagination,omitempty"`
}

type GameSearchRequestDto struct {
	Query     string    `json:"query"`
	Category  *[]int    `json:"category,omitempty"`
	Status    *[]int    `json:"status,omitempty"`
	Limit     *int      `json:"limit,omitempty"`
	Page      *int      `json:"page,omitempty"`
	Profile   *bool     `json:"profile,omitempty"`
}

type UserManticoreSearchSource struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

type UserManticoreResponseHit struct {
	ID     string                    `json:"_id"`
	Score  int                       `json:"_score"`
	Source UserManticoreSearchSource `json:"_source"`
}

type UserManticoreSearchHits struct {
	MaxScore      *float64                   `json:"max_score,omitempty"`
	Total         *int                       `json:"total,omitempty"`
	TotalRelation *string                    `json:"total_relation,omitempty"`
	Hits          []UserManticoreResponseHit `json:"hits,omitempty"`
}

type UserManticoreSearchResponse struct {
	Took         *int                     `json:"took,omitempty"`
	TimedOut     *bool                    `json:"timed_out,omitempty"`
	Aggregations map[string]interface{}   `json:"aggregations,omitempty"`
	Hits         *UserManticoreSearchHits `json:"hits,omitempty"`
	Profile      *map[string]interface{}  `json:"profile,omitempty"`
	Warning      map[string]interface{}   `json:"warning,omitempty"`
}

type UserDto struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
}

type UserSearchResponseData struct {
	Took    *int                    `json:"took,omitempty"`
	Items   *[]UserDto              `json:"items,omitempty"`
	Profile *map[string]interface{} `json:"profile,omitempty"`
}

type UserSearchResponseDto struct {
	Data       *UserSearchResponseData `json:"data"`
	Pagination *PaginationInfo         `json:"pagination"`
}

type UserSearchRequestDto struct {
	Query string `json:"query"`
	Limit *int   `json:"limit"`
	Page  *int   `json:"page"`
}

const DEFAULT_LIMIT int = 20
