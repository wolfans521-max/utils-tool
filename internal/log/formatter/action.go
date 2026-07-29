package formatter

import (
	"os"
	"sort"
	"strings"
	"time"

	"git.intra.weibo.com/search_fe/wbutil-go/facades/context"
)

// Action
// 行为码日志的字段结构如下：
//
//	0  time        打码时间
//	1  uid         用户 uid
//	2  act_code    行为码
//	3  oid         行为目标 id
//	4  uicode      本级 ui 编码
//	5  fid         本级页面标记主 id
//	6  lfid        上级页面附属信息
//	7  luicode     上级页面 ui 码
//	8  cardid      cardid
//	9  lcardid     产品 id
//	10 featurecode 根页面 ui 码
//	11 from        版本信息
//	12 wm          渠道值
//	13 oldwm       原始渠道值
//	14 ip          IP
//	15 version     WAP 版本
//	16 aid         aid 信息
//	17 ext         扩展字段
//
// 输出格式：SERVER_IP|time`uid`act_code`...`ext\n
// 其中 SERVER_IP 优先取 SEARCH_ENV_IP，其次 SERVER_ADDR，最后回退 127.0.0.1。
// Action 根据上下文 ctx 格式化行为日志
func Action(_ any, ctx *context.Context) string {
	// 每次调用都获取当前时间，确保日志时间戳准确
	ts := time.Now().Format("2006-01-02 15:04:05")

	dynamic := getDynamicParams(ctx) // 1..10
	fixed := getFixedParams(ctx)     // 11..16
	ext := formatExtForAction(ctx)   // 17

	recordParams := make([]string, 0, 1+len(dynamic)+len(fixed)+1)
	recordParams = append(recordParams, ts)
	recordParams = append(recordParams, dynamic...)
	recordParams = append(recordParams, fixed...)
	recordParams = append(recordParams, ext)

	// 前缀 SERVER_IP：SEARCH_ENV_IP > SERVER_ADDR > 127.0.0.1
	ip := os.Getenv("SEARCH_ENV_IP")
	if ip == "" {
		ip = os.Getenv("SERVER_ADDR")
		if ip == "" {
			ip = "127.0.0.1"
		}
	}

	var b strings.Builder
	b.WriteString(ip)
	b.WriteString("|")
	b.WriteString(strings.Join(recordParams, "`"))
	b.WriteString("\n")
	return b.String()
}

// getDynamicParams 获取易变字段：uid、act_code、oid、uicode、fid、lfid、luicode、cardid、lcardid、featurecode。
func getDynamicParams(ctx *context.Context) []string {
	// uid：login_uid 优先，其次 ctx.GetUserId()
	uid := ""
	if ctx != nil {
		if v := ctx.DefaultRequestParam("login_uid", ""); v != "" {
			uid = v
		} else {
			uid = ctx.GetUserId()
		}
	}

	vals := make(map[string]string, 10)
	vals["uid"] = uid

	keys := []string{"act_code", "oid", "uicode", "fid", "lfid", "luicode", "cardid", "lcardid", "featurecode"}
	for _, k := range keys {
		if ctx != nil {
			vals[k] = ctx.DefaultRequestParam(k, "")
		} else {
			vals[k] = ""
		}
	}

	// uicode 拆分：如 "xxx_yyy_zzz" => luicode=xxx, uicode=zzz, lfid="yyy" 或 "yyy_..."
	uicode := vals["uicode"]
	if strings.Contains(uicode, "_") {
		parts := strings.Split(uicode, "_")
		if len(parts) > 0 {
			vals["uicode"] = parts[len(parts)-1]
			vals["luicode"] = parts[0]
			if len(parts) > 2 {
				vals["lfid"] = strings.Join(parts[1:len(parts)-1], "_")
			}
		}
	}

	return []string{
		vals["uid"],
		vals["act_code"],
		vals["oid"],
		vals["uicode"],
		vals["fid"],
		vals["lfid"],
		vals["luicode"],
		vals["cardid"],
		vals["lcardid"],
		vals["featurecode"],
	}
}

// getFixedParams 获取相对固定的字段：from、wm、oldwm、ip、version、aid。
func getFixedParams(ctx *context.Context) []string {
	var from, wm, oldwm, ip, version, aid string
	if ctx != nil {
		from = ctx.GetFrom()
		wm = ctx.Wm()
		oldwm = ctx.DefaultRequestParam("oldwm", "")
		ip = ctx.ClientIP()
		version = logVersion(ctx)
		aid = ctx.DefaultRequestParam("aid", "")
	}
	if oldwm == "" {
		oldwm = wm
	}
	return []string{from, wm, oldwm, ip, version, aid}
}

// logVersion 获取记录日志所属的版本号
// 返回值：4 代表请求来自触屏版（H5），0 代表来自接口，空字符串代表来自客户端
func logVersion(ctx *context.Context) string {
	if ctx == nil {
		return ""
	}
	// isFromH5: 'h5' == request_client || 'h5' == c
	if ctx.DefaultRequestParam("request_client", "") == "h5" || ctx.C() == "h5" {
		return "4"
	}
	return "0"
}

// formatExtForAction 格式化日志的 ext 参数
func formatExtForAction(ctx *context.Context) string {
	var b strings.Builder

	// path 直接从 ctx.Path() 获取，如果为空则保持空字符串
	path := ""
	if ctx != nil {
		path = ctx.Path()
	}
	// 去掉前导斜杠
	path = strings.TrimPrefix(path, "/")
	b.WriteString("_path:")
	b.WriteString(path)
	b.WriteString("|")

	// launchid：始终输出键，值可为空
	launchid := ""
	if ctx != nil {
		launchid = ctx.DefaultRequestParam("launchid", "")
	}
	b.WriteString("launchid:")
	b.WriteString(launchid)
	b.WriteString("|")

	// 热启动 & 折叠屏参数，仅在 ctx 中存在时输出
	if ctx != nil {
		if s := ctx.DefaultRequestParam("ul_hid", ""); s != "" {
			b.WriteString("ul_hid:")
			b.WriteString(s)
			b.WriteString("|")
		}
		if s := ctx.DefaultRequestParam("ul_sid", ""); s != "" {
			b.WriteString("ul_sid:")
			b.WriteString(s)
			b.WriteString("|")
		}
		if s := ctx.DefaultRequestParam("f_s", ""); s != "" {
			b.WriteString("f_s:")
			b.WriteString(s)
			b.WriteString("|")
		}
		if s := ctx.DefaultRequestParam("a_e", ""); s != "" {
			b.WriteString("a_e:")
			b.WriteString(s)
			b.WriteString("|")
		}

		// 是否 push 进入
		if s := ctx.DefaultRequestParam("ispush", ""); s != "" {
			b.WriteString("ispush:")
			b.WriteString(s)
			b.WriteString("|")
		}

		// gsidnil：如果gsid为空，则输出gsidnil:1
		gsid := ctx.DefaultRequestParam("gsid", "")
		if gsid == "" {
			b.WriteString("gsidnil:1|")
		}
	}

	// ext 参数从 ctx 获取：优先从 ctx.Get("ext") 获取（支持 any 类型），其次从 RequestParam 获取（string 类型）
	var ext any
	if ctx != nil {
		if v, ok := ctx.Get("ext"); ok && v != nil {
			ext = v
		} else {
			ext = ctx.RequestParam("ext")
		}
	}

	switch v := ext.(type) {
	case string:
		// 直接附加字符串
		b.WriteString(v)
	case []any:
		// 处理 numeric-key 的列表：每个元素直接作为一个字段；嵌套 array 用 "|" 拼接
		fields := make([]string, 0, len(v))
		for _, it := range v {
			switch vv := it.(type) {
			case []any:
				inner := make([]string, 0, len(vv))
				for _, x := range vv {
					inner = append(inner, toString(x))
				}
				fields = append(fields, strings.Join(inner, "|"))
			default:
				fields = append(fields, toString(vv))
			}
		}
		b.WriteString(strings.Join(fields, "|"))
	case []string:
		// 处理 []string 类型：每个元素直接作为一个字段
		b.WriteString(strings.Join(v, "|"))
	case map[string]any:
		// 字段顺序与线上日志输出顺序一致
		order := []string{
			// 初始化字段（按线上日志顺序）
			"page", "since_id", "mid", "page_pro", "ext",
			"last", "empty", "orifid", "oriuicode", "request_referer",
			"filter_label_word", "search_ext",
			// 可选字段
			"volume", "scenes", "srid", "search_ssid", "search_vsid",
			"feed_cate", "root_mid", "is_positivity", "reply_type",
			"is_half_screen", "q_extra", "source_mid", "have_mblog",
			"stats_exclude", "searchbar_source", "search_mode", "search_layout",
			"is_blink_flag",
		}
		used := map[string]struct{}{}
		fields := make([]string, 0, len(v))

		for _, k := range order {
			if val, ok := v[k]; ok {
				// ext 字段特殊处理：如果是数组，直接输出数组元素的值（不带key名）
				if k == "ext" {
					switch extVal := val.(type) {
					case []string:
						// []string 类型：输出 "ext:" 后跟数组元素用 "|" 分隔
						fields = append(fields, "ext:"+strings.Join(extVal, "|"))
					case []any:
						// []any 类型：输出 "ext:" 后跟数组元素用 "|" 分隔
						extFields := make([]string, 0, len(extVal))
						for _, item := range extVal {
							extFields = append(extFields, toString(item))
						}
						fields = append(fields, "ext:"+strings.Join(extFields, "|"))
					default:
						fields = append(fields, k+":"+toString(val))
					}
				} else {
					fields = append(fields, k+":"+toString(val))
				}
				used[k] = struct{}{}
			}
		}

		// 追加未在 order 中声明的其它键
		extraKeys := make([]string, 0, len(v))
		for k := range v {
			if _, ok := used[k]; !ok {
				extraKeys = append(extraKeys, k)
			}
		}
		if len(extraKeys) > 0 {
			sort.Strings(extraKeys)
			for _, k := range extraKeys {
				fields = append(fields, k+":"+toString(v[k]))
			}
		}

		b.WriteString(strings.Join(fields, "|"))
	default:
		if ext != nil {
			b.WriteString(toString(ext))
		}
	}

	b.WriteString("|")
	b.WriteString("go_server:1")

	return b.String()
}
