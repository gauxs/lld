package connectfour

import (
	"github.com/gauxs/lld/connect_four/enum"
)

type Player struct {
	name string
	disk enum.Disc
}

func (p *Player) GetName() string {
	return p.name
}
func (p *Player) GetDisk() enum.Disc {
	return p.disk
}
