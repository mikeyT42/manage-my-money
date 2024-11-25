package pagination

// =============================================================================
//
//	Struct Definition
//
// =============================================================================
type PaginationState struct {
	Previous         int
	Next             int
	DisabledPrevious bool
	DisabledNext     bool
}
