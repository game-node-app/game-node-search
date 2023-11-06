package search

import (
	"bytes"
	"fmt"
	"game-node-search/schema"
	"game-node-search/util"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"strings"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func validateSearchRequest(r *http.Request) (*schema.GameSearchRequestDto, error) {
	dto := &schema.GameSearchRequestDto{
		Index: "gamenode",
	}

	err := json.NewDecoder(r.Body).Decode(dto)
	if err != nil {
		return nil, err
	}

	return dto, nil

}

func parseManticoreResponse(mr *schema.ManticoreSearchResponse) (*schema.GameSearchResponseDto, error) {
	mrBytes, _ := json.Marshal(mr)
	resultJson := make(map[string]interface{})
	err := json.Unmarshal(mrBytes, &resultJson)
	if err != nil {
		fmt.Println("json.Unmarshal failed:", err)
		return nil, err
	}

	convertedMap := snakeToCamelCase(resultJson)
	convertedMapBytes, _ := json.Marshal(convertedMap)
	var searchResponse schema.GameSearchResponseDto
	err = json.Unmarshal(convertedMapBytes, &searchResponse)
	if err != nil {
		fmt.Println("json.Unmarshal failed:", err)
		return nil, err
	}

	return &searchResponse, nil
}

// Converts all keys of a given type from snake_case to camelCase.
// Since Go doesn't support type unions, we need to use interface{} and manually check for a map.
// It's ugly, but it works.
func snakeToCamelCase(src interface{}) *map[string]interface{} {
	dst := make(map[string]interface{})
	mappedSrc, isMap := src.(map[string]interface{})
	if isMap && len(mappedSrc) > 0 {
		for k, v := range mappedSrc {
			switch v := v.(type) {
			case map[string]interface{}:
				camelCaseKey := stringToCamelCase(k)
				dst[camelCaseKey] = snakeToCamelCase(v)
			case []interface{}:
				var convertedSlice []interface{}
				for _, av := range v {
					convertedSlice = append(convertedSlice, snakeToCamelCase(av))
				}
				camelCaseKey := stringToCamelCase(k)
				dst[camelCaseKey] = convertedSlice
			default:
				// Convert the key to camel case and assign the value
				camelCaseKey := stringToCamelCase(k)
				dst[camelCaseKey] = v
			}
		}
	}

	return &dst
}

func stringToCamelCase(s string) string {
	// Removes leading underscore from fields like "_score" and "_id"
	if strings.HasPrefix(s, "_") {
		s = s[1:]
	}
	ss := strings.Split(s, "_")
	for i := 1; i < len(ss); i++ {
		ss[i] = strings.Title(ss[i])
	}
	return strings.Join(ss, "")
}

func Handler(w http.ResponseWriter, r *http.Request) *schema.GameSearchResponseDto {
	defer r.Body.Close()

	dto, err := validateSearchRequest(r)
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

	result, err := parseManticoreResponse(&manticoreResponseDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	return result
}
