package card

// C57 股票行情卡片
type C57 struct {
	BaseItem
	Date           string      `json:"date,omitempty"`            // 右侧右下文字
	DisplayArrow   int8        `json:"display_arrow,omitempty"`   // 是否显示箭头
	OscillatePrice string      `json:"oscillate_price,omitempty"` // 右侧左上文字（涨跌额）
	OscillateRate  string      `json:"oscillate_rate,omitempty"`  // 右侧左下文字（涨跌幅）
	Price          string      `json:"price,omitempty"`           // 左侧最大文字（当前价格）
	PriceColorType int8        `json:"price_color_type,omitempty"`
	StateDesc      string      `json:"state_desc,omitempty"`      // 右侧右上文字（状态描述）
	StockInfo      *C57StockInfo `json:"stock_info,omitempty"`    // 下部股票信息
}

// C57StockInfo 股票详细信息
type C57StockInfo struct {
	Scheme          string             `json:"scheme,omitempty"`
	StockPriceInfos []*C57StockPriceInfo `json:"stock_price_infos,omitempty"` // 下部实用信息数组
}

// C57StockPriceInfo 股票价格信息项
type C57StockPriceInfo struct {
	Price string `json:"price,omitempty"` // 价格
	Desc  string `json:"desc,omitempty"`  // 描述（今开、昨收、最低、最高等）
}

// NewCard57 创建Card57实例
func NewCard57() *C57 {
	return &C57{BaseItem: BaseItem{Base: Base{CardType: 57}}}
}
