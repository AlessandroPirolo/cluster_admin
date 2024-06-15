package inforequest

import (
	"raft/admin/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	rpc "raft/admin/src/internal/rpcs"
	"time"

	"google.golang.org/protobuf/proto"
)

type InfoRequest struct {
  pMex protobuf.InfoRequest 
}

func NewInfoRequest() rpc.Rpc  {
  var req = &InfoRequest{
      pMex: protobuf.InfoRequest{
        Timestamp: time.Now().String(),
        ReqType: protobuf.AdminOp_CHANGE_CONF_CHANGE,
      },
    }
  return req
}

func (this *InfoRequest) Execute() {
  panic("Not implemented")
}

func (this *InfoRequest) ToString() string {
  return " "
} 

func (this *InfoRequest) Encode() ([]byte, error) {
    return proto.Marshal(&(this).pMex)
}

func (this *InfoRequest) Decode(rawMex []byte) error {
  return proto.Unmarshal(rawMex, &this.pMex)
} 

