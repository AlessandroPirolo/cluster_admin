package inforesponse

import (
	"raft/admin/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"google.golang.org/protobuf/proto"
	rpc "raft/admin/src/internal/rpcs"
)

type InfoResponse struct {
  pMex protobuf.InfoResponse 
}

func NewInfoResonse() rpc.Rpc  {
  var req = &InfoResponse{
  }
  return req
}

func (this *InfoResponse) Execute() {
  panic("Not implemented")
}

func (this *InfoResponse) ToString() string {
  return " "
} 

func (this *InfoResponse) Encode() ([]byte, error) {
    return proto.Marshal(&(this).pMex)
}

func (this *InfoResponse) Decode(rawMex []byte) error {
  return proto.Unmarshal(rawMex, &this.pMex)
} 

