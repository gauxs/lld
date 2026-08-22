## Requirements
7 × 6 board.
Two players initially.
Players take turns dropping pieces into columns.
Pieces occupy the lowest available position.
Four consecutive pieces wins.
Horizontal, vertical, and diagonal wins.
Full board without winner = draw.
No moves after game completion.
Future support for more players.
Future configurable consecutive-piece count and winning directions.

## Entities
Game
Player
Board
Piece

## Relationships
Game contains Players and Board.
Board contains Pieces.
Piece belongs to a Player.
Game coordinates turns and game state.
Board manages board state/placement.