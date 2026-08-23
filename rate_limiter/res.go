package ratelimiter

type Resource struct {
	id string
}

func (r *Resource) GetID() string {
	return ""
}
