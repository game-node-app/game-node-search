package games

import (
	"game-node-search/schema"
	Manticoresearch "github.com/manticoresoftware/manticoresearch-go"
	"github.com/mitchellh/mapstructure"
	"log/slog"
	"strconv"
	"time"
)

func buildPaginationInfo(mr *Manticoresearch.SearchResponse, limit *int32, currentPage *int32) *schema.PaginationInfo {
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
			var convertedHit ManticoreResponseHit
			err := mapstructure.Decode(hit, convertedHit)
			if err != nil {
				slog.Error("Error while decoding results to return type.", "err", err)
				return nil, err
			}

			hitIdNumber, _ := strconv.ParseUint(convertedHit.ID, 10, 64)

			searchGame := SearchGame{
				ID:                     hitIdNumber,
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

	return &data, nil
}
