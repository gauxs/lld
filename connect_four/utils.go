package connectfour

import "github.com/gauxs/lld/connectfour/enum"

func nextAvailaibleColor(players []*Player) enum.Color {
	m := make(map[enum.Color]bool)
	for _, p := range players {
		if p != nil {
			m[p.GetPiece().Color()] = true
		}
	}

	allColors := []enum.Color{
		enum.COLOR_BLUE,
		enum.COLOR_RED,
	}

	for _, color := range allColors {
		if !m[color] {
			return color
		}
	}

	return enum.COLOR_INVALID
}
