package miniIM

const eventName = "simple_im_message" // Actually,it's correspond to the topic name?

type Event struct {
	Msg Message

	Receiver int64 // member id
}
