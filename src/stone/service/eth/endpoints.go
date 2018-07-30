package eth

import (
	"net/url"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"stone/common"
)

type endpoint struct {
	weight  int
	url     string
	isOk    bool
	details *EthereumNodeInfo
}

func (e *endpoint) rpc(result interface{}, method string, args ...interface{}) error {
	client, err := rpc.Dial(e.url)
	if err != nil {
		common.Logger.Error("dial error in rpc: ", e.url)
		return err
	}
	err = client.Call(result, method, args...)
	if err != nil {
		common.Logger.Error("client.Call error", err)
		return err
	}
	return nil
}

func (e *endpoint) heartbeat() (bool, *EthereumNodeInfo) {
	var res string
	details := new(EthereumNodeInfo)
	err := e.rpc(&res, "net_version")
	if err != nil {
		common.Logger.Info("heartbeat error: ", e.url)
		return false, nil
	}
	err = e.rpc(&res, "eth_coinbase")
	if err != nil {
		common.Logger.Info("heartbeat error: ", e.url)
		return false, nil
	}
	details.Miner = res
	var mining bool
	err = e.rpc(&mining, "eth_mining")
	if err != nil {
		common.Logger.Info("heartbeat error: ", e.url)
		return false, nil
	}
	var lastest string
	err = e.rpc(&lastest, "eth_blockNumber")
	if err != nil {
		common.Logger.Info("heartbeat error: ", e.url)
		return false, nil
	}
	details.Is_mining = mining
	/*
	snapshot := make(map[string]interface{})
	err = e.rpc(&snapshot, "clique_getSnapshot", lastest)
	if err != nil {
		common.Logger.Info("heartbeat error: ", e.url)
		return false, nil
	}
	captial := snapshot["signers"].(map[string]interface{})
	if captial[details.Miner] != nil {
		details.Is_mining = mining
	} else {
		details.Is_mining = false
	}

	nodeinfo := make(map[string]interface{})
	err = e.rpc(&nodeinfo, "admin_nodeInfo")
	if err != nil {
		common.Logger.Info("heartbeat error: ", e.url)
		return false, nil
	}
	details.Version = nodeinfo["name"].(string)
	*/
	details.URL = e.url
	details.Is_alive = e.isOk
	return true, details
}

// EndpointsManager endpoints of ethereum
type EndpointsManager struct {
	endpoints       []*endpoint
	rAliveEndpoints []*endpoint
	rwMutex         sync.RWMutex
	exit            chan bool
	closed          chan bool
}

// NewEndPointsManager create a endPoint manager
func NewEndPointsManager() *EndpointsManager {
	return &EndpointsManager{
		endpoints:       []*endpoint{},
		rAliveEndpoints: []*endpoint{},
		exit:            make(chan bool),
		closed:          make(chan bool),
	}
}

func (e *EndpointsManager) AddEndPoint(endpointURL string, weight int) {
	e.rwMutex.Lock()
	defer e.rwMutex.Unlock()
	endpoint := &endpoint{
		url:    endpointURL,
		weight: weight}
	e.endpoints = append(e.endpoints, endpoint)
	e.rAliveEndpoints = append(e.rAliveEndpoints, endpoint)
}
func (e *EndpointsManager) GetEndPoints() []*EthereumNodeInfo {
	e.rwMutex.Lock()
	defer e.rwMutex.Unlock()
	nodes := []*EthereumNodeInfo{}
	for _, item := range e.endpoints {
		nodes = append(nodes, item.details)
	}
	return nodes
}

// Run endpoints run, monitor alive Endpoint
func (e *EndpointsManager) Run() {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			e.watchAliveEndpoint()
		case <-e.exit:
			close(e.closed)
			common.Logger.Info("service done!!!")
			return
		}
	}
}

// Stop Stop
func (e *EndpointsManager) Stop() {
	close(e.exit)
	// wait for stop
	<-e.closed
}

func (e *EndpointsManager) watchAliveEndpoint() error {
	for _, item := range e.endpoints {
		item.isOk, item.details = item.heartbeat()
	}
	e.updateAliveEndpoint()
	if len(e.rAliveEndpoints) != 0 {
		common.Logger.Info("endpoint watch: ", len(e.rAliveEndpoints), ": ", e.rAliveEndpoints[0].url)
	} else {
		common.Logger.Info("endpoint watch: ", len(e.rAliveEndpoints))
	}

	return nil
}

func (e *EndpointsManager) updateAliveEndpoint() {
	e.rwMutex.Lock()
	defer e.rwMutex.Unlock()
	res := []*endpoint{}
	for _, item := range e.endpoints {
		if item.isOk {
			res = append(res, item)
		}
	}
	e.rAliveEndpoints = res
}

// RPC rpc
func (e *EndpointsManager) RPC(result interface{}, method string, args ...interface{}) (err error) {
	e.rwMutex.RLock()
	defer e.rwMutex.RUnlock()
	for _, item := range e.rAliveEndpoints {
		err = item.rpc(result, method, args...)
		if _, ok := err.(*url.Error); ok {
			continue
		} else {
			break
		}
	}
	return
}
