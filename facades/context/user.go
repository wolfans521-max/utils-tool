package context

type User struct {
	Id string
}

func (ctx *Context) GetUserId() string {
	if cuid := ctx.DefaultRequestParam("cuid", ""); cuid != "" {
		return cuid
	}

	if xLoginUid := ctx.DefaultHeader("x-login-uid", ""); xLoginUid != "" {
		return xLoginUid
	}

	if loginUid := ctx.DefaultRequestParam("login_uid", ""); loginUid != "" {
		return loginUid
	}

	return ""
}

func (ctx *Context) IsTeenager() bool {
	if isTeenagerflag := ctx.RequestInt("is_teenager"); isTeenagerflag == 1 {
		return true
	}

	if xTeenagerFlag := ctx.DefaultHeader("X-Teenager-Flag", ""); xTeenagerFlag != "" {
		return true
	}

	return false
}

func (ctx *Context) IsGuest() bool {
	uid := ctx.GetUserId()

	return len(uid) >= 13 || len(uid) <= 3
}
