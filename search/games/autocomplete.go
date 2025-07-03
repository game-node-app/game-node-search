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
//	@Param        query   body      games.GameSearchRequestDto  true  "Request body"
//	@Success      200  {object}   games.GameAutocompleteResponseDto
//	@Router       /search/games/autocomplete [post]
func Autocomplete(dto *GameSearchRequestDto) (*GameAutocompleteResponseDto, error) {
	autocompleteQuery := fmt.Sprintf("%s*", dto.Query)
	limit := int32(5)

	dto.Query = autocompleteQuery
	dto.Limit = &limit

	searchResponseDto, err := Search(dto)
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
