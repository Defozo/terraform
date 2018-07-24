package terraform

import (
	"github.com/hashicorp/terraform/states"
)

// NodeDestroyableDataResource represents a resource that is "destroyable":
// it is ready to be destroyed.
type NodeDestroyableDataResource struct {
	*NodeAbstractResourceInstance
}

// GraphNodeEvalable
func (n *NodeDestroyableDataResource) EvalTree() EvalNode {
	addr := n.ResourceInstanceAddr()

	// stateId is the ID to put into the state
	stateId := NewLegacyResourceInstanceAddress(addr).stateId()

	// Just destroy it.
	var state *states.ResourceInstanceObject
	return &EvalWriteState{
		Addr:  addr.Resource,
		State: &state, // state is nil here
	}
}
