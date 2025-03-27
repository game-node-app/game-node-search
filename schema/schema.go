package schema

type AnyMap map[string]interface{}

const DefaultLimit int32 = 20

type PaginationInfo struct {
	TotalItems  uint64 `json:"totalItems"`
	TotalPages  uint16 `json:"totalPages"`
	HasNextPage bool   `json:"hasNextPage"`
}
