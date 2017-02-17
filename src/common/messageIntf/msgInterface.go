package messageIntf

type RegisterAgentMsg struct {
	MsgType	 string
	Hostname string
	Ip       string
	Id       string
}

type AddVdiskMsg struct {
	MsgType string
	AgentId string
	VmId string
	Path string
}

type FeedbackMsg struct {
	MsgType string
	OpResult string	//SUCCESS/FAIL
	Body string
}
