package mp

import "strconv"

func Error(msg string, args ...string) Meta {
	return New("Error").
		WithAttr("message", msg).
		WithAttr(args...)
}

func AjaxError(code int, msg string, args ...string) Meta {
	return New("Error").
		WithAttr("code", strconv.Itoa(code)).
		WithAttr("message", msg).
		WithAttr(args...)
}

func ErrorFor(m Meta, err error, args ...string) Meta {
	ret := Error(err.Error(), "kind", m.Kind())
	if m.HasTag("method") {
		ret = ret.WithTag("method", m.Tag("method"))
	}
	return ret.WithAttr(args...) //.WithSub("for", m)
}
