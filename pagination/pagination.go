package pagination

// =============================================================================
//
//	Struct Definition
//
// =============================================================================
type PaginationState struct {
	// TODO: I could add an array of integers with the middle element being
	// guaranteed to be the current page. This could then be used in the middle
	// of the current and previous page buttons.
	Previous         int
	Next             int
	DisabledPrevious bool
	DisabledNext     bool
}
