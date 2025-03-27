package games

import (
	"game-node-search/schema"
	Manticoresearch "github.com/manticoresoftware/manticoresearch-go"
	"log/slog"
	"time"
)

func buildPaginationInfo(mr *Manticoresearch.SearchResponse, limit *int32, currentPage *int32) *schema.PaginationInfo {
	limitToUse := limit
	if limit == nil || *limitToUse == 0 {
		i := schema.DefaultLimit
		limitToUse = &i
	}

	if currentPage == nil {
		p := int32(1)
		currentPage = &p
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

func buildResponseData(mr *Manticoresearch.SearchResponse) (*ResponseData, error) {
	var searchGames []SearchGame

	var data = ResponseData{
		Took:    mr.Took,
		Profile: mr.Profile,
	}

	if mr.Hits != nil && mr.Hits.Hits != nil && len(mr.Hits.Hits) > 0 {

		for _, hit := range mr.Hits.Hits {
			hitJson, err := json.Marshal(hit)
			if err != nil {
				slog.Error("Error while decoding results to JSON.", "err", err)
				return nil, err
			}

			var convertedHit ManticoreResponseHit

			err = json.Unmarshal(hitJson, &convertedHit)
			if err != nil {
				slog.Error("Error while decoding results from JSON.", "err", err)
				return nil, err
			}

			searchGame := SearchGame{
				ID:                     convertedHit.ID,
				Name:                   convertedHit.Source.Name,
				Slug:                   convertedHit.Source.Slug,
				Summary:                convertedHit.Source.Summary,
				Storyline:              convertedHit.Source.Storyline,
				Checksum:               convertedHit.Source.Checksum,
				AggregatedRating:       convertedHit.Source.AggregatedRating,
				AggregatedRatingCount:  convertedHit.Source.AggregatedRatingCount,
				Category:               convertedHit.Source.Category,
				Status:                 convertedHit.Source.Status,
				CoverUrl:               convertedHit.Source.CoverUrl,
				GenresNames:            convertedHit.Source.GenresNames,
				ThemesNames:            convertedHit.Source.ThemesNames,
				PlatformsNames:         convertedHit.Source.PlatformsNames,
				PlatformsAbbreviations: convertedHit.Source.PlatformsAbbreviations,
				KeywordsNames:          convertedHit.Source.KeywordsNames,
				NumViews:               convertedHit.Source.NumViews,
				NumLikes:               convertedHit.Source.NumLikes,
				Source:                 convertedHit.Source.Source,
				CreatedAt:              time.Unix(int64(convertedHit.Source.CreatedAt), 0),
				FirstReleaseDate:       time.Unix(int64(convertedHit.Source.FirstReleaseDate), 0),
				UpdatedAt:              time.Unix(int64(convertedHit.Source.UpdatedAt), 0),
			}
			searchGames = append(searchGames, searchGame)
		}
	}

	data.Items = &searchGames

	if searchGames == nil {
		data.Items = &[]SearchGame{}
	}

	return &data, nil
}
