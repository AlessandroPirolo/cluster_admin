package vmmanager

import "raft/admin/src/internal/cluster"

type VMManager interface {
  Init() cluster.Cluster  
  Terminate() bool
  AddNode() cluster.Cluster
  RemoveNode(IP string)
}

func NewClusterManager() VMManager {
  return &vMManagerImpl{
    sources: *newSources(), 
  }
}
