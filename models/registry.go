package models

type RegisterNodeReq struct {
	Port int
}

type NewNodeHook struct {
	Address string
}

type ListPeersResp struct {
	Peers []string
}
