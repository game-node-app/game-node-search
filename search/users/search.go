package users

import (
	"context"
	"fmt"
	"game-node-search/schema"
	"game-node-search/util"
	jsoniter "github.com/json-iterator/go"
	Manticoresearch "github.com/manticoresoftware/manticoresearch-go"
	"log/slog"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func ValidateUserSearchRequest(dtoBytes []byte) (*UserSearchRequestDto, error) {
	var request = UserSearchRequestDto{}
	fmt.Println(string(dtoBytes))
	err := json.Unmarshal(dtoBytes, &request)
	if err != nil {
		return nil, err
	}
	if &request.Query == nil || len(request.Query) == 0 {
		return nil, fmt.Errorf("invalid search parameters")
	}

	return &request, nil
}

func buildManticoreMatchString(dto *UserSearchRequestDto, request *Manticoresearch.SearchRequest) {
	matchObj := map[string]interface{}{
		"username": fmt.Sprintf("%s*", dto.Query),
	}

	if !request.HasQuery() {
		request.Query = Manticoresearch.NewSearchQuery()
	}

	request.Query.SetMatch(matchObj)
}

func buildManticorePaginationString(dto *UserSearchRequestDto, request *Manticoresearch.SearchRequest) {
	limit := dto.Limit
	page := dto.Page

	if limit == nil || int32(*limit) == 0 {
		i := schema.DefaultLimit
		limit = &i
	}

	if page == nil || int32(*page) == 0 {
		i := int32(1)
		page = &i
	}

	offset := (*page - 1) * *limit

	request.Limit = limit
	request.Offset = &offset
}

func buildManticoreSearchRequest(dto *UserSearchRequestDto) (Manticoresearch.SearchRequest, error) {
	searchRequest := Manticoresearch.NewSearchRequest("users")

	buildManticoreMatchString(dto, searchRequest)
	buildManticorePaginationString(dto, searchRequest)

	return *searchRequest, nil
}

// Search users handler
//
//	@Summary      Searches for users using Manticore engine
//	@Description  Returns a parsed search response from the Manticore engine
//	@Tags         search
//	@Accept       json
//	@Produce      json
//	@Param        query   body      schema.UserSearchRequestDto  true  "Account ID"
//	@Success      200  {object}   schema.UserSearchResponseDto
//	@Router       /search/users [post]
func Search(dto *UserSearchRequestDto) (*UserSearchResponseDto, error) {
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

	if err != nil {
		return nil, err
	}

	pagination := buildPaginationInfo(mr, dto.Limit)

	response := UserSearchResponseDto{
		Data:       data,
		Pagination: pagination,
	}

	return &response, nil
}
