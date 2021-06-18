package states

import (
	"context"
	"github.com/zilm13/zrnt/eth2/beacon/common"
	"github.com/protolambda/ztyp/tree"
	"sync"
	"sync/atomic"
)

type MemDB struct {
	// beacon.Root -> tree.Node (backing of BeaconStateView)
	data        sync.Map
	removalLock sync.Mutex
	stats       DBStats
	spec        *common.Spec
}

func NewMemDB(spec *common.Spec) *MemDB {
	return &MemDB{spec: spec}
}

func (db *MemDB) Store(ctx context.Context, state common.BeaconState) (exists bool, err error) {
	// Released when the block is removed from the DB
	root := state.HashTreeRoot(tree.GetHashFn())
	_, loaded := db.data.LoadOrStore(root, state)
	if !loaded {
		atomic.AddInt64(&db.stats.Count, 1)
		db.stats.LastWrite = root
	}
	return loaded, nil
}

func (db *MemDB) Get(ctx context.Context, root common.Root) (state common.BeaconState, exists bool, err error) {
	dat, ok := db.data.Load(root)
	if !ok {
		return nil, false, nil
	}
	exists = true
	state, ok = dat.(common.BeaconState)
	if !ok {
		panic("in-memory db was corrupted with unexpected state type")
	}
	return
}

func (db *MemDB) Remove(root common.Root) (exists bool, err error) {
	db.removalLock.Lock()
	defer db.removalLock.Unlock()
	_, ok := db.data.Load(root)
	if ok {
		atomic.AddInt64(&db.stats.Count, -1)
	}
	db.data.Delete(root)
	return ok, nil
}

func (db *MemDB) Stats() DBStats {
	// return a copy (struct is small and has no pointers)
	return db.stats
}

func (db *MemDB) List() (out []common.Root) {
	out = make([]common.Root, 0, db.stats.Count)
	db.data.Range(func(key, value interface{}) bool {
		id := key.(common.Root)
		out = append(out, id)
		return true
	})
	return out
}

func (db *MemDB) Path() string {
	return ""
}

func (db *MemDB) Spec() *common.Spec {
	return db.spec
}
