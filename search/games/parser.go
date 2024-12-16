package games

import (
	"game-node-search/schema"
	"regexp"
	"strconv"
	"time"
)

func parseQuery(q string) string {
	var result = q
	quotesRegx := regexp.MustCompile(`['"]`)
	result = quotesRegx.ReplaceAllString(result, "")

	return result
}

func buildPaginationInfo(mr *schema.ManticoreSearchResponse, limit *int, currentPage *int) *schema.PaginationInfo {
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
		limitU := uint64(*limitToUse)
		hitsLen := uint64(len(mr.Hits.Hits))
		totalItems := uint64(*mr.Hits.Total)
		totalPages := uint16((totalItems + limitU - 1) / limitU)

		currentPageU := uint16(*currentPage)

		paginationInfo.TotalItems = totalItems
		paginationInfo.TotalPages = totalPages
		paginationInfo.HasNextPage = totalItems > hitsLen && currentPageU < totalPages
	}

	return &paginationInfo

}

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
				GenresNames:            hit.Source.GenresNames,
				ThemesNames:            hit.Source.ThemesNames,
				PlatformsNames:         hit.Source.PlatformsNames,
				PlatformsAbbreviations: hit.Source.PlatformsAbbreviations,
				KeywordsNames:          hit.Source.KeywordsNames,
				NumViews:               hit.Source.NumViews,
				NumLikes:               hit.Source.NumLikes,
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
