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

type Selected struct {
	checker  rune
	position Point
}

type model struct {
	checkerboard Board
	cursor       Point
	selected     Selected
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
		selected: Selected{checker: ' ', position: Point{x: -1, y: -1}},
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
		case " ":
			if m.selected.checker == ' ' {
				checker := m.checkerboard[m.cursor.y][m.cursor.x]
				if checker == 'b' || checker == 'w' { // save current checker in Selected
					m.selected = Selected{
						checker: checker, // checker 'w' or 'b'
						position: Point{
							x: m.cursor.x, // point of selected
							y: m.cursor.y, // checker
						},
					}
				}
			} else {
				// if checker already selected and user press space
				// remove checker from old position and unselect
				m.checkerboard[m.cursor.y][m.cursor.x] = m.selected.checker
				m.checkerboard[m.selected.position.y][m.selected.position.x] = '◦'
				m.selected.checker = ' '
			}
		}
	}
	return m, nil
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

	s += fmt.Sprintf("x: %d y: %d, selected: %s", m.cursor.x, m.cursor.y, string(m.selected.checker))
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
