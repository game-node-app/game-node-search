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
	assert.NotNil(t, searchResponse.Hits)
	assert.NotNil(t, searchResponse.Hits.Hits)
	assert.Greater(t, *searchResponse.Hits.Total, 0)
	assert.Greater(t, len(searchResponse.Hits.Hits), 0)

}
