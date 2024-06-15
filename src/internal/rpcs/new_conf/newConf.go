package newConf

import (
	"log"
	"raft/admin/src/internal/rpcs"
	"raft/admin/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"

	"google.golang.org/protobuf/proto"
)

type NewConf struct {
    pMex protobuf.ChangeConfReq
}

func NewnewConfRPC(op protobuf.AdminOp, conf []string) rpcs.Rpc {
    return &NewConf{
        pMex: protobuf.ChangeConfReq{
            Op: op,
            Conf: &protobuf.ClusterConf{
                Conf: conf,
            },
        },
    }
}

func (this *NewConf) Execute() {
  panic("Not implemented")
}


// ToString implements rpcs.Rpc.
func (this *NewConf) ToString() string {
    return this.pMex.String()
}

func (this *NewConf) Encode() ([]byte, error) {
    var mess []byte
    var err error

    mess, err = proto.Marshal(&(*this).pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote: ", err)
    }

	return mess, err
}
func (this *NewConf) Decode(b []byte) error {
	err := proto.Unmarshal(b,&this.pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote: ", err)
    }
	return err
}

