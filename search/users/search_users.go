package users

import (
	"errors"
	"fmt"
	"game-node-search/schema"
	"game-node-search/util"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	url2 "net/url"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func ValidateUserSearchRequest(dtoBytes []byte) (*schema.UserSearchRequestDto, error) {
	var request = schema.UserSearchRequestDto{}
	fmt.Println(string(dtoBytes))
	err := json.Unmarshal(dtoBytes, &request)
	if err != nil {
		return nil, err
	}
	if &request.Query == nil || len(request.Query) == 0 {
		return nil, fmt.Errorf("Invalid search parameters.")
	}

	return &request, nil
}

func buildManticoreMatchString(dto *schema.UserSearchRequestDto) string {
	return fmt.Sprintf("@username %s", dto.Query)
}

func buildManticorePaginationString(dto *schema.UserSearchRequestDto) string {
	limit := dto.Limit
	page := dto.Page
	if limit != nil || int(*limit) == 0 {
		i := schema.DEFAULT_LIMIT
		limit = &i
	}

	if page != nil || int(*page) == 0 {
		i := 1
		page = &i
	}

	offset := (*page - 1) * *limit

	return fmt.Sprintf("LIMIT %d OFFSET %d", *limit, offset)
}

func buildManticoreSearchRequest(dto *schema.UserSearchRequestDto) (string, error) {
	matchString := buildManticoreMatchString(dto)
	paginationString := buildManticorePaginationString(dto)

	selectString := fmt.Sprintf("SELECT * FROM users WHERE MATCH(%s) %s", matchString, paginationString)

	return selectString, nil
}

func UserSearchHandler(dto *schema.UserSearchRequestDto) (*schema.UserSearchResponseDto, error) {
	reqString, err := buildManticoreSearchRequest(dto)
	if err != nil {
		return nil, err
	}

	urlParams := url2.Values{}
	urlParams.Set("query", reqString)

	url := util.GetEnv("MANTICORE_URL", "http://localhost:9308")
	urlWithQuery := fmt.Sprintf("%s%s?%s", url, "/sql", urlParams.Encode())
	searchRequest, err := http.NewRequest(http.MethodGet, urlWithQuery, nil)
	if err != nil {
		return nil, err
	}
	searchRequest.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	manticoreResponseObject, err := client.Do(searchRequest)

	var manticoreSearchResponse schema.UserManticoreSearchResponse

	err = json.NewDecoder(manticoreResponseObject.Body).Decode(&manticoreSearchResponse)

	if err != nil {
		return nil, errors.New("Manticore is unavailable")
	}

	responseData := buildResponseData(&manticoreSearchResponse)
	paginationInfo := buildPaginationInfo(&manticoreSearchResponse, dto.Limit)

	finalResponseDto := schema.UserSearchResponseDto{
		Data:       responseData,
		Pagination: paginationInfo,
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(manticoreResponseObject.Body)

	return &finalResponseDto, nil

}
