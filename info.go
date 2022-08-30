package mp

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	UtcDateFormat = "2006-01-02"
	UtcTimeFormat = "2006-01-02T15:04:05.000Z"
)

var truth = map[string]bool{
	"1": true, "true": true, "yes": true, "on": true,
	"0": false, "false": false, "no": false, "off": false,
}

type Info interface {
	IsNil() bool

	Kind() string
	Ns() string
	Method() string

	Gid() string
	Label() string

	Tag(name string) string
	HasTag(args ...string) bool
	HasTags(tags map[string]string) bool
	TagNames() []string
	TagCount() int
	TagMap() map[string]string // cloned
	Is(kind, ns string, tags ...string) bool

	Attr(name string) string
	IntAttr(name string) (int, error)
	IntAttrOr(name string, otherwise int) int
	BoolAttr(name string) (bool, error)
	IsTrueAttr(name string) bool
	IsFalseAttr(name string) bool
	FloatAttr(name string) (float64, error)
	DateAttr(name string, loc ...*time.Location) (time.Time, error)
	TimeAttr(name, layout string, loc ...*time.Location) (time.Time, error)
	UtcAttr(name string) time.Time
	IntsAttr(name, sep string) ([]int, error)
	HasAttr(args ...string) bool
	HasAttrs(attrs map[string]string) bool
	AttrNames() []string
	AttrCount() int
	AttrMap(skips ...string) map[string]string // cloned

	Payload() string
}

// base impl

type info struct {
	kind    string
	mthd    string
	ns      string
	gid     string
	tags    map[string]string
	attrs   map[string]string
	payload string
}

func (m info) Kind() string {
	return m.kind
}

func (m info) Ns() string {
	return m.ns
}

func (m info) Method() string {
	return m.mthd
}

func (m info) Gid() string {
	return m.gid
}

func (m info) Label() string {
	if m.mthd != "" {
		return fmt.Sprintf("%s:%s.%s:%s", m.ns, m.kind, m.mthd, m.gid)
	} else {
		return fmt.Sprintf("%s:%s:%s", m.ns, m.kind, m.gid)
	}
}

//func (m info) IsNil() bool {
//	return m.kind == ""
//}

func (m info) Tag(name string) string {
	return m.tags[name]
}

func (m info) HasTag(args ...string) bool {
	return hasValue(m.tags, args)
}

func hasValue(vals map[string]string, args []string) bool {
	argn := len(args)
	for i, n := 0, argn/2; i < n; i++ {
		if val, ok := vals[args[2*i]]; !ok || val != args[2*i+1] {
			return false
		}
	}
	ok := true
	if argn%2 != 0 {
		_, ok = vals[args[argn-1]]
	}
	return ok
}

func (m info) HasTags(ts map[string]string) bool {
	return hasValues(m.tags, ts)
}

func hasValues(vals, ts map[string]string) bool {
	for k, v := range ts {
		if val, ok := vals[k]; !ok || val != v {
			return false
		}
	}
	return true
}

func (m info) TagCount() int {
	return len(m.tags)
}

func (m info) TagNames() []string {
	return mapKeys(m.tags)
}

func mapKeys(v map[string]string) []string {
	keys := make([]string, 0, len(v))
	for key, _ := range v {
		keys = append(keys, key)
	}
	return keys
}

func copyStrMap(src map[string]string, more int) map[string]string {
	ret := make(map[string]string, len(src)+more)
	for k, v := range src {
		ret[k] = v
	}
	return ret
}

func (m info) TagMap() map[string]string {
	ret := make(map[string]string, len(m.tags))
	for k, v := range m.tags {
		ret[k] = v
	}
	return ret
}

func (m info) Is(kind, ns string, tags ...string) bool {
	return m.Kind() == kind && m.Ns() == ns && m.HasTag(tags...)
}

func (m info) Attr(name string) string {
	return m.attrs[name]
}
func (m info) IntAttr(name string) (int, error) {
	return strconv.Atoi(m.attrs[name])
}
func (m info) IntAttrOr(name string, otherwise int) int {
	if val, ok := m.attrs[name]; ok {
		if ret, err := strconv.Atoi(val); err == nil {
			return ret
		}
	}
	return otherwise
}

func (m info) BoolAttr(name string) (bool, error) {
	val := m.Attr(name)
	if b, ok := truth[strings.ToLower(val)]; ok {
		return b, nil
	}
	return false, fmt.Errorf("Invalid bool %s", val)
}
func (m info) IsTrueAttr(name string) bool {
	val := strings.ToLower(m.Attr(name))
	b, ok := truth[val]
	return ok && b
}
func (m info) IsFalseAttr(name string) bool {
	val := strings.ToLower(m.Attr(name))
	b, ok := truth[val]
	return ok && !b
}

func (m info) FloatAttr(name string) (float64, error) {
	return strconv.ParseFloat(m.attrs[name], 64)
}
func (m info) DateAttr(name string, loc ...*time.Location) (time.Time, error) {
	if len(loc) > 0 && loc[0] != nil {
		return time.ParseInLocation(UtcDateFormat, m.attrs[name], loc[0])
	}
	return time.Parse(UtcDateFormat, m.attrs[name]) // UTC
}
func (m info) TimeAttr(name, layout string, loc ...*time.Location) (time.Time, error) {
	if len(loc) > 0 && loc[0] != nil {
		return time.ParseInLocation(layout, m.attrs[name], loc[0])
	}
	return time.Parse(layout, m.attrs[name])
}
func (m info) UtcAttr(name string) time.Time {
	if ret, err := time.Parse(UtcTimeFormat, m.attrs[name]); err == nil {
		return ret
	}
	return time.Time{}
}

func (m info) IntsAttr(name, sep string) ([]int, error) {
	v := m.attrs[name]
	if v == "" {
		return []int{}, nil
	}
	words := strings.Split(v, sep)
	ret := make([]int, len(words))
	var err error
	for i, word := range words {
		if ret[i], err = strconv.Atoi(word); err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func (m info) HasAttr(args ...string) bool {
	return hasValue(m.attrs, args)
}

func (m info) HasAttrs(ts map[string]string) bool {
	return hasValues(m.attrs, ts)
}

func (m info) AttrCount() int {
	return len(m.attrs)
}

func (m info) AttrNames() []string {
	return mapKeys(m.attrs)
}

func (m info) AttrMap(skips ...string) map[string]string {
	b := len(skips) == 0
	ret := make(map[string]string, len(m.attrs))
	for k, v := range m.attrs {
		if b || !contains(skips, k) {
			ret[k] = v
		}
	}
	return ret
}

func contains(arr []string, val string) bool {
	for i := range arr {
		if arr[i] == val {
			return true
		}
	}
	return false
}

// payload

func (m info) Payload() string {
	return m.payload
}

// utils

func CopyTags(m Info) map[string]string {
	ret := make(map[string]string, m.TagCount())
	for _, name := range m.TagNames() {
		ret[name] = m.Tag(name)
	}
	return ret
}

func CopyAttrs(m Info) map[string]string {
	ret := make(map[string]string, m.AttrCount())
	for _, name := range m.AttrNames() {
		ret[name] = m.Attr(name)
	}
	return ret
}
