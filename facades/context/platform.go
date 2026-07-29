package context

func (ctx *Context) IsFoldDevice() bool {
	fs := ctx.DefaultRequestParam("f_s", "")
	if fs != "" {
		return true
	}

	return false
}

func (ctx *Context) IsPad() bool {
	isPad := ctx.DefaultRequestParam("is_pad", "")
	if isPad == "true" {
		return true
	}

	return false
}
