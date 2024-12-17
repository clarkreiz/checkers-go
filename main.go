// Terminology:
// checkers - name of the game
// checkerboard - game desk for playing checkers
// checker - name of one `figure`, may be white or black (w,b in this game)

package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Board = [][]rune
type pieceType = rune

const (
	white pieceType = 'w'
	black pieceType = 'b'
	empty pieceType = '◦'
)

type Point struct {
	x int
	y int
}

type model struct {
	checkerboard Board     // game desk for checkers
	cursor       Point     // cursor use for moving through checkerboard and select checker
	selected     *Point    // for selected checker (default is nil)
	infoMessage  string    // special field for info text message
	turnMessage  string    // who turn?
	currentTurn  pieceType // who turn?
}

func (m model) Init() tea.Cmd {
	return nil
}

func initialModel() model {
	game := model{
		checkerboard: Board{ // i can use white/black const,
			// but more i like this code
			{'◦', 'b', '◦', 'b', '◦', 'b', '◦', 'b'},
			{'b', '◦', 'b', '◦', 'b', '◦', 'b', '◦'},
			{'◦', 'b', '◦', 'b', '◦', 'b', '◦', 'b'},
			{'◦', '◦', '◦', '◦', '◦', '◦', '◦', '◦'},
			{'◦', '◦', '◦', '◦', '◦', '◦', '◦', '◦'},
			{'w', '◦', 'w', '◦', 'w', '◦', 'w', '◦'},
			{'◦', 'w', '◦', 'w', '◦', 'w', '◦', 'w'},
			{'w', '◦', 'w', '◦', 'w', '◦', 'w', '◦'},
		},
		cursor:      Point{x: 0, y: 7},
		selected:    nil,
		infoMessage: "Go-go!",
		currentTurn: white,
		turnMessage: "White turn.",
	}
	return game
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Exit the Game
		case "ctrl+c", "q":
			return m, tea.Quit
		// Cursor movement
		case "h", "j", "k", "l":
			return moveCursor(msg.String(), m), nil
		// Selection and all core game logic handled here ↓ :)
		case " ":
			if m.selected == nil { // first time we simply cath the checker
				checker := m.checkerboard[m.cursor.y][m.cursor.x]
				if checker == black || checker == white {
					m.selected = &Point{x: m.cursor.x, y: m.cursor.y}
				}
			} else { // secondly we make a move
				return makeMove(*m.selected, m.cursor, m), nil
			}
		}
	}
	return m, nil
}

func abs(n int) int {
	if n < 0 {
		return -n
	} else {
		return n
	}
}

func moveCursor(msg string, m model) model {
	switch msg {
	case "h", "left":
		if m.cursor.x > 0 {
			m.cursor.x--
		}
	case "j", "down":
		if m.cursor.y < 7 {
			m.cursor.y++
		}
	case "k", "up":
		if m.cursor.y > 0 {
			m.cursor.y--
		}
	case "l", "right":
		if m.cursor.x < 7 {
			m.cursor.x++
		}
	}
	return m
}

func makeMove(from, to Point, m model) model {
	var validMat [8][8]int

	for i, line := range validMat {
		for j := range line {
			if i%2 == 0 {
				validMat[i][j] = (j % 2) & pieceToValid(m.checkerboard[i][j])
			} else {
				validMat[i][j] = (1 ^ j%2) & pieceToValid(m.checkerboard[i][j])
			}
		}
	}
	dy := abs(to.y - from.y)
	dx := abs(to.x - from.x)
	if dy <= 1 && dx <= 1 && validMat[to.y][to.x] != 0 {
		m.checkerboard[m.cursor.y][m.cursor.x] = m.checkerboard[m.selected.y][m.selected.x]
		m.checkerboard[m.selected.y][m.selected.x] = empty
		m.selected = nil
		m.infoMessage = "Okay, good!"
		if m.currentTurn == white {
			m.currentTurn = black
			m.turnMessage = "Black turn!"
		} else {
			m.currentTurn = white
			m.turnMessage = "White turn!"
		}
		return m
	} else {
		// TODO: add the ability to capture a checker
		// if dy > 1 || dx > 1 && validMat[to.y][to.x] != 0 {
		// 	makeCapture
		// } else
		m.infoMessage = "Something went wrong"
		return m
	}
}

func pieceToValid(piece pieceType) int {
	// A checker may be placed only on an empty field
	if piece == empty {
		return 1
	} else {
		return 0
	}
}

func (m model) View() string {
	var s string
	for i, line := range m.checkerboard {
		for j, cell := range line {
			if i == m.cursor.y && j == m.cursor.x {
				switch cell {
				case empty:
					s += " ○ "
				case white:
					s += " W "
				case black:
					s += " B "
				}
			} else {
				s += " "
				s += string(cell)
				s += " "
			}
		}
		s += "\n"
	}

	s += fmt.Sprintf("x: %d y: %d", m.cursor.x, m.cursor.y)
	s += fmt.Sprint("\n", m.infoMessage)
	s += fmt.Sprint(" ", m.turnMessage)
	return s
}

func main() {
	fmt.Println("Hi! This is Checkers")
	fmt.Println("hjkl to move cursor,\nspace to select\nq to exit")

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, theres been error: %v", err)
		os.Exit(1)
	}

}
