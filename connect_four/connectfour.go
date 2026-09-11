package connectfour

import (
	"fmt"

	"github.com/gauxs/lld/connect_four/enum"
)

func Execute() {
	fmt.Println("Starting the game ConnectFour")

	game := NewGame(6, 7)
	if err := game.AddPlayer("Player-A", enum.DISK_RED); err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := game.AddPlayer("Player-B", enum.DISK_BLUE); err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := game.StartGame(); err != nil {
		fmt.Println(err.Error())
		return
	}

	for game.GetGameState() == enum.GAMESTATE_PLAYING {
		col := 0
		curPlayer := game.GetCurrentPlayer()
		fmt.Printf("Enter column for player %v: ", curPlayer.GetName())
		fmt.Scanln(&col)

		if err := game.MakeMove(col); err != nil {
			fmt.Println(err.Error())
		}
	}

	fmt.Printf("Game ended in %v \n", game.GetGameState().String())
	if game.GetGameState() == enum.GAMESTATE_WON {
		fmt.Printf("The winner of te game is: %v \n", game.GetWinner().GetName())
	}
}
