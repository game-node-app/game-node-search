package search

import (
	"bytes"
	"game-node-search/schema"
	"game-node-search/util"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func ValidateSearchRequest(r *http.Request) (*schema.GameSearchRequestDto, error) {
	dto := &schema.GameSearchRequestDto{
		Index: "gamenode",
	}

	err := json.NewDecoder(r.Body).Decode(dto)
	if err != nil {
		return nil, err
	}

	return dto, nil

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
func Handler(w http.ResponseWriter, r *http.Request) *schema.GameSearchResponseDto {
	defer r.Body.Close()

	dto, err := ValidateSearchRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	// Marshals the struct back to Json bytes (so it's lowercase)
	jsonDto, err := json.Marshal(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	url := util.GetEnv("MANTICORE_URL"+"/search", "http://localhost:9308/search")
	searchRequest, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonDto))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	searchRequest.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	manticoreResponseObject, err := client.Do(searchRequest)

	if err != nil || manticoreResponseObject == nil || manticoreResponseObject.StatusCode != 200 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	var manticoreResponseDto schema.ManticoreSearchResponse
	err = json.NewDecoder(manticoreResponseObject.Body).Decode(&manticoreResponseDto)

	if err != nil {
		var errorResponse schema.ManticoreErrorResponse
		if json.NewDecoder(manticoreResponseObject.Body).Decode(&errorResponse) == nil {
			errorBytes, err := io.ReadAll(manticoreResponseObject.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write(errorBytes)
				return nil
			}
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	result, err := ParseManticoreResponse(&manticoreResponseDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	return result
}
