package cluster

import (
	"raft/admin/src/internal/rpcs"
	"raft/admin/src/internal/utility"
)

type Cluster interface {
  GetConfig() (rpcs.Rpc, error)
  RemoveNode(IP string) (rpcs.Rpc, error)
  EnstablishConnection() error
  ConnectToLeader() error 
  SendConfig(op ConfigChangeOp) (rpcs.Rpc, error) 
}

func NewCluster(IPs []utility.Pair[string,string]) *clusterImpl {
  return &clusterImpl{
    IPs: IPs,
    leaderIP: "",
    conn: nil,
  }
} 
