package search

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

func ValidateSearchRequest(dtoBytes []byte) (*schema.GameSearchRequestDto, error) {
	var request = schema.GameSearchRequestDto{}
	err := json.Unmarshal(dtoBytes, &request)
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// buildManticoreMatchString
// Builds the string to be inserted in the 'match' SQL function.
func buildManticoreMatchString(dto *schema.GameSearchRequestDto) (string, error) {
	var matchString string
	query := &dto.Query
	genres := dto.Genres
	platforms := dto.Platforms

	if query == nil {
		return "", errors.New("query parameter empty")
	}

	// Matches all fields
	matchString = fmt.Sprintf("@* %s", *query)
	var genresMatchString = ""
	var platformsMatchString = ""

	if genres != nil && len(*genres) > 0 {
		genresMatchString = "@genres_names "
		for i, v := range *genres {
			if i > 0 {
				genresMatchString = fmt.Sprintf("%s|%s", genresMatchString, v)
				continue
			}
			genresMatchString = fmt.Sprintf("%s%s", genresMatchString, v)
		}

	}

	if platforms != nil && len(*platforms) > 0 {
		platformsMatchString = "@(platforms_names,platforms_abbreviations) "
		for i, v := range *platforms {
			if i > 0 {
				platformsMatchString = fmt.Sprintf("%s|%s", platformsMatchString, v)
				continue
			}
			platformsMatchString = fmt.Sprintf("%s%s", platformsMatchString, v)
		}
	}

	return matchString + genresMatchString + platformsMatchString, nil

}

func buildManticoreFilterString(dto *schema.GameSearchRequestDto) (string, error) {
	limit := dto.Limit
	if limit == nil || *limit == 0 {
		u := 20
		limit = &u
	}
	page := dto.Page
	if page == nil || *page == 0 {
		u := 1
		page = &u
	}

	offset := (*page - 1) * *limit
	paginationString := fmt.Sprintf("OFFSET %d LIMIT %d", offset, limit)
	var filterString string

	if dto.Category != nil && len(*dto.Category) > 0 {
		var categoryFilterArrayString = ""
		for _, v := range *dto.Category {
			categoryFilterArrayString = fmt.Sprintf("%s,%d", categoryFilterArrayString, v)
		}

		filterString = fmt.Sprintf("AND category IN (%s)", categoryFilterArrayString)
	}
	if dto.Status != nil && len(*dto.Status) > 0 {
		var statusFilterArrayString = ""
		for _, v := range *dto.Status {
			statusFilterArrayString = fmt.Sprintf("%s,%d", statusFilterArrayString, v)
		}

		filterString = fmt.Sprintf("AND category IN (%s)", statusFilterArrayString)

	}

}

func buildManticoreSearchRequest(dto *schema.GameSearchRequestDto) (string, error) {

	matchString, _ := buildManticoreMatchString(dto)

	selectString := fmt.Sprintf("SELECT * FROM gamenode WHERE match('%s')%s", matchString, "")

	return selectString, nil

}

// Handler search handler
//
//	@Summary      Searches using Manticore engine
//	@Description  Returns a parsed search response from the Manticore engine
//	@Tags         search
//	@Accept       json
//	@Produce      json
//	@Param        query   body      schema.GameSearchRequestDto  true  "Account ID"
//	@Success      200  {array}   schema.GameSearchResponseDto
//	@Router       /search [post]
func Handler(dto *schema.GameSearchRequestDto) (*schema.GameSearchResponseDto, error) {

	reqString, err := buildManticoreSearchRequest(dto)
	if err != nil {
		return nil, err
	}

	fmt.Print(reqString)

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
	defer manticoreResponseObject.Body.Close()

	if err != nil || manticoreResponseObject == nil || manticoreResponseObject.StatusCode != 200 {
		return nil, err
	}

	var manticoreResponseDto schema.ManticoreSearchResponse
	err = json.NewDecoder(manticoreResponseObject.Body).Decode(&manticoreResponseDto)

	if err != nil {
		var errorResponse schema.ManticoreErrorResponse
		if json.NewDecoder(manticoreResponseObject.Body).Decode(&errorResponse) == nil {
			_, err := io.ReadAll(manticoreResponseObject.Body)
			if err != nil {
				return nil, err
			}
		}

		if &errorResponse != nil {
			return nil, errors.New(errorResponse.Error)
		}
		return nil, err
	}

	result, err := ParseManticoreResponse(&manticoreResponseDto)
	if err != nil {
		return nil, err
	}

	return result, nil
}
