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
		checkerboard: Board{
			{'◦', 'b', '◦', 'b', '◦', 'b', '◦', 'b'},
			{'b', '◦', 'b', '◦', 'b', '◦', 'b', '◦'},
			{'◦', 'b', '◦', 'b', '◦', 'b', '◦', 'b'},
			{'◦', '◦', '◦', '◦', '◦', '◦', '◦', '◦'},
			{'◦', '◦', '◦', '◦', '◦', '◦', '◦', '◦'},
			{'w', '◦', 'w', '◦', 'w', '◦', 'w', '◦'},
			{'◦', 'w', '◦', 'w', '◦', 'w', '◦', 'w'},
			{'w', '◦', 'w', '◦', 'w', '◦', 'w', '◦'},
		},
		cursor:      Point{x: 0, y: 0},
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
		// Exit Game
		case "ctrl+c", "q":
			return m, tea.Quit
		// Movement
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
		// Selecting
		case " ": // " " it's a `space`
			if m.selected == nil {
				checker := m.checkerboard[m.cursor.y][m.cursor.x]
				if checker == black || checker == white {
					// save current checker in Selected
					m.selected = &Point{x: m.cursor.x, y: m.cursor.y}
				}
			} else {
				dy := abs(m.cursor.y - m.selected.y)
				dx := abs(m.cursor.x - m.selected.x)
				if dy == 0 || dx == 0 {
					return m.stepBack("Try to make a move!"), nil
				}
				// Basic move without capture
				if dx == 1 && dy == 1 {
					if m.checkerboard[m.cursor.y][m.cursor.x] != empty {
						return m.stepBack("You can't place checker here :("), nil
					} else {
						return m.makeMove(), nil
					}
				}
				// Capture
				if dy == 2 && dx == 2 {
					enemyPoint := Point{y: (m.cursor.y + m.selected.y) / 2,
						x: (m.cursor.x + m.selected.x) / 2,
					}
					if enemy := m.checkerboard[enemyPoint.y][enemyPoint.x]; enemy == black || enemy == white {
						m.checkerboard[enemyPoint.y][enemyPoint.x] = empty
						return m.makeMove(), nil
					}
				}
			}
		}
	}
	return m, nil
}

func (m model) stepBack(msg string) model {
	m.infoMessage = msg
	return m
}

func (m model) makeMove() model {

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
}

func abs(n int) int {
	if n < 0 {
		return -n
	} else {
		return n
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
	fmt.Println("↑←↓→ to move cursor,\nspace to select\nq to exit")

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, theres been error: %v", err)
		os.Exit(1)
	}

}
