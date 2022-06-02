// Copyright (c) 2022 Jing-Ying Chen. MIT License. See https://github.com/jyrobin/mp
package mp

import (
	"context"
	"fmt"
)

type Mpi interface {
	IsNil() bool

	Call(ctx context.Context, method string, m Meta, opts ...Meta) (Meta, error)
	//List(ctx context.Context, m Meta, opts ...Meta) (Meta, error)
	//Find(ctx context.Context, m Meta, opts ...Meta) (Meta, error)
	//Create(ctx context.Context, m Meta, opts ...Meta) (Meta, error)
	//Make(ctx context.Context, m Meta, opts ...Meta) (Meta, error)

	// Later...

	// find the Actor targeting First(m, filters...)
	//FirstActor(ctx context.Context, m Meta, filters ...Meta) (Actor, error)

	//GetActor(ctx context.Context, gid string) (Actor, error)

	//Process(ctx context.Context, gid string, m Meta) (Meta, error)

	// same as FirstActor(m, filters...) then call it's Process(m)
	//Do(ctx context.Context, m Meta, filters ...Meta) (Meta, error)
}

type Domain interface {
	Mpi

	Meta() Meta

	Parent() Domain
	Subs() []Domain
	Sub(name string) Domain
	WithParent(parent Domain) Domain
	WithSubs(subs ...Domain) Domain

	Actors() []Actor
	WithActors(actors ...Actor) Domain

	Indexer() Indexer
}

type DomainConfig struct {
	IndexerFn func(Domain) Indexer
	InfoFn    func(Domain) Meta
}

type domain struct {
	cfg DomainConfig

	meta    Meta
	parent  Domain
	subs    []Domain
	actors  []Actor
	indexer Indexer
}

func SimpleDomain(m Meta, cfgs ...DomainConfig) *domain {
	var cfg DomainConfig
	if len(cfgs) > 0 {
		cfg = cfgs[0]
	}
	return &domain{cfg, m, nil, nil, nil, nil} // make sure (also below) indexer init with literal nil
}

func (dom *domain) IsNil() bool {
	return dom == nil || dom.meta == nil || dom.meta.IsNil()
}

func (dom *domain) Parent() Domain {
	return dom.parent
}

func (dom *domain) Root() Domain {
	var root Domain = dom
	for root.Parent() != nil {
		root = root.Parent()
	}
	return root
}

func (dom *domain) Meta() Meta {
	return dom.meta
}

func (dom *domain) Sub(name string) Domain {
	for _, sub := range dom.subs {
		if sub.Meta().Ns() == name {
			return sub
		}
	}
	return nil
}

func (dom *domain) Subs() []Domain {
	return dom.subs
}

func (dom *domain) WithParent(parent Domain) Domain {
	return &domain{dom.cfg, dom.meta, parent, dom.subs, dom.actors, nil}
}

func (dom *domain) WithSubs(subs ...Domain) Domain {
	newSubs := make([]Domain, len(subs))
	for i, sub := range subs {
		newSubs[i] = sub.WithParent(dom)
	}
	return &domain{dom.cfg, dom.meta, dom.parent, newSubs, dom.actors, nil}
}

func (dom *domain) Actors() []Actor {
	return dom.actors
}

func (dom *domain) WithActors(actors ...Actor) Domain {
	return &domain{dom.cfg, dom.meta, dom.parent, dom.subs, actors, nil}
}

func (dom *domain) Call(ctx context.Context, mthd string, m Meta, opts ...Meta) (Meta, error) {
	if actor := dom.Indexer().ActorWithMethod(m.Kind(), mthd); actor != nil {
		ret, err := actor.Process(ctx, m, opts...)
		if ret == nil {
			ret = Nil
		}
		return ret, err
	}
	return Nil, fmt.Errorf("Actor %s for %s not found", mthd, m.Kind())
}

func (dom *domain) List(ctx context.Context, m Meta, filters ...Meta) (Meta, error) {
	return dom.Call(ctx, "list", m, filters...)
}

func (dom *domain) Find(ctx context.Context, m Meta, filters ...Meta) (Meta, error) {
	return dom.Call(ctx, "find", m, filters...)
}

func (dom *domain) Create(ctx context.Context, m Meta, filters ...Meta) (Meta, error) {
	return dom.Call(ctx, "create", m, filters...)
}

func (dom *domain) Make(ctx context.Context, m Meta, filters ...Meta) (Meta, error) {
	return dom.Call(ctx, "make", m, filters...)
}

/* Later
func (dom *domain) List(ctx context.Context, m Meta, filters ...Meta) (Meta, error) {
	// TODO: hack for the time being; filters later
	if actor := dom.Indexer().ActorWithMethod(m.Kind(), "list"); actor != nil {
		return actor.Process(ctx, m)
	}
	return nil, fmt.Errorf("Lister for %s not found", m.Kind())
}

func (dom *domain) GetActor(ctx context.Context, gid string) (Actor, error) {
	return dom.Indexer().ActorWithGid(gid), nil
}

func (dom *domain) FirstActor(ctx context.Context, m Meta, filters ...Meta) (Actor, error) {
	actors := dom.Indexer().ActorsWithKind(m.Kind())
	for _, actor := range actors {
		if actor.Meta().Specializes(m) { // filters later
			return actor, nil
		}
	}
	return nil, fmt.Errorf("Actor for %s not found", m.Kind())
}

func (dom *domain) Process(ctx context.Context, gid string, m Meta) (Meta, error) {
	if actor, err := dom.GetActor(ctx, gid); err != nil {
		return nil, err
	} else {
		return actor.Process(ctx, m)
	}
}

func (dom *domain) Do(ctx context.Context, m Meta, filters ...Meta) (Meta, error) {
	if actor, err := dom.FirstActor(ctx, m, filters...); err != nil {
		return nil, err
	} else {
		return actor.Process(ctx, m)
	}
}
*/

// reentrant-able as domain itself is constant
func (dom *domain) Indexer() Indexer {
	indexer := dom.indexer
	if indexer == nil { // make sure dom initiated with literal nil
		if dom.cfg.IndexerFn != nil {
			indexer = dom.cfg.IndexerFn(dom)
		} else {
			indexer = SimpleIndexer(dom)
		}
		dom.indexer = indexer
	}
	return indexer
}
