                           ∩~~∩ 
                          ξ ･×･ξ 
                          ξ  ~ ξ 
                          ξ    ξ 
                          ξ   “~～~～〇
                          ξ           ξ	
           d8888 888      ξ ξ ξ~～~ξ  ξ                                
          d88888 888      ξ_ξξ_ξ　ξ_ξξ_ξ                               
         d88P888 888                                     
        d88P 888 888 88888b.   8888b.   .d8888b  8888b.  
       d88P  888 888 888 "88b     "88b d88P"        "88b 
      d88P   888 888 888  888 .d888888 888      .d888888 
     d8888888888 888 888 d88P 888  888 Y88b.    888  888 
    d88P     888 888 88888P"  "Y888888  "Y8888P "Y888888 
                 888                                 
                 888                                 
                 888     


## Overview:
Open-source chess engine written in Golang. UCI and XBoard compatible.
It does not include its own GUI for chess playing, but games can be played from the terminal in console mode.
By leveraging the UCI/XBoard protocols we can use different chess GUI programs like:
- [Arena](http://www.playwitharena.de/)
- [Cute-Chess](https://cutechess.com/)

Alpaca is currently under development.

### Console mode:
Available commands in console mode:
```
Alpaca > help
Commands:
quit - quit game
force - computer will not think
print - show board
post - show thinking
nopost - do not show thinking
new - start a new game
go - set computer thinking
depth x - set depth to x
time x - set thinking time to x seconds (depth still applies if set)
view - show current depth and movetime settings
setboard x - set position to fen x
** note ** - to reset time and depth, set to 0
enter moves using b7b8q notation
```

### Board:
- 120/64 squares board representation
- use of extra padding squares to perform boundary checks (2 ranks top and bottom, 1 file left and right)
- use of bitboards for pawns

### Search:
- [Iterative Deepening](https://www.chessprogramming.org/Iterative_Deepening)
- [Alpha-Beta pruning](https://www.chessprogramming.org/Alpha-Beta)
- [Quiescence Search](https://www.chessprogramming.org/Quiescence_Search)
- [History Heuristic](https://www.chessprogramming.org/History_Heuristic)
- [Killer Move Heuristic](https://www.chessprogramming.org/Killer_Move)
- [MVV-LVA Heuristic](https://www.chessprogramming.org/MVV-LVA)
- [Principal Variation](https://www.chessprogramming.org/Principal_Variation)
- [Null Move Pruning](https://www.chessprogramming.org/Null_Move_Pruning)
- [Transposition Table](https://www.chessprogramming.org/Transposition_Table)


## Perft:

To run all perft tests: 

```
go run cmd/perft/main.go
```

## Build:

```
go build cmd/alpaca/main.go
```

## License:
Alpaca is licensed under the MIT License. See the LICENSE file.

## Resources:

- Minimax: https://ocw.cs.pub.ro/courses/pa/tutoriale/minimax 
- Alpha-beta prunning: https://www.youtube.com/watch?v=xBXHtz4Gbdo&ab_channel=CS188Spring2013
- http://web.archive.org/web/20070707012511/http://www.brucemo.com/compchess/programming/index.htm
- Chess Programming Wiki: https://www.chessprogramming.org/
- https://www.sjeng.org/
- Programming A Chess Engine in C By Bluefever Software: https://www.youtube.com/watch?v=bGAfaepBco4&list=PLZ1QII7yudbc-Ky058TEaOstZHVbT-2hg&ab_channel=BluefeverSoftware
- https://pages.cs.wisc.edu/~psilord/blog/data/chess-pages/index.html
- https://www.youtube.com/@chessprogramming591