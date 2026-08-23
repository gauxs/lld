package pkg

type Request struct {
	clientID string
	api      string
}

func NewRequest(cID string, a string) *Request {
	return &Request{
		clientID: cID,
		api:      a,
	}
}

func (r *Request) GetClientID() string {
	return ""
}

func (r *Request) GetAPI() string {
	return ""
}
