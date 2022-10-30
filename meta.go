// Copyright (c) 2022 Jing-Ying Chen. MIT License. See https://github.com/jyrobin/mp

package mp

import (
	"encoding/json"
	"fmt"
)

var Nil Meta = &meta{}

type Meta interface {
	Info

	IsError() bool
	IsValid() bool
	Error() string
	ErrorCode(otherwise int) int

	WithMethod(mthd string) Meta
	WithGid(gid string) Meta

	WithTag(name, value string, rest ...string) Meta
	WithTags(map[string]string) Meta
	WithAttr(args ...string) Meta
	WithNonEmptyAttr(args ...string) Meta
	WithAttrs(map[string]string) Meta

	WithPayload(data string) Meta

	WithSub(name string, sub Meta) Meta
	WithSubs(map[string]Meta) Meta
	HasSub(name string) bool
	Sub(name string) Meta
	SubNames() []string
	SubCount() int

	WithRel(name string, rel Meta) Meta
	WithRels(map[string]Meta) Meta
	HasRel(name string) bool
	Rel(name string) Meta
	RelNames() []string
	RelCount() int

	WithList(list []Meta, trims ...bool) Meta
	List() []Meta
	ListItem(index int) Meta

	Generalizes(m Meta) bool
	Specializes(m Meta) bool

	Walk(v Visitor)

	Json(opts ...string) string
}

type Visitor interface {
	BeginMeta(m Meta, kind, method, ns, gid string) // add uppers later
	OnTag(m Meta, name, value string)
	OnAttr(m Meta, name, value string)
	OnSub(m Meta, name string, sub Meta)
	OnRel(m Meta, name string, rel Meta)
	OnListItem(m Meta, index int, item Meta)
	EndMeta(m Meta)
}

type meta struct {
	info
	subs map[string]Meta // assert subs has no nil
	rels map[string]Meta // assert rels has no nil
	list []Meta          // assert list has no nil
}

func New(kind string, args ...string) Meta {
	var mthd, ns, gid string
	switch len(args) {
	case 0:
	case 1:
		mthd = args[0]
	case 2:
		mthd, ns = args[0], args[1]
	default:
		mthd, ns, gid = args[0], args[1], args[2]
	}
	return &meta{info{kind, mthd, ns, gid, nil, nil, ""}, nil, nil, nil}
}

func IsNil(m Meta) bool {
	return m == nil || m.IsNil()
}

func (m *meta) IsNil() bool {
	return m == nil || m.kind == ""
}

func (m *meta) IsError() bool {
	return m != nil && m.kind == "Error"
}

func (m *meta) IsValid() bool {
	return m != nil && m.kind != "" && m.kind != "Error"
}
func (m *meta) Error() string {
	if m.IsNil() {
		return "Nil"
	} else if m.IsError() {
		ret := m.Attr("message")
		if ret == "" {
			ret = "Error"
		}
		code := m.Attr("code")
		if code != "" {
			ret = fmt.Sprintf("%s (CODE %s)", ret, code)
		}
		return ret
	}
	return ""
}

func (m *meta) ErrorCode(otherwise int) int {
	if m.IsError() && m.HasAttr("code") {
		if code, err := m.IntAttr("code"); err == nil {
			return code
		}
	}
	return otherwise
}

func (m *meta) WithMethod(mthd string) Meta {
	return &meta{info{m.kind, mthd, m.ns, m.gid, m.tags, m.attrs, m.payload}, m.subs, m.rels, m.list}
}

func (m *meta) WithGid(gid string) Meta {
	return &meta{info{m.kind, m.mthd, m.ns, gid, m.tags, m.attrs, m.payload}, m.subs, m.rels, m.list} // better to enumerate all
}

func (m *meta) WithTag(name, value string, rest ...string) Meta {
	restn := len(rest) / 2
	tags := copyStrMap(m.tags, 1+restn)
	tags[name] = value
	for i := 0; i < restn; i++ {
		tags[rest[2*i]] = rest[2*i+1]
	}
	return &meta{info{m.kind, m.mthd, m.ns, m.gid, tags, m.attrs, m.payload}, m.subs, m.rels, m.list} // better to enumerate all
}

func (m *meta) WithTags(ts map[string]string) Meta {
	if len(ts) == 0 {
		return m
	}

	tags := copyStrMap(m.tags, len(ts))
	for k, v := range ts {
		tags[k] = v
	}
	return &meta{info{m.kind, m.mthd, m.ns, m.gid, tags, m.attrs, m.payload}, m.subs, m.rels, m.list}
}

func (m *meta) WithAttr(args ...string) Meta {
	return m.withAttr(false, args)
}

func (m *meta) WithNonEmptyAttr(args ...string) Meta {
	return m.withAttr(true, args)
}

func (m *meta) withAttr(skip bool, args []string) Meta {
	argn := len(args) / 2
	changed := 0
	for i := 0; i < argn; i++ {
		if !skip || args[2*i+1] != "" {
			changed += 1
		}
	}
	if changed == 0 {
		return m
	}

	attrs := copyStrMap(m.attrs, changed)
	for i := 0; i < argn; i++ {
		if val := args[2*i+1]; !skip || val != "" {
			attrs[args[2*i]] = val
		}
	}
	return &meta{info{m.kind, m.mthd, m.ns, m.gid, m.tags, attrs, m.payload}, m.subs, m.rels, m.list}
}

func (m *meta) WithAttrs(ts map[string]string) Meta {
	if len(ts) == 0 {
		return m
	}

	attrs := copyStrMap(m.attrs, len(ts))
	for k, v := range ts {
		attrs[k] = v
	}
	return &meta{info{m.kind, m.mthd, m.ns, m.gid, m.tags, attrs, m.payload}, m.subs, m.rels, m.list}
}

// payload

func (m *meta) WithPayload(payload string) Meta {
	if payload == m.payload {
		return m
	}
	return &meta{info{m.kind, m.mthd, m.ns, m.gid, m.tags, m.attrs, payload}, m.subs, m.rels, m.list}
}

// subs

func (m *meta) Sub(name string) Meta {
	if ret, ok := m.subs[name]; ok {
		return ret
	}
	return Nil
}

func (m *meta) HasSub(name string) bool {
	sub, ok := m.subs[name]
	return ok && !sub.IsNil() // make sure subs no nil
}

func (m *meta) WithSub(name string, sub Meta) Meta {
	if subs, changed := withSub(m.subs, name, sub); changed {
		return &meta{info{m.kind, m.mthd, m.ns, m.gid, m.tags, m.attrs, m.payload}, subs, m.rels, m.list}
	}
	return m
}

func withSub(subs map[string]Meta, name string, sub Meta) (map[string]Meta, bool) {
	ch, ok := subs[name]
	if sub == nil || ok && sub == ch {
		return subs, false
	}

	more := 1
	if ok {
		more = 0
	}

	newSubs := make(map[string]Meta, len(subs)+more)
	for k, v := range subs {
		newSubs[k] = v
	}
	newSubs[name] = sub
	return newSubs, true
}

func (m *meta) WithSubs(subs map[string]Meta) Meta {
	newSubs := m.subs
	mod, changed := false, false
	for k, v := range subs {
		if v != nil { // defensive
			newSubs, changed = withSub(newSubs, k, v)
			mod = mod || changed
		}
	}

	if !mod {
		return m
	}
	return &meta{info{m.kind, m.mthd, m.ns, m.gid, m.tags, m.attrs, m.payload}, newSubs, m.rels, m.list}
}

func (m *meta) SubCount() int {
	return len(m.subs)
}

func (m *meta) SubNames() []string {
	names := make([]string, 0, len(m.subs))
	for name := range m.subs {
		names = append(names, name)
	}
	return names
}

func (m *meta) WithList(list []Meta, trims ...bool) Meta {
	if len(list) == 0 && len(m.list) == 0 {
		return m
	}

	trim := true // IMPORTANT: default trim nil
	if len(trims) > 0 {
		trim = trims[0]
	}
	newList := make([]Meta, 0, len(list)) // for not changing size
	for _, item := range list {
		if item != nil && !item.IsNil() {
			newList = append(newList, item)
		} else if !trim {
			newList = append(newList, Nil)
		}
	}

	if len(newList) == 0 && len(m.list) == 0 {
		return m
	}

	return &meta{info{m.kind, m.mthd, m.ns, m.gid, m.tags, m.attrs, m.payload}, m.subs, m.rels, newList}
}

// rels

func (m *meta) Rel(name string) Meta {
	if ret, ok := m.rels[name]; ok {
		return ret
	}
	return Nil
}

func (m *meta) HasRel(name string) bool {
	rel, ok := m.rels[name]
	return ok && !rel.IsNil()
}

func (m *meta) WithRel(name string, rel Meta) Meta {
	if rels, changed := withSub(m.rels, name, rel); changed {
		return &meta{info{m.kind, m.mthd, m.ns, m.gid, m.tags, m.attrs, m.payload}, m.subs, rels, m.list}
	}
	return m
}

func (m *meta) WithRels(rels map[string]Meta) Meta {
	newRels := m.rels
	mod, changed := false, false
	for k, v := range rels {
		if v != nil {
			newRels, changed = withSub(newRels, k, v)
			mod = mod || changed
		}
	}

	if !mod {
		return m
	}
	return &meta{info{m.kind, m.mthd, m.ns, m.gid, m.tags, m.attrs, m.payload}, m.subs, newRels, m.list}
}

func (m *meta) RelCount() int {
	return len(m.rels)
}

func (m *meta) RelNames() []string {
	names := make([]string, 0, len(m.rels))
	for name := range m.rels {
		names = append(names, name)
	}
	return names
}

// list

func (m *meta) List() []Meta {
	return m.list
}

func (m *meta) ListItem(idx int) Meta {
	if idx >= 0 && idx < len(m.list) {
		return m.list[idx]
	}
	return nil
}

func (m *meta) Generalizes(n Meta) bool {
	if m.Kind() == n.Kind() {
		for k, v := range m.tags {
			if !n.HasTag(k, v) {
				return false
			}
		}
	}
	return true
}

func (m *meta) Specializes(n Meta) bool {
	return n.Generalizes(m)
}

// traverse

func (m *meta) Walk(v Visitor) {
	Walk(m, v)
}

func Walk(m Meta, v Visitor) {
	v.BeginMeta(m, m.Kind(), m.Method(), m.Ns(), m.Gid())
	for _, name := range m.TagNames() {
		v.OnTag(m, name, m.Tag(name))
	}
	for _, name := range m.SubNames() {
		v.OnSub(m, name, m.Sub(name))
		Walk(m.Sub(name), v)
	}
	for _, name := range m.RelNames() {
		v.OnRel(m, name, m.Rel(name))
		Walk(m.Rel(name), v)
	}

	for idx, item := range m.List() {
		v.OnListItem(m, idx, item)
		Walk(item, v)
	}

	v.EndMeta(m)
}

func (m *meta) MarshalJSON() ([]byte, error) {
	return json.Marshal(MetaToJson(m))
}

func (m *meta) Json(opts ...string) string {
	return MetaToJson(m).Json(opts...)
}

// utils

func First(ms []Meta) Meta {
	ret := Nil
	if len(ms) > 0 && ms[0] != nil {
		ret = ms[0]
	}
	return ret
}

func FirstAttr(ms []Meta, name string) string {
	return First(ms).Attr(name)
}

func FirstIs(ms []Meta, kind, ns string, tags ...string) Meta {
	ret := First(ms)
	if !ret.IsNil() && ret.Is(kind, ns, tags...) {
		return ret
	}
	return Nil
}
