package games

import (
	"fmt"
)

// Autocomplete games autocomplete handler
//
//	@Summary      Autocomplete games names using Manticore engine
//	@Description  Returns a parsed search response from the Manticore engine
//	@Tags         search
//	@Accept       json
//	@Produce      json
//	@Param        query   body      games.GameAutocompleteRequestDto  true  ""
//	@Success      200  {object}   games.GameAutocompleteResponseDto
//	@Router       /search/games/autocomplete [post]
func Autocomplete(dto *GameAutocompleteRequestDto) (*GameAutocompleteResponseDto, error) {
	autocompleteQuery := fmt.Sprintf("^\"%s*\"", dto.Query)

	limit := int32(5)

	searchRequest := GameSearchRequestDto{
		Query: autocompleteQuery,
		Limit: &limit,
	}

	searchResponseDto, err := Search(&searchRequest)
	if err != nil {
		return nil, err
	}

	items := searchResponseDto.Data.Items

	response := GameAutocompleteResponseDto{Data: make([]string, 0), Total: 0}

	if items != nil && len(*items) > 0 {
		for _, game := range *items {
			response.Data = append(response.Data, game.Name)
		}

		response.Total = uint(len(*items))
	}

	return &response, nil
}
