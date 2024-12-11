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

type model struct {
	checkerboard Board
	cursor       Point
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
		cursor: Point{x: 0, y: 0},
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
	}
	return m, nil
}

func (m model) View() string {
	var s string
	for _, lines := range m.checkerboard {
		switch m.checkerboard[m.cursor.y][m.cursor.x] {
		case '◦':
			m.checkerboard[m.cursor.y][m.cursor.x] = '○'
		}

		s += fmt.Sprintf("%c\n", lines)
	}

	s += fmt.Sprintf("x: %d y: %d", m.cursor.x, m.cursor.y)
	return s
}

func main() {
	// fmt.Println("Hi! This is Checkers")
	// fmt.Println("hjkl to move cursor, space to select\nq to exit")

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, theres been error: %v", err)
		os.Exit(1)
	}

}
