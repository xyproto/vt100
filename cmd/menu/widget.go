package main

import (
	"github.com/xyproto/vt100"
)

type MenuWidget struct {
	title      string // title
	w          uint   // width
	h          uint   // height (number of menu items)
	y          uint   // current position
	oldy       uint   // previous position
	fg         string // foreground color
	active     string // active (selected) color
	hi         string // highlight color
	marginLeft int    // margin, may be negative?
	marginTop  int    // margin, may be negative?
	choices    []string
	selected   int
	titleColor string
}

func NewMenuWidget(title, titleColor string, choices []string, fg, hi, active string, canvasWidth, canvasHeight uint) *MenuWidget {
	maxlen := uint(0)
	for _, choice := range choices {
		if uint(len(choice)) > uint(maxlen) {
			maxlen = uint(len(choice))
		}
	}
	marginLeft := 10
	if int(canvasWidth)-(int(maxlen)+marginLeft) <= 0 {
		marginLeft = 0
	}
	marginTop := 10
	if int(canvasHeight)-(len(choices)+marginTop) <= 0 {
		marginTop = 0
	}
	return &MenuWidget{
		title:      title,
		w:          uint(marginLeft + int(maxlen)),
		h:          uint(len(choices)),
		y:          0,
		oldy:       0,
		fg:         fg,
		hi:         hi,
		active:     active,
		marginLeft: marginLeft,
		marginTop:  marginTop,
		choices:    choices,
		selected:   -1,
		titleColor: titleColor,
	}
}

func (m *MenuWidget) Selected() int {
	return m.selected
}

func (m *MenuWidget) Draw(c *vt100.Canvas) {
	// Draw the title
	titleHeight := 2
	for x, r := range m.title {
		c.PlotColor(uint(m.marginLeft+x), uint(m.marginTop), vt100.LightColorMap[m.titleColor], r)
	}
	// Draw the menu entries, with various colors
	ulenChoices := uint(len(m.choices))
	for y := uint(0); y < m.h; y++ {
		itemString := "---"
		if y < ulenChoices {
			itemString = "-> " + m.choices[y] + " ---"
		}
		for x := uint(0); x < m.w; x++ {
			r := '-'
			if x < uint(len([]rune(itemString))) {
				r = []rune(itemString)[x]
			}
			if x < 2 {
				c.PlotColor(uint(m.marginLeft+int(x)), uint(m.marginTop+int(y)+titleHeight), vt100.LightBlue, r)
			} else if y == m.y {
				c.PlotColor(uint(m.marginLeft+int(x)), uint(m.marginTop+int(y)+titleHeight), vt100.LightColorMap[m.hi], r)
			} else {
				c.PlotColor(uint(m.marginLeft+int(x)), uint(m.marginTop+int(y)+titleHeight), vt100.LightColorMap[m.fg], r)
			}
		}
	}
}

func (m *MenuWidget) SelectDraw(c *vt100.Canvas) {
	old := m.hi
	m.hi = m.active
	m.Draw(c)
	m.hi = old
}

func (m *MenuWidget) Select() {
	m.selected = int(m.y)
}

func (m *MenuWidget) Up(c *vt100.Canvas) bool {
	m.oldy = m.y
	if m.y <= 0 {
		m.y = m.h - 1
	} else {
		m.y--
	}
	return true
}

func (m *MenuWidget) Down(c *vt100.Canvas) bool {
	m.oldy = m.y
	m.y++
	if m.y >= m.h {
		m.y = 0
	}
	return true
}

// Terminal was resized
func (m *MenuWidget) Resize() {
	//m.hi = "Magenta"
}

// Select a specific index, if possible. Returns false if it was not possible.
func (m *MenuWidget) SelectIndex(n uint) bool {
	if n >= m.h {
		return false
	}
	m.oldy = m.y
	m.y = n
	return true
}

func (m *MenuWidget) SelectFirst() bool {
	return m.SelectIndex(0)
}

func (m *MenuWidget) SelectLast() bool {
	m.oldy = m.y
	m.y = m.h - 1
	return true
}
