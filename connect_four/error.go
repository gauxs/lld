package connectfour

import "errors"

var ErrInvalidMove = errors.New("invalid move")
var ErrDiskAlreadyTaken = errors.New("disk with this color is already taken")
var ErrGameNotInCorrectState = errors.New("game is not in correct state")
