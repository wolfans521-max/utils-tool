package card

type MeasureHeight struct {
	Expression        string `json:"expression,omitempty"`
	BundleVersionCode string `json:"bundleVersionCode,omitempty"`
}

type C236 struct {
	BaseItem
	ReadtimeType string `json:"readtimetype,omitempty"`
	Display      struct {
		MaxHeight         int           `json:"maxHeight,omitempty"`
		CoverImageUrl     string        `json:"coverImageUrl,omitempty"`
		CoverImageUrlDark string        `json:"coverImageUrlDark,omitempty"`
		DisplayHeight     int           `json:"displayHeight,omitempty"`
		DisplayRatio      float32       `json:"displayRatio,omitempty"`
		MeasureHeight     *MeasureHeight `json:"measureHeight,omitempty"`
	} `json:"display"`
	DisplayBackup struct {
		DisplayHeight     int     `json:"displayHeight,omitempty"`
		DisplayRatio      float32 `json:"displayRatio,omitempty"`
		CoverImageUrl     string  `json:"coverImageUrl,omitempty"`
		CoverImageUrlDark string  `json:"coverImageUrlDark,omitempty"`
	}
	Scheme                string         `json:"scheme,omitempty"`
	WboxScheme            string         `json:"wboxScheme"`
	WboxMinRuntimeVersion int            `json:"wboxMinRuntimeVersion,omitempty"`
	WboxMinSDKVersion     int            `json:"wboxMinSDKVersion,omitempty"`
	WboxParam             map[string]any `json:"wboxParam"`
}

func NewCard236() *C236 {
	return &C236{BaseItem: BaseItem{Base: Base{CardType: 236}}}
}
