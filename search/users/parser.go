package users

import (
	"game-node-search/schema"
	Manticoresearch "github.com/manticoresoftware/manticoresearch-go"
	"log/slog"
)

func buildPaginationInfo(mr *Manticoresearch.SearchResponse, limit *int32) *schema.PaginationInfo {
	limitToUse := limit
	if limit == nil || *limitToUse == 0 {
		i := schema.DefaultLimit
		limitToUse = &i
	}
	var paginationInfo = schema.PaginationInfo{
		TotalItems:  0,
		TotalPages:  0,
		HasNextPage: false,
	}

	if mr.Hits != nil && mr.Hits.Hits != nil && len(mr.Hits.Hits) > 0 {
		hitsLen := uint64(len(mr.Hits.Hits))
		totalItems := uint64(*mr.Hits.Total)
		limitU := uint64(*limitToUse)

		paginationInfo.TotalItems = totalItems
		paginationInfo.TotalPages = uint16((totalItems + limitU - 1) / limitU)
		paginationInfo.HasNextPage = totalItems > hitsLen
	}

	return &paginationInfo
}

func buildResponseData(mr *Manticoresearch.SearchResponse) (*UserSearchResponseData, error) {
	var users []UserDto

	var data = UserSearchResponseData{
		Took:    mr.Took,
		Profile: mr.Profile,
		Items:   &[]UserDto{},
	}

	if mr.Hits != nil && mr.Hits.Hits != nil && len(mr.Hits.Hits) > 0 {
		for _, hit := range mr.Hits.Hits {
			hitJson, err := json.Marshal(hit)
			if err != nil {
				slog.Error("Error while decoding results to JSON.", "err", err)
				return nil, err
			}

			var convertedHit UserManticoreResponseHit

			err = json.Unmarshal(hitJson, &convertedHit)
			if err != nil {
				slog.Error("Error while decoding results from JSON.", "err", err)
				return nil, err
			}

			userDto := UserDto{
				UserId:   convertedHit.Source.UserId,
				Username: convertedHit.Source.Username,
			}

			users = append(users, userDto)
		}
	}

	if users != nil {
		data.Items = &users
	}

	return &data, nil
}
