// Copyright (c) 2022 Jing-Ying Chen. MIT License. See https://github.com/jyrobin/mp
package mp

import (
	"encoding/json"
)

type MetaJson struct {
	Kind    string              `json:"kind,omitempty"`
	Method  string              `json:"method,omitempty"`
	Ns      string              `json:"ns,omitempty"`
	Gid     string              `json:"gid,omitempty"`
	Tags    map[string]string   `json:"tags,omitempty"`
	Attrs   map[string]string   `json:"attrs,omitempty"`
	Payload string              `json:"payload,omitempty"`
	Subs    map[string]MetaJson `json:"subs,omitempty"`
	Rels    map[string]MetaJson `json:"rels,omitempty"`
	List    []MetaJson          `json:"list,omitempty"`
}

func ParseMeta(buf []byte) (Meta, error) {
	var mj MetaJson
	if err := json.Unmarshal(buf, &mj); err != nil {
		return Nil, err
	}
	return JsonToMeta(mj), nil
}

/*func ParseMetas(buf []byte) ([]Meta, error) {
	var mjs []MetaJson
	if err := json.Unmarshal(buf, &mjs); err != nil {
		return nil, err
	}
	return JsonsToMetas(mjs), nil
}
*/

func MetaToJson(m Meta) MetaJson {
	if m == nil || m.IsNil() {
		return MetaJson{}
	}

	return MetaJson{
		Kind:    m.Kind(),
		Method:  m.Method(),
		Ns:      m.Ns(),
		Gid:     m.Gid(),
		Tags:    CopyTags(m),
		Attrs:   CopyAttrs(m),
		Payload: m.Payload(),
		Subs:    subMetaJsons(m),
		Rels:    relMetaJsons(m),
		List:    listMetaJsons(m),
	}
}

func MetasToJsons(ms []Meta) []MetaJson {
	items := make([]MetaJson, len(ms))
	for i, m := range ms {
		items[i] = MetaToJson(m)
	}
	return items
}

func subMetaJsons(m Meta) map[string]MetaJson {
	ret := make(map[string]MetaJson, m.SubCount())
	for _, name := range m.SubNames() {
		ret[name] = MetaToJson(m.Sub(name))
	}
	return ret
}

func relMetaJsons(m Meta) map[string]MetaJson {
	ret := make(map[string]MetaJson, m.RelCount())
	for _, name := range m.RelNames() {
		ret[name] = MetaToJson(m.Rel(name))
	}
	return ret
}

func listMetaJsons(m Meta) []MetaJson {
	items := m.List()
	ret := make([]MetaJson, len(items))
	for idx, item := range items {
		ret[idx] = MetaToJson(item)
	}
	return ret
}

func JsonToMeta(mj MetaJson) Meta {
	if mj.IsNil() {
		return Nil
	}

	ret := New(mj.Kind, mj.Method, mj.Ns, mj.Gid).
		WithTags(mj.Tags).
		WithAttrs(mj.Attrs).
		WithPayload(mj.Payload)
	for name, subJson := range mj.Subs {
		ret = ret.WithSub(name, JsonToMeta(subJson))
	}
	for name, relJson := range mj.Rels {
		ret = ret.WithRel(name, JsonToMeta(relJson))
	}
	if len(mj.List) > 0 {
		ret = ret.WithList(JsonsToMetas(mj.List))
	}

	return ret
}

func JsonsToMetas(ml []MetaJson) []Meta {
	n := len(ml)
	items := make([]Meta, 0, n)
	for i := 0; i < n; i++ {
		items = append(items, JsonToMeta(ml[i]))
	}
	return items
}

func (mj MetaJson) IsNil() bool {
	return mj.Kind == ""
}

func (mj MetaJson) Meta() Meta {
	return JsonToMeta(mj)
}

func (mj MetaJson) String() string {
	return mj.Json("  ")
}
func (mj MetaJson) Json(opts ...string) string {
	var buf []byte
	switch len(opts) {
	case 0:
		buf, _ = json.Marshal(mj)
	case 1:
		buf, _ = json.MarshalIndent(mj, "", opts[0])
	default:
		buf, _ = json.MarshalIndent(mj, opts[0], opts[1])
	}
	return string(buf)
}

// MetaListJson

type MetaListJson struct {
	Meta MetaJson `json:"meta"`
	Gids []string `json:"gids,omitempty"`
}
