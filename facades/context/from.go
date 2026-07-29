package context

import (
	"strings"

	strutil "git.intra.weibo.com/search_fe/wbutil-go/tools/utils"
)

type From struct {
	raw string // AABBBCXDDE

	app        string // AA, 产品
	year       string // B, 年份
	month      string // B, 月份
	count      string // B, 当月第几个版本
	issued     string // C, 公测/正式
	osPlatform string // X, 系统平台，3-iOS，5-android，A-鸿蒙
	osVersion  string // DD, 系统版本号，一般为01
	channel    string // E, 主版本/预装等

	appVersion string // BBB, app版本号,常用
	os         string // XDD, 常用
	platform   string // XDDE, 后四位，常用

	appPlatform      string // AAXDDE 用于版本比较
	appVersionIssued string // BBBC 用于版本比较
}

type CompareStruct struct {
	Operation string
	Version   string
}

const (
	FromLen = 10

	PlatAndroid = "android"
	PlatIPhone  = "iphone"
	PlatIPad    = "ipad"
	PlatHarmony = "harmony"
)

var platforms = map[string]string{
	"501": PlatAndroid,
	"506": PlatAndroid,
	"301": PlatIPhone,
	"901": PlatIPad,    // iPad HD 版
	"902": PlatIPad,    // iPad
	"A01": PlatHarmony, // 鸿蒙平台
	"A02": PlatHarmony, // 鸿蒙单框架版微博客户端
}

func NewFrom(from string) *From {
	f := &From{raw: from}

	if len(from) != FromLen {
		return f
	}

	f.app = from[:2]
	f.year = from[2:3]
	f.month = from[3:4]
	f.count = from[4:5]
	f.issued = from[5:6]
	f.osPlatform = from[6:7]
	f.osVersion = from[7:9]
	f.channel = from[9:10]

	f.appVersion = from[2:5]
	f.os = from[6:9]
	f.platform = from[6:10]

	f.appPlatform = f.app + f.platform
	f.appVersionIssued = f.appVersion + f.issued

	return f
}

func (f *From) Version() string {
	return f.appVersion
}

func (f *From) SetAppVersion(version string) {
	f.appVersion = version
}

func (f *From) Compare(operator, version string) bool {
	switch operator {
	case "=":
		return f.appVersion == version
	case ">":
		return f.appVersion > version
	case ">=":
		return f.appVersion >= version
	case "<":
		return f.appVersion < version
	case "<=":
		return f.appVersion <= version
	}

	return false
}

// IsVersion 判断请求的客户端是否符合版本要求
// 支持 >, <, >=, <=, !, 空为等于
func (f *From) IsVersion(conditions string) bool {
	if len(conditions) > 2 && (conditions[:2] == ">=" || conditions[:2] == "<=") {
		return f.Compare(conditions[:2], conditions[2:])
	}
	if len(conditions) > 1 && (conditions[:1] == ">" || conditions[:1] == "<") {
		return f.Compare(conditions[:1], conditions[1:])
	}
	if len(conditions) > 1 && string(conditions[0]) == "!" {
		return f.appVersion != conditions[1:]
	}
	return f.Compare("=", conditions)
}

func (f *From) CreateCompareStruct(targetVersion string) CompareStruct {
	compareObj := CompareStruct{}
	targetVersion = strings.TrimSpace(targetVersion)
	if targetVersion == "" {
		return compareObj
	}

	if strings.Contains(targetVersion, ">=") {
		compareObj.Operation = ">="
	} else if strings.Contains(targetVersion, "<=") {
		compareObj.Operation = "<="
	} else if strings.Contains(targetVersion, ">") {
		compareObj.Operation = ">"
	} else if strings.Contains(targetVersion, "<") {
		compareObj.Operation = "<"
	} else {
		compareObj.Operation = "="
	}

	length := len(targetVersion)
	index := length - 3
	compareObj.Version = targetVersion[index:]
	return compareObj
}

func (f *From) GreaterThan(versions map[string]string) bool {
	return f.isFrom(">=", versions)
}

func (f *From) LessThan(versions map[string]string) bool {
	return f.isFrom("<=", versions)
}

func (f *From) isFrom(operator string, versions map[string]string) bool {
	if len(f.raw) != 10 {
		return false
	}

	if _, ok := versions[f.appPlatform]; !ok {
		return false
	}

	target := versions[f.appPlatform]

	return f.Compare(operator, target)
}

// IsApp 通过from前两位判断是否为某类APP
// example:
//
//	f.IsApp("10") //判断是否主端
//	f.IsApp("18") //判断是否为小程序
func (f *From) IsApp(app string) bool {
	return f.app == app
}

func (f *From) IsAndroid() bool {
	return platforms[f.os] == PlatAndroid
}

func (f *From) IsIPhone() bool {
	return platforms[f.os] == PlatIPhone
}

func (f *From) IsIPad() bool {
	return platforms[f.os] == PlatIPad
}

// IsIPadHD 是否是 iPad HD 版（901）
func (f *From) IsIPadHD() bool {
	return f.os == "901"
}

func (f *From) IsIos() bool {
	return f.IsIPhone() || f.IsIPad()
}

// IsHarmony 是否是鸿蒙端
func (f *From) IsHarmony() bool {
	return f.IsHarmonyDKF() || f.IsHarmonySKF()
}

// IsHarmonyDKF 是否是A01是双框架鸿蒙
func (f *From) IsHarmonyDKF() bool {
	return f.os == "A01"
}

// IsHarmonySKF 是否是A02是单框架鸿蒙
func (f *From) IsHarmonySKF() bool {
	return f.os == "A02"
}

// IsWeiboAppByFrom 微博主端
func (f *From) IsWeiboAppByFrom() bool {
	return f.app == "10"
}

// IsLiteWeiboApp 是否为简版微博App
func (f *From) IsLiteWeiboApp() bool {
	return f.channel == "8"
}

// IsVisonPro 是否为VisonPro
func (f *From) IsVisonPro() bool {
	tmpPlatform := f.platform
	if strutil.Substr(tmpPlatform, 0, 1) == "B" && f.appVersion >= "D62" {
		return true
	}

	return false
}

func (f *From) Platform() string {
	p := platforms[f.os]
	if p == "" {
		p = "unknown"
	}

	return p
}

// IsWeiboPro 判断是否为微博专业版
func (f *From) IsWeiboPro() bool {
	return strings.HasPrefix(f.raw, "1F") && strings.HasSuffix(f.raw, "6039")
}
