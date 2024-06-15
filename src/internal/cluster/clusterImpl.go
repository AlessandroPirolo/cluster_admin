package cluster

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"raft/admin/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"raft/admin/src/internal/rpcs"
	inforequest "raft/admin/src/internal/rpcs/info_request"
	inforesponse "raft/admin/src/internal/rpcs/info_response"
	conf "raft/admin/src/internal/rpcs/new_conf"
	clientReturnValue "raft/admin/src/internal/rpcs/return_value"
	"raft/admin/src/internal/utility"
	"slices"
	"strings"

	"google.golang.org/protobuf/proto"
)

type ConfigChangeOp int 

const (
  NEW ConfigChangeOp = 0
  CHANGE ConfigChangeOp = 1
)


type clusterImpl struct {
  IPs []utility.Pair[string,string]
  leaderIP string
  conn net.Conn
}

/*
 * Enstablish a connection to a random node in the cluster 
 * (only at beginning when no leader exists)
*/
func (this *clusterImpl) EnstablishConnection() error {
    var conn net.Conn
    var err error
    var addr string = this.IPs[0].Fst

    conn, err = net.Dial("tcp",addr+":8080")
    this.conn = conn
    return err
}

/*
 * Connect to cluster leader
*/
func (this *clusterImpl) ConnectToLeader() error {
    var conn net.Conn
    var err error
    var mex []byte
    var leaderIp protobuf.LeaderIp
    var addr string = this.IPs[0].Fst

    for {
        conn,err = net.Dial("tcp",addr+":8080")
        if err != nil {
            return err
        }
        mex, err = recv(conn)
        err = proto.Unmarshal(mex,&leaderIp)
        if err != nil {
            return err
        }

        log.Println("new leader :", leaderIp.Ip)
        if strings.Contains(leaderIp.Ip,"ok") {
            break
        }
        conn.Close()
        this.leaderIP = leaderIp.Ip
        err = nil
    }

    this.conn = conn
    return nil
}

/*
 * Removing a node by its public IPs
*/
func (this *clusterImpl) RemoveNode(IP string) (rpcs.Rpc, error) {
  var found bool = slices.ContainsFunc(this.IPs, func(element utility.Pair[string, string]) bool {
                                              return IP == element.Fst}) 

  if found == false {
    return nil, errors.New("IP address not found")
  }

  var IPs []utility.Pair[string,string]
  for _,i := range this.IPs {
    if i.Fst != IP {
      IPs = append(IPs, i)
    }
  }
  this.IPs = IPs

  return this.SendConfig(CHANGE)
}

/*
 * Make an info request of type "config" and send it to the cluster
 * in order to get the current configuration
*/
func (this *clusterImpl) GetConfig() (rpcs.Rpc, error) {
  var req, resp rpcs.Rpc
  var reqByte, respByte []byte
  var err, errResp, errDec error

  req = inforequest.NewInfoRequest()
  resp = inforesponse.NewInfoResonse()

  reqByte, err = req.Encode()
  if err != nil {
    return nil, err
  }
  this.conn.Write(reqByte)

  respByte, errResp = recv(this.conn)
  if errResp != nil {
    return nil, errResp
  }
  
  errDec = resp.Decode(respByte)
  if errDec != nil {
    return nil, errDec
  }

  return resp, nil
}

/*
 * Send a new configuration to the cluster, it distinguishes 
 * between a NEW configuration (i.e. starting config) 
 * and a CHANGE configuration (i.e. modified config)
*/
func (this *clusterImpl) SendConfig(op ConfigChangeOp) (rpcs.Rpc, error) {
  var req rpcs.Rpc
  var reqByte, resp []byte
  var err, errResp, errReturn error
  var config []string = this.getPrivateIPs()
  var returnVal clientReturnValue.ClientReturnValue

  switch op {
    case CHANGE: 
      req = conf.NewnewConfRPC(protobuf.AdminOp_CHANGE_CONF_CHANGE, config)  

    default: 
      req = conf.NewnewConfRPC(protobuf.AdminOp_CHANGE_CONF_NEW, config)  
   }

  reqByte, err = req.Encode()
  if err != nil {
    log.Println(err)
  }
  this.conn.Write(reqByte)
  resp, errResp = recv(this.conn)
  if errResp != nil {
    return nil, errResp
  }

  errReturn = returnVal.Decode(resp)
  if errReturn != nil {
    return nil, errReturn
  }

  return &returnVal, nil

}

/*
 * Receiver method
*/
func recv(conn net.Conn) ([]byte, error) {

    buffer := &bytes.Buffer{}

	// Create a temporary buffer to store incoming data
	var tmp []byte = make([]byte, 1024) // Initial buffer size
    var bytesRead int = len(tmp)
    var errConn error
    var errSavi error
    for bytesRead == len(tmp){
		// Read data from the connection
		bytesRead, errConn = conn.Read(tmp)

        // Write the read data into the buffer
        _, errSavi = buffer.Write(tmp[:bytesRead])
        
        // check error saving
        if errSavi != nil {
            return nil, errSavi
        }

		if errConn != nil {
			if errConn != io.EOF {
				// Handle other errConnors
				return nil, errConn
			}

            if errConn == io.EOF {
                return nil, errConn
            }
			break
		}
	}
	return buffer.Bytes(), nil 

}

func (this *clusterImpl) getPrivateIPs() []string {
  var IPs []string
  for _,i := range this.IPs {
    IPs = append(IPs, i.Snd)
  }
  return IPs
}

func (this *clusterImpl) getPublicIPs() []string {
  var IPs []string
  for _,i := range this.IPs {
    IPs = append(IPs, i.Fst)
  }
  return IPs
}
