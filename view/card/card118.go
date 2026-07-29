package card

type C118CardPadding struct {
	Left   int `json:"left"`
	Top    int `json:"top"`
	Right  int `json:"right"`
	Bottom int `json:"bottom"`
}

type C118 struct {
	BaseItem
	LeftPadding   int             `json:"left_padding"`
	RightPadding  int             `json:"right_padding"`
	LoopInterval  int             `json:"loop_interval,omitempty"`
	NewStyle      int             `json:"newStyle"`
	CardPadding   C118CardPadding `json:"card_padding"`
	Items         []C119          `json:"items,omitempty"`
	TopPadding    int             `json:"top_padding,omitempty"`
	BottomPadding int             `json:"bottom_padding,omitempty"`
}

func NewCard118() *C118 {
	return &C118{BaseItem: BaseItem{Base: Base{CardType: 118}}}
}
