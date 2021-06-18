package pool

import (
	"github.com/zilm13/zrnt/eth2/beacon/common"
	"github.com/zilm13/zrnt/eth2/beacon/phase0"
	"github.com/protolambda/ztyp/tree"
	"sync"
)

type VersionedRoot struct {
	Root    common.Root
	Version common.Version
}

type AttesterSlashingPool struct {
	sync.RWMutex
	spec      *common.Spec
	slashings map[VersionedRoot]*phase0.AttesterSlashing
}

func NewAttesterSlashingPool(spec *common.Spec) *AttesterSlashingPool {
	return &AttesterSlashingPool{
		spec:      spec,
		slashings: make(map[VersionedRoot]*phase0.AttesterSlashing),
	}
}

// This does not filter slashings that are a subset of other slashings.
// The pool merely collects them. Make sure to protect against spam elsewhere as a caller.
func (asp *AttesterSlashingPool) AddAttesterSlashing(sl *phase0.AttesterSlashing, version common.Version) (exists bool) {
	root := sl.HashTreeRoot(asp.spec, tree.GetHashFn())
	asp.Lock()
	defer asp.Unlock()
	key := VersionedRoot{Root: root, Version: version}
	if _, ok := asp.slashings[key]; ok {
		return true
	}
	asp.slashings[key] = sl
	return false
}

func (asp *AttesterSlashingPool) All() []*phase0.AttesterSlashing {
	asp.RLock()
	defer asp.RUnlock()
	out := make([]*phase0.AttesterSlashing, 0, len(asp.slashings))
	for _, a := range asp.slashings {
		out = append(out, a)
	}
	return out
}

// Pack n slashings, removes the slashings from the pool. A reward estimator is used to pick the best slashings.
// Slashings with negative rewards will not be packed.
func (asp *AttesterSlashingPool) Pack(estReward func(sl *phase0.AttesterSlashing) int, n uint) []*phase0.AttesterSlashing {
	// TODO
	return nil
}
