package mp

import (
	"context"
	"reflect"
	"strings"
	"unicode"
)

type Actor interface {
	Meta() Meta
	Process(ctx context.Context, m Meta, opts ...Meta) (Meta, error)
}

type BaseActor struct {
	meta Meta
}

func NewBaseActor(m Meta) BaseActor {
	return BaseActor{m}
}

func (a BaseActor) Meta() Meta {
	return a.meta
}

func SimpleActor(kind, mthd string, tags ...string) BaseActor {
	m := New("Actor", mthd).WithTag("target", kind)
	return BaseActor{m}
}

func SimpleLister(kind string, tags ...string) BaseActor {
	return SimpleActor(kind, "list")
}

func SimpleFinder(kind string, tags ...string) BaseActor {
	return SimpleActor(kind, "find")
}

func SimpleCreator(kind string, tags ...string) BaseActor {
	return SimpleActor(kind, "create")
}

func SimpleMaker(kind string, tags ...string) BaseActor {
	return SimpleActor(kind, "make")
}

func SimpleRemover(kind string, tags ...string) BaseActor {
	return SimpleActor(kind, "remove")
}

func caseFirstLetter(s string, upper bool) string {
	if len(s) == 0 {
		return s
	}
	r := []rune(s)
	if upper && unicode.IsUpper(r[0]) || !upper && unicode.IsLower(r[0]) {
		return s
	}
	if upper {
		r[0] = unicode.ToUpper(r[0])
	} else {
		r[0] = unicode.ToLower(r[0])
	}
	return string(r)
}

func toUpper(s string) string {
	return caseFirstLetter(s, true)
}
func toLower(s string) string {
	return caseFirstLetter(s, false)
}

func ReflectActorList(val interface{}, kind string, prefixAndNames ...string) []Actor {
	t := reflect.TypeOf(val)

	var prefix string
	var names []string
	argn := len(prefixAndNames)
	if argn > 0 {
		prefix = strings.TrimSpace(prefixAndNames[0])
		names = prefixAndNames[1:]
	}
	if prefix == "" {
		prefix = "Mpi"
	}

	actorMap := map[string]bool{}
	actors := []Actor{}
	for _, name := range names {
		name = toUpper(strings.TrimSpace(name))
		if name != "" && !actorMap[name] {
			if mthd, ok := t.MethodByName(name); ok {
				actorMap[name] = true
				actors = append(actors, createActor(val, t, mthd, kind, name))
			}
		}
	}

	for i := 0; i < t.NumMethod(); i++ {
		mthd := t.Method(i)
		if strings.HasPrefix(mthd.Name, prefix) {
			name := toUpper(mthd.Name[len(prefix):])
			if name != "" && !actorMap[name] {
				actors = append(actors, createActor(val, t, mthd, kind, name))
			}
		}
	}

	return actors
}

var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()
var metaType = reflect.TypeOf((*Meta)(nil)).Elem()
var errorType = reflect.TypeOf((*error)(nil)).Elem()

func createActor(val interface{}, vt reflect.Type, mthd reflect.Method, kind string, mthdName string) Actor {
	fullName := vt.Name() + "." + mthd.Name

	method := reflect.ValueOf(val).MethodByName(mthd.Name)
	if !method.IsValid() {
		panic(fullName + ": invalid method")
	}

	t := method.Type()
	if !t.IsVariadic() {
		panic(fullName + ": not variadic")
	}
	if t.NumIn() != 3 {
		panic(fullName + ": not 3 input params")
	}
	if t.In(0) != contextType {
		panic(fullName + ": first parameter not context.Context")
	}
	if t.In(1) != metaType {
		panic(fullName + ": second parameter not mp.Meta")
	}
	if t.In(2).Kind() != reflect.Slice || t.In(2).Elem() != metaType {
		panic(fullName + ": third parameter not variadic mp.Meta")
	}

	if t.NumOut() != 2 {
		panic(fullName + ": not 2 return values")
	}
	if t.Out(0) != metaType {
		panic(fullName + ": first return value not mp.Meta")
	}
	if t.Out(1) != errorType {
		panic(fullName + ": second return value not error")
	}

	return &reflectActor{
		val,
		method,
		New("Actor", toLower(mthdName)).WithTag("target", kind),
	}
}

type reflectActor struct {
	val    interface{}
	method reflect.Value
	meta   Meta
}

func (a *reflectActor) Meta() Meta {
	return a.meta
}

func (a *reflectActor) Process(ctx context.Context, m Meta, opts ...Meta) (Meta, error) {
	inputs := make([]reflect.Value, 2+len(opts))
	inputs[0] = reflect.ValueOf(ctx)
	inputs[1] = reflect.ValueOf(m)
	for i, opt := range opts {
		inputs[i+2] = reflect.ValueOf(opt)
	}

	out := a.method.Call(inputs)
	ret := out[0].Interface().(Meta)
	er := out[1].Interface()
	var err error
	if er != nil {
		err = er.(error)
	}
	return ret, err
}
