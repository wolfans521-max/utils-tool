package card

type ItemCard struct {
	ItemBase
	Data    any      `json:"data"`
	Style   *Style   `json:"style,omitempty"`
	ItemExt *ItemExt `json:"itemExt,omitempty"`
}

func NewItemCard() *ItemCard {
	return &ItemCard{
		ItemBase: ItemBase{
			Category: CategoryCard,
		},
	}
}
