package mp

import (
	"context"
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
