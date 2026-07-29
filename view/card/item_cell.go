package card

type cellBase struct {
	ItemBase
	Type string `json:"type"`
}

type CellSpanSpan struct {
	Style *Style `json:"style,omitempty"`
}

type CellSpan struct {
	cellBase
	Span    *CellSpanSpan `json:"span,omitempty"`
	Style   *Style        `json:"style,omitempty"`
	ItemExt *ItemExt      `json:"itemExt,omitempty"`
}

func NewItemCellSpan() *CellSpan {
	return &CellSpan{
		cellBase: cellBase{
			ItemBase: ItemBase{Category: CategoryCell},
			Type:     "span",
		},
	}
}
