package search

import (
	"fmt"
	"game-node-search/schema"
	"game-node-search/util"
	"strconv"
	"strings"
	"time"
)

func ParseManticoreResponse(mr *schema.ManticoreSearchResponse) (*schema.GameSearchResponseDto, error) {
	mrBytes, _ := json.Marshal(mr)
	resultJson := make(map[string]interface{})
	err := json.Unmarshal(mrBytes, &resultJson)
	if err != nil {
		fmt.Println("json.Unmarshal failed:", err)
		return nil, err
	}

	convertedMap := normalizeManticoreResponse(resultJson)
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
func normalizeManticoreResponse(src interface{}) *map[string]interface{} {
	dst := make(map[string]interface{})
	dateFields := []string{"first_release_date", "created_at", "updated_at"}

	mappedSrc, isMap := src.(map[string]interface{})
	if isMap && len(mappedSrc) > 0 {
		for k, v := range mappedSrc {
			camelCaseKey := stringToCamelCase(k)

			switch v := v.(type) {
			case map[string]interface{}:

				dst[camelCaseKey] = normalizeManticoreResponse(v)

			case []interface{}:
				var convertedSlice []interface{}
				/**
				This iterates over the array of hits []ManticoreResponseHit
				*/
				for _, av := range v {
					avMap, isAvMap := av.(map[string]interface{})
					if isAvMap {
						id, hasId := avMap["_id"]
						idString, isIdString := id.(string)
						// _id comes as string by default for some reason
						if isIdString {
							idAsUint, _ := strconv.ParseUint(idString, 10, 64)
							id = idAsUint
							avMap["_id"] = idAsUint
						}
						source, hasSource := avMap["_source"]
						sourceMap, isSourceMap := source.(map[string]interface{})
						if hasId && hasSource && isSourceMap {
							sourceMap["id"] = id
						}
					}
					convertedSlice = append(convertedSlice, normalizeManticoreResponse(av))
				}
				dst[camelCaseKey] = convertedSlice
			case float64:
				if util.Contains(dateFields, k) {
					dst[camelCaseKey] = time.Unix(int64(v), 0)
				} else {
					dst[camelCaseKey] = v
				}
			case int:
				if util.Contains(dateFields, k) {
					dst[camelCaseKey] = time.Unix(int64(v), 0)
				} else {
					dst[camelCaseKey] = v
				}

			default:
				// Convert the key to camel case and assign the value
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
