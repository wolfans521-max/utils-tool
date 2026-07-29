package card

type StyleSpacing []float32

type BackgroundStyle struct {
	Type              string   `json:"type,omitempty"`
	Color             string   `json:"color,omitempty"`
	ColorDark         string   `json:"colorDark,omitempty"`
	ColorKey          string   `json:"colorKey,omitempty"`
	ImgUrl            string   `json:"imgUrl,omitempty"`
	DarkMode          string   `json:"darkMode,omitempty"`
	ScaleType         string   `json:"scaleType,omitempty"`
	Corners           []int    `json:"corners,omitempty"`           // 圆角 [左上, 右上, 右下, 左下]
	Colors            []string `json:"colors,omitempty"`            // 渐变颜色数组
	ColorsKey         []string `json:"colorsKey,omitempty"`         // 渐变颜色key数组
	GradientDirection int      `json:"gradientDirection,omitempty"` // 渐变方向 1=从上到下
}

type Divider struct {
	Color string `json:"color,omitempty"`
	Width uint16 `json:"width,omitempty"`
}

type Style struct {
	Background *BackgroundStyle `json:"background,omitempty"`
	Margin     StyleSpacing     `json:"margin,omitempty"`
	Padding    StyleSpacing     `json:"padding,omitempty"`
	Grid       bool             `json:"grid,omitempty"`
	Height     float32          `json:"height,omitempty"`
	Divider    *Divider         `json:"divider,omitempty"`
	Column     int8             `json:"column,omitempty"`
}
