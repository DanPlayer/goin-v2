package delaytask

import "encoding/json"

const (
	WaitQueue = "DelayTaskWaitQueue"
	ProcessQueue = "DelayTaskProcessQueue"
	ConsumeQueue = "DelayTaskConsumeQueue"
)

type QueueMsgSchema struct {
	NameSpace string
	ID string
	Value string
}

func (r QueueMsgSchema) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}


func (r QueueMsgSchema) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &r)
}