package card

type ItemFeed struct {
	ItemBase
	Data    map[string]any `json:"data"`
	ItemExt *ItemExt       `json:"itemExt,omitempty"`
	Style   *Style         `json:"style,omitempty"`
}

func NewItemFeed() *ItemFeed {
	return &ItemFeed{
		ItemBase: ItemBase{
			Category: CategoryFeed,
		},
	}
}
