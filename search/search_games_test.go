package search

import (
	"game-node-search/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSearch(t *testing.T) {
	reqDto := &schema.GameSearchRequestDto{
		Query: "the witcher 3",
	}

	searchResponse, err := Handler(reqDto)

	assert.NoError(t, err)
	assert.NotNil(t, searchResponse.Data.Items)
	assert.NotNil(t, searchResponse.Pagination)
	greaterThan := uint64(0)
	assert.Greater(t, searchResponse.Pagination.TotalItems, greaterThan)
	assert.Greater(t, len(*searchResponse.Data.Items), 0)

}
