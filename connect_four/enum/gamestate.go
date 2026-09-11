package enum

type GameState int

const (
	GAMESTATE_INVALID = iota
	GAMESTATE_NOT_PLAYING
	GAMESTATE_PLAYING
	GAMESTATE_WON
	GAMESTATE_DRAW
)

func (gs GameState) String() string {
	switch gs {
	case GAMESTATE_NOT_PLAYING:
		return "Not Playing"
	case GAMESTATE_PLAYING:
		return "Playing"
	case GAMESTATE_WON:
		return "Won"
	case GAMESTATE_DRAW:
		return "Draw"
	}

	return "Invalid"
}
