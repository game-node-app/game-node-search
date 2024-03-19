package games

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

func ValidateGameSearchRequest(dtoBytes []byte) (*schema.GameSearchRequestDto, error) {
	var request = schema.GameSearchRequestDto{}
	fmt.Println(string(dtoBytes))
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
	themes := dto.Themes
	platforms := dto.Platforms

	if query == nil {
		return "", errors.New("query parameter empty")
	}

	// Matches all fields
	matchString = fmt.Sprintf("@* %s", *query)
	var genresMatchString = ""
	var themesMatchString = ""
	var platformsMatchString = ""

	if genres != nil && len(*genres) > 0 {
		genresMatchString = " @genres_names "
		for i, v := range *genres {
			if i > 0 {
				genresMatchString = fmt.Sprintf("%s|%s", genresMatchString, v)
				continue
			}
			genresMatchString = fmt.Sprintf("%s%s", genresMatchString, v)
		}

	}

	if themes != nil && len(*themes) > 0 {
		themesMatchString = " @themes_names "
		for i, v := range *themes {
			if i > 0 {
				themesMatchString = fmt.Sprintf("%s|%s", themesMatchString, v)
				continue
			}
			themesMatchString = fmt.Sprintf("%s|%s", themesMatchString, v)
		}

	}

	if platforms != nil && len(*platforms) > 0 {
		platformsMatchString = " @(platforms_names,platforms_abbreviations) "
		for i, v := range *platforms {
			if i > 0 {
				platformsMatchString = fmt.Sprintf("%s|%s", platformsMatchString, v)
				continue
			}
			platformsMatchString = fmt.Sprintf("%s%s", platformsMatchString, v)
		}
	}

	finalMatchString := fmt.Sprintf("%s%s%s%s", matchString, genresMatchString, themesMatchString, platformsMatchString)

	return finalMatchString, nil

}

func buildManticoreFilterString(dto *schema.GameSearchRequestDto) (string, error) {
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

func buildManticorePaginationString(dto *schema.GameSearchRequestDto) string {
	limit := dto.Limit
	if limit == nil || *limit == 0 {
		i := schema.DEFAULT_LIMIT
		limit = &i
	}
	page := dto.Page
	if page == nil || *page == 0 {
		u := 1
		page = &u
	}

	offset := (*page - 1) * *limit
	paginationString := fmt.Sprintf("LIMIT %d OFFSET %d", *limit, offset)

	return paginationString
}

func buildManticoreOrderString() string {
	return "ORDER BY num_likes DESC, num_views DESC"
}

func buildManticoreSearchRequest(dto *schema.GameSearchRequestDto) (string, error) {

	matchString, _ := buildManticoreMatchString(dto)
	filterString, _ := buildManticoreFilterString(dto)
	paginationString := buildManticorePaginationString(dto)
	orderString := buildManticoreOrderString()

	selectString := fmt.Sprintf("SELECT * FROM games WHERE match('%s') %s %s %s;", matchString, filterString, orderString, paginationString)

	return selectString, nil

}

// GameSearchHandler search handler
//
//	@Summary      Searches for games using Manticore engine
//	@Description  Returns a parsed search response from the Manticore engine
//	@Tags         search
//	@Accept       json
//	@Produce      json
//	@Param        query   body      schema.GameSearchRequestDto  true  "Account ID"
//	@Success      200  {object}   schema.GameSearchResponseDto
//	@Router       /search/games [post]
func GameSearchHandler(dto *schema.GameSearchRequestDto) (*schema.GameSearchResponseDto, error) {

	reqString, err := buildManticoreSearchRequest(dto)
	if err != nil {
		return nil, err
	}

	fmt.Println(reqString)

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

	if err != nil {
		return nil, errors.New("Manticore is unavailable")
	}

	var manticoreResponseDto schema.ManticoreSearchResponse
	err = json.NewDecoder(manticoreResponseObject.Body).Decode(&manticoreResponseDto)
	if err != nil || manticoreResponseDto.Hits == nil {
		return nil, errors.New("failed to fetch data. check query parameters")
	}

	var response = schema.GameSearchResponseDto{}
	data := buildResponseData(&manticoreResponseDto)
	pagination := buildPaginationInfo(&manticoreResponseDto, dto.Limit)
	response.Data = *data
	response.Pagination = *pagination

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(manticoreResponseObject.Body)

	return &response, nil
}
