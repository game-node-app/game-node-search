package users

import "game-node-search/schema"

type UserManticoreSearchSource struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

type UserManticoreResponseHit struct {
	ID     uint                      `json:"_id"`
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
	Took    *int32                 `json:"took,omitempty"`
	Items   *[]UserDto             `json:"items,omitempty"`
	Profile map[string]interface{} `json:"profile,omitempty"`
}

type UserSearchResponseDto struct {
	Data       *UserSearchResponseData `json:"data"`
	Pagination *schema.PaginationInfo  `json:"pagination"`
}

type UserSearchRequestDto struct {
	Query string `json:"query"`
	Limit *int32 `json:"limit"`
	Page  *int32 `json:"page"`
}
