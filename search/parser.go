package search

import (
	"game-node-search/schema"
	"strconv"
	"time"
)

func buildResponseData(mr *schema.ManticoreSearchResponse) *schema.ResponseData {
	var searchGames []schema.SearchGame

	var data = schema.ResponseData{
		Took:    mr.Took,
		Profile: mr.Profile,
	}

	if mr.Hits != nil && mr.Hits.Hits != nil && len(mr.Hits.Hits) > 0 {
		for _, hit := range mr.Hits.Hits {
			hitIdNumber, _ := strconv.ParseUint(hit.ID, 10, 64)
			searchGame := schema.SearchGame{
				ID:                     hitIdNumber,
				Name:                   hit.Source.Name,
				Slug:                   hit.Source.Slug,
				Summary:                hit.Source.Summary,
				Storyline:              hit.Source.Storyline,
				Checksum:               hit.Source.Checksum,
				AggregatedRating:       hit.Source.AggregatedRating,
				AggregatedRatingCount:  hit.Source.AggregatedRatingCount,
				Category:               hit.Source.Category,
				Status:                 hit.Source.Status,
				CoverUrl:               hit.Source.CoverUrl,
				NumViews:               hit.Source.NumViews,
				NumLikes:               hit.Source.NumLikes,
				GenresNames:            hit.Source.GenresNames,
				ThemesNames:            hit.Source.ThemesNames,
				PlatformsNames:         hit.Source.PlatformsNames,
				PlatformsAbbreviations: hit.Source.PlatformsAbbreviations,
				KeywordsNames:          hit.Source.KeywordsNames,
				Source:                 hit.Source.Source,
				CreatedAt:              time.Unix(int64(hit.Source.CreatedAt), 0),
				FirstReleaseDate:       time.Unix(int64(hit.Source.FirstReleaseDate), 0),
				UpdatedAt:              time.Unix(int64(hit.Source.UpdatedAt), 0),
			}
			searchGames = append(searchGames, searchGame)
		}
	}

	data.Items = &searchGames

	return &data
}

func BuildPaginationInfo(mr *schema.ManticoreSearchResponse, limit *int) *schema.PaginationInfo {
	limitToUse := limit
	if limit == nil || *limitToUse == 0 {
		i := schema.DEFAULT_LIMIT
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
