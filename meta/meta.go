// Copyright (c) 2022 Jing-Ying Chen. MIT License. See https://github.com/jyrobin/mp

package meta

type Tag struct {
	Name  string
	Value string
}

type Meta interface {
	Kind() string
	Tag(name string) string
	HasTag(name string) bool
	WithTag(name, value string, rest ...string) Meta
	WithSub(name string, sub Meta) Meta
	HasSub(name string) bool
	Sub(name string) Meta
}

type meta struct {
	kind string
	tags map[string]string
	subs map[string]Meta
}

func New(kind string) *meta {
	return &meta{kind: kind}
}

func (m *meta) Kind() string {
	return m.kind
}

func (m *meta) Tag(name string) string {
	return m.tags[name]
}

func (m *meta) HasTag(name string) bool {
	_, ok := m.tags[name]
	return ok
}

func (m *meta) WithTag(name, value string, rest ...string) Meta {
	restn := len(rest) / 2
	tags := copyStrMap(m.tags, 1+restn)
	tags[name] = value
	for i := 0; i < restn; i++ {
		tags[rest[2*i]] = rest[2*i+1]
	}
	return &meta{kind: m.kind, tags: tags, subs: m.subs}
}

func copyStrMap(src map[string]string, more int) map[string]string {
	ret := make(map[string]string, len(src)+more)
	for k, v := range src {
		ret[k] = v
	}
	return ret
}

func (m *meta) Sub(name string) Meta {
	return m.subs[name]
}

func (m *meta) HasSub(name string) bool {
	_, ok := m.subs[name]
	return ok
}

func (m *meta) WithSub(name string, sub Meta) Meta {
	ch, ok := m.subs[name]
	if sub == nil || ok && sub == ch {
		return m
	}

	more := 1
	if ok {
		more = 0
	}

	subs := make(map[string]Meta, len(m.subs)+more)
	for k, v := range m.subs {
		subs[k] = v
	}
	return &meta{kind: m.kind, tags: m.tags, subs: subs}
}
