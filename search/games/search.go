package games

import (
	"context"
	"errors"
	"fmt"
	"game-node-search/schema"
	"game-node-search/util"
	jsoniter "github.com/json-iterator/go"
	Manticoresearch "github.com/manticoresoftware/manticoresearch-go"
	"log/slog"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func ValidateGameSearchRequest(dtoBytes []byte) (*GameSearchRequestDto, error) {
	var request = GameSearchRequestDto{}
	fmt.Println(string(dtoBytes))
	err := json.Unmarshal(dtoBytes, &request)
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func buildManticoreMatchString(dto *GameSearchRequestDto, request *Manticoresearch.SearchRequest) error {
	query := &dto.Query

	if query == nil {
		return errors.New("query parameter empty")
	}

	if !request.HasQuery() {
		request.Query = Manticoresearch.NewSearchQuery()
	}

	if !request.Query.HasBool() {
		request.Query.Bool = Manticoresearch.NewBoolFilter()
	}

	if !request.Query.Bool.HasMust() {
		request.Query.Bool.Must = make([]Manticoresearch.QueryFilter, 0)
	}

	filter := Manticoresearch.NewQueryFilter()

	matchObj := map[string]interface{}{
		"name,alternative_names": fmt.Sprintf("^%s", *query),
	}

	filter.SetMatch(matchObj)

	request.Query.Bool.Must = append(request.Query.Bool.Must, *filter)

	return nil

}

func buildManticoreFilterString(dto *GameSearchRequestDto, request *Manticoresearch.SearchRequest) error {
	if dto.Category != nil && len(dto.Category) > 0 {
		inObj := map[string]interface{}{
			"category": dto.Category,
		}

		categoryFilter := Manticoresearch.QueryFilter{
			In: inObj,
		}

		request.Query.Bool.Must = append(request.Query.Bool.Must, categoryFilter)
	}
	if dto.Status != nil && len(dto.Status) > 0 {
		inObj := map[string]interface{}{
			"status": dto.Status,
		}

		statusFilter := Manticoresearch.QueryFilter{
			In: inObj,
		}

		request.Query.Bool.Must = append(request.Query.Bool.Must, statusFilter)
	}

	return nil
}

func buildManticorePaginationString(dto *GameSearchRequestDto, request *Manticoresearch.SearchRequest) error {
	limit := dto.Limit
	if limit == nil || *limit == 0 {
		i := schema.DefaultLimit
		limit = &i
	}
	page := dto.Page
	if page == nil || *page == 0 {
		u := int32(1)
		page = &u
	}

	offset := (*page - 1) * *limit

	if offset < 0 {
		return errors.New("invalid resulting offset")
	}

	request.Limit = limit
	request.Offset = &offset

	return nil
}

func buildManticoreOrderString(request *Manticoresearch.SearchRequest) {
	request.Sort = map[string]interface{}{
		"_score":             "desc",
		"num_views":          "desc",
		"first_release_date": "desc",
	}
}

func buildManticoreSearchRequest(dto *GameSearchRequestDto) (Manticoresearch.SearchRequest, error) {
	searchRequest := Manticoresearch.NewSearchRequest("games")

	var err error = nil

	err = buildManticoreMatchString(dto, searchRequest)
	err = buildManticoreFilterString(dto, searchRequest)
	err = buildManticorePaginationString(dto, searchRequest)
	buildManticoreOrderString(searchRequest)

	if err != nil {
		slog.Error("Last error while building match filter: ", "err", err)
		return Manticoresearch.SearchRequest{}, err
	}

	return *searchRequest, err

}

// Search search handler
//
//	@Summary      Searches for games using Manticore engine
//	@Description  Returns a parsed search response from the Manticore engine
//	@Tags         search
//	@Accept       json
//	@Produce      json
//	@Param        query   body      games.GameSearchRequestDto  true  "Request body"
//	@Success      200  {object}   games.GameSearchResponseDto
//	@Router       /search/games [post]
func Search(dto *GameSearchRequestDto) (*GameSearchResponseDto, error) {

	request, err := buildManticoreSearchRequest(dto)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Generated request: %v\n", request)

	manticore := util.GetManticoreInstance()

	mr, r, err := manticore.SearchAPI.Search(context.Background()).SearchRequest(request).Execute()

	if err != nil {
		slog.Error("error while calling Manticore instance: ", "err", err, "response", r)
		return nil, err
	}

	data, err := buildResponseData(mr)

	pagination := buildPaginationInfo(mr, dto.Limit, dto.Page)

	response := GameSearchResponseDto{
		Data:       *data,
		Pagination: *pagination,
	}

	return &response, nil
}
