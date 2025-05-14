package exception

type ItemListing struct {
}

func (e *ItemListing) Error() string {
	return "Failed to list items"
}
