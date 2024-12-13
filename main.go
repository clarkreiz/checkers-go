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

type Point struct {
	x int
	y int
}

// type Selected struct {
// 	checker  rune // default is ' ', may be 'b' or 'w'
// 	position Point
// }

type model struct {
	checkerboard Board  // game desk for checkers
	cursor       Point  // cursor use for moving through checkerboard and select checker
	selected     *Point // struct for selected checker with position and rune (default is ' ')
	info         string // special field for info text message
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
		cursor:   Point{x: 0, y: 0},
		selected: nil,
		info:     "You awesome!",
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
				if checker == 'b' || checker == 'w' {
					// save current checker in Selected
					m.selected = &Point{x: m.cursor.x, y: m.cursor.y}
					m.info = "Put a checker where you want!"
				}
			} else {
				dy := abs(m.cursor.y - m.selected.y)
				dx := abs(m.cursor.x - m.selected.x)
				if dx > 1 || dy > 1 || dy == 0 || dx == 0 {
					m.info = "You cannot put checker here :("
					return m, nil
				}
				m.checkerboard[m.cursor.y][m.cursor.x] = m.checkerboard[m.selected.y][m.selected.x]
				m.checkerboard[m.selected.y][m.selected.x] = '◦'
				m.selected = nil
				m.info = "Okay, good!"
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

func (m model) View() string {
	var s string
	for i, line := range m.checkerboard {
		for j, cell := range line {
			if i == m.cursor.y && j == m.cursor.x {
				switch cell {
				case '◦':
					s += " ○ "
				case 'w':
					s += " W "
				case 'b':
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
	s += fmt.Sprint("\n", m.info)
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
