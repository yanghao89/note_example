package zkstore

import (
	"encoding/json"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

type ZkClient struct {
	zkServer []string `json:"zk_server"`
	zkConn   *zk.Conn
	zkRoot   string
}

func NewZkClient(zkServer []string, root string, timeout int) (*ZkClient, error) {
	client := &ZkClient{
		zkServer: zkServer,
	}
	conn, _, err := zk.Connect(zkServer, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	client.zkConn = conn
	client.zkRoot = root
	if err := client.EnsurePath(root); err != nil {
		client.Close()
		return nil, err
	}
	return client, nil
}

func (s *ZkClient) Close() {
	s.zkConn.Close()
}

func (z *ZkClient) EnsurePath(path string) error {
	exists, _, err := z.zkConn.Exists(path)
	if err != nil {
		return err
	}
	if !exists {
		if _, err := z.zkConn.Create(path, []byte(""), 0, zk.WorldACL(zk.PermAll)); err != nil && err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}

type ServerNode struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (z *ZkClient) Register(node *ServerNode) error {
	path := z.zkRoot + "/" + node.Name
	if err := z.EnsurePath(path); err != nil {
		return err
	}
	path = path + "/n"
	data, err := json.Marshal(node)
	if err != nil {
		return err
	}
	if _, err := z.zkConn.CreateProtectedEphemeralSequential(path, data, zk.WorldACL(zk.PermAll)); err != nil {
		return err
	}
	return nil
}

func (z *ZkClient) GetNodes(name string) ([]*ServerNode, error) {
	path := z.zkRoot + "/" + name
	child, _, err := z.zkConn.Children(path)
	if err != nil {
		if err == zk.ErrNoNode {
			return []*ServerNode{}, nil
		}
		return nil, err
	}
	nodes := make([]*ServerNode, 0, len(child))
	for _, v := range child {
		fullPath := path + "/" + v
		data, _, err := z.zkConn.Get(fullPath)
		if err != nil {
			if err == zk.ErrNoNode {
				continue
			}
			return nil, err
		}
		node := new(ServerNode)
		err = json.Unmarshal(data, node)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}
