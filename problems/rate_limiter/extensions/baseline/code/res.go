package ratelimiter

type Resource struct {
	id string
}

func NewResource(id string) *Resource {
	return &Resource{
		id: id,
	}
}

func (r *Resource) GetID() string {
	return r.id
}
