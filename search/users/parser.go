package users

import (
	"game-node-search/schema"
	Manticoresearch "github.com/manticoresoftware/manticoresearch-go"
	"github.com/mitchellh/mapstructure"
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

func buildResponseData(mr *Manticoresearch.SearchResponse) (*schema.UserSearchResponseData, error) {
	var users []schema.UserDto
	var data = schema.UserSearchResponseData{
		Took:    mr.Took,
		Profile: mr.Profile,
	}

	if mr.Hits != nil && mr.Hits.Hits != nil && len(mr.Hits.Hits) > 0 {
		for _, hit := range mr.Hits.Hits {
			var convertedHit schema.UserManticoreResponseHit

			err := mapstructure.Decode(hit, convertedHit)

			if err != nil {
				return nil, err
			}

			userDto := schema.UserDto{
				UserId:   convertedHit.Source.UserId,
				Username: convertedHit.Source.Username,
			}

			users = append(users, userDto)
		}
	}

	data.Items = &users

	return &data, nil
}
