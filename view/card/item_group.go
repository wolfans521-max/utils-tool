package card

type ItemGroup struct {
	Category      Category `json:"category"`
	ItemExt       *ItemExt `json:"itemExt,omitempty"`
	Items         []any    `json:"items"`
	Style         *Style   `json:"style,omitempty"`
	Type          string   `json:"type"`
	ReadtimeType  string   `json:"readtimetype,omitempty"`
	AnalysisExtra string   `json:"analysis_extra,omitempty"`
	ItemCategory  string   `json:"item_category,omitempty"`
	ItemID        string   `json:"itemId,omitempty"`
}

func NewItemGroup() *ItemGroup {
	return &ItemGroup{
		Category: CategoryGroup,
	}
}
