package mp

type Indexer interface {
	GidActorMap() map[string]Actor
	KindActorsMap() map[string][]Actor
	MethodActorMap() map[string]Actor
	ActorWithGid(gid string) Actor
	ActorsWithKind(kind string) []Actor
	ActorWithMethod(kind, method string) Actor
}

type indexer struct {
	dom        Domain
	kindActors map[string][]Actor
	gidActors  map[string]Actor
	mthdActors map[string]Actor
}

func SimpleIndexer(dom Domain) *indexer {
	gmap, kmap, mmap := map[string]Actor{}, map[string][]Actor{}, map[string]Actor{}
	for _, actor := range dom.Actors() {
		am := actor.Meta()
		if gid := am.Gid(); gid != "" {
			if _, ok := gmap[gid]; !ok { // no override...
				gmap[gid] = actor
			}
		}

		kind := am.Tag("target") // override target rel for now
		if kind == "" {
			if m := am.Rel("target"); m != nil && m.Kind() != "" {
				kind = m.Kind()
				if m.Ns() != "" {
					kind = m.Ns() + ":" + kind
				}
			}
		}
		if kind != "" {
			kmap[kind] = append(kmap[kind], actor)

			if mthd := am.Method(); mthd != "" {
				key := kind + "." + mthd
				if _, ok := mmap[key]; !ok { // no override...
					mmap[key] = actor
				}
			}
		}
	}

	for _, sub := range dom.Subs() {
		idx := sub.Indexer()
		gmap2, kmap2, mmap2 := idx.GidActorMap(), idx.KindActorsMap(), idx.MethodActorMap()
		for k, v := range gmap2 {
			if _, ok := gmap[k]; !ok {
				gmap[k] = v
			}
		}
		for k, v := range kmap2 {
			kmap[k] = append(kmap[k], v...)
		}
		for k, v := range mmap2 {
			if _, ok := mmap[k]; !ok {
				mmap[k] = v
			}
		}
	}
	return &indexer{dom, kmap, gmap, mmap}
}

func (idx *indexer) ActorWithGid(gid string) Actor {
	return idx.gidActors[gid]
}
func (idx *indexer) ActorsWithKind(kind string) []Actor {
	return idx.kindActors[kind]
}
func (idx *indexer) ActorWithMethod(kind, method string) Actor {
	return idx.mthdActors[kind+"."+method]
}

func (idx *indexer) GidActorMap() map[string]Actor {
	return idx.gidActors
}
func (idx *indexer) KindActorsMap() map[string][]Actor {
	return idx.kindActors
}
func (idx *indexer) MethodActorMap() map[string]Actor {
	return idx.mthdActors
}
