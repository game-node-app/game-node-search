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

	matchObj := map[string]interface{}{
		"name,alternative_names": *query,
	}

	request.Query.SetMatch(matchObj)

	return nil

}

func buildManticoreFilterString(dto *GameSearchRequestDto) (string, error) {
	var filterString = ""
	if dto.Category != nil && len(*dto.Category) > 0 {
		var categoryFilterArrayString = ""
		for i, v := range *dto.Category {
			if i > 0 {
				categoryFilterArrayString = fmt.Sprintf("%s,%d", categoryFilterArrayString, v)
				continue
			}
			categoryFilterArrayString = fmt.Sprintf("%d", v)

		}

		filterString += fmt.Sprintf(" AND category IN (%s)", categoryFilterArrayString)
	}
	if dto.Status != nil && len(*dto.Status) > 0 {
		var statusFilterArrayString = ""
		for i, v := range *dto.Status {
			if i > 0 {
				statusFilterArrayString = fmt.Sprintf("%s,%d", statusFilterArrayString, v)
				continue
			}
			statusFilterArrayString = fmt.Sprintf("%d", v)
		}

		filterString += fmt.Sprintf(" AND status IN (%s)", statusFilterArrayString)

	}

	return filterString, nil
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

func buildManticoreOrderString() string {
	return "ORDER BY num_likes DESC, num_views DESC"
}

func buildManticoreSearchRequest(dto *GameSearchRequestDto) (Manticoresearch.SearchRequest, error) {
	searchRequest := Manticoresearch.NewSearchRequest("games")

	options := map[string]interface{}{
		"fuzzy": true,
	}

	searchRequest.SetOptions(options)

	var err error = nil

	err = buildManticoreMatchString(dto, searchRequest)
	err = buildManticorePaginationString(dto, searchRequest)

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
//	@Param        query   body      schema.GameSearchRequestDto  true  "Account ID"
//	@Success      200  {object}   schema.GameSearchResponseDto
//	@Router       /search/games [post]
func Search(dto *GameSearchRequestDto) (*GameSearchResponseDto, error) {

	request, err := buildManticoreSearchRequest(dto)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Generated request: %v\n", request)

	manticore := util.GetManticoreInstance()

	mr, _, err := manticore.SearchAPI.Search(context.Background()).SearchRequest(request).Execute()

	if err != nil {
		slog.Error("error while calling Manticore instance: ", "err", err)
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
