package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"raft/admin/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"raft/admin/src/internal/admin_cli"
	clusterform "raft/admin/src/internal/admin_cli/model/cluster_form"
	"raft/admin/src/internal/cluster"
  "raft/admin/src/internal/vm_manager"
	
  "strings"

	"google.golang.org/protobuf/proto"
	"libvirt.org/go/libvirt"
)

//virsh --connect=qemu:///system list
//virsh --connect=qemu:///system domifaddr --source arp --domain VM_RAFT_3
//ip range of cluster 192.168.122.2 192.168.122.255
const (
  green string = "\x1b[32m"
  white string = "\x1b[0m"
  red string = "\x1b[31m"
)

func main()  {

    cl := vmmanager.NewClusterManager()
    log.Printf("Creating virtual machine... ") 
    clust := cl.Init()
    log.Printf("%sSUCCESS\n%s", green, white)
    log.Printf("Connecting to the cluster")
    err := clust.EnstablishConnection()
    if err != nil {
      log.Printf("%sFAILED\n%s", red, white)
      log.Println(err)
    } else {
      log.Printf("%sSUCCESS\n%s", green, white)
    }
    res, err:= clust.SendConfig(cluster.NEW)
    if err != nil {
      log.Printf("%sFAILED\n%s", red, white)
      log.Println(err)
    } else {
      
      log.Printf("%sSUCCESS\n%s", green, white)
    }
    
    /*var conn net.Conn
    var mex []byte
    var ipAddr string*/

    /*
  TODO:   
        implement different operation type with an interaction by the user,
        also implement add data in the payload field, when needed.
        OPERATION TO ADD:
            1- CREATE create ad file with name specified in the payload
            2- READ (add the name of the file to read in the payload)
  3- WRITE (write data in the file, if it does not exist the file will be created)
            4- RENAME (rename a file)
            5- DELETE (delete a file from the cluster)
    */

    var cli = admincli.NewCli()
    parameters, err := cli.Start()
    if err != nil {
      log.Panicln("Error in CLI: ", err)
    }

    clust.ConnectToLeader()
    
    switch parameters["operation"] {
      case string(clusterform.ADD): 
        /* ADD NODE*/
        clust = cl.AddNode()
        clust.SendConfig(cluster.CHANGE)
      case string(clusterform.REMOVE): 
        /* REMOVE NODE*/
        res, err := clust.RemoveNode(parameters["ip"])
        if err != nil {
          log.Printf("%s: %s", err, parameters["ip"])
        }
        log.Println(res)
        cl.RemoveNode(parameters["ip"])
      case string(clusterform.GET): 
        /* GET CONFIG*/
        clust.GetConfig()

    }
}

func EncodeMessage(req *protobuf.ClientReq) []byte{
    var mex []byte
    var err error

    mex, err = proto.Marshal(req)
    if err != nil {
        panic("error encoding")
    }
    return mex
}

func SendCluster(conn net.Conn, mex []byte){
    var err error

    _,err = conn.Write(mex)
    if err != nil {
        panic("error sending data")
    }
}

func ConnectToLeader(ipAddr string) net.Conn{
    var conn net.Conn
    var err error
    var mex []byte
    var leaderIp protobuf.LeaderIp

    for {
        conn,err = net.Dial("tcp",ipAddr+":8080")
        if err != nil {
            log.Panicf("error dialing conection %v\n",err)
        }
        mex, err = Recv(conn)
        err = proto.Unmarshal(mex,&leaderIp)
        if err != nil {
            panic("failed decoding ip of leader")
        }

        log.Println("new leader :", leaderIp.Ip)
        if strings.Contains(leaderIp.Ip,"ok") {
            break
        }
        conn.Close()
        ipAddr = leaderIp.Ip
        err = nil
    }

    return conn
}

func GetCluterNodeIp() string{
    var connHyper, err =libvirt.NewConnect("qemu:///system")
    if err != nil {
        panic("failed connection to hypervisor")
    }

    var doms []libvirt.Domain
    doms,err = connHyper.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
    if err != nil {
        panic("failed retrieve domains")
    }

    var ip []libvirt.DomainInterface
    ip, err = doms[0].ListAllInterfaceAddresses(libvirt.DOMAIN_INTERFACE_ADDRESSES_SRC_ARP)

    var ipAddr string = ip[0].Addrs[0].Addr

    if !strings.Contains(ipAddr,"200.168.122"){
        ipAddr = ip[1].Addrs[0].Addr
    }

    return ipAddr
}

func  Recv(conn net.Conn) ([]byte, error) {

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
