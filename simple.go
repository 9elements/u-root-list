package list

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"
)

var maxLength int = 120

type SimpleItem struct {
	Title, Desc string
	Disabled    bool

	Options        []string
	SelectedOption string
}

type SimpleItemList []SimpleItem

var _ fuzzy.Source = (SimpleItemList)(nil)

func (s SimpleItemList) Len() int {
	return len(s)
}

func (s SimpleItemList) String(i int) string {
	return s[i].Title
}

type SimpleAdapter struct {
	items             SimpleItemList
	filterResult      []fuzzy.Match
	lastFilterPattern string

	StyleNormal, StyleDimmed *SimpleAdapterStyle
}

func NewSimpleAdapter(items SimpleItemList) *SimpleAdapter {
	styleNormal, styleDimmed := SimpleDefaultStyle()
	return &SimpleAdapter{
		items: items,

		StyleNormal: styleNormal,
		StyleDimmed: styleDimmed,
	}
}

var _ Adapter = (*SimpleAdapter)(nil)

func (s *SimpleAdapter) Len() int {
	if s.filterResult == nil {
		return len(s.items)
	}
	return len(s.filterResult)
}

func (s *SimpleAdapter) Sep() string {
	return "\n\n"
}

func (s *SimpleAdapter) View(pos, focus int, expanded bool) string {
	item := s.items[pos]

	if focus == pos {
		return buildTitleAndHelptext(item, true, expanded)
	} else {
		return buildTitleAndHelptext(item, false, false)
	}
}

func buildTitleAndHelptext(item SimpleItem, focus bool, expanded bool) string {
	baseStyle := lipgloss.NewStyle()

	if focus {
		baseStyle = baseStyle.BorderStyle(lipgloss.NormalBorder()).BorderLeft(true)
	} else {
		baseStyle = baseStyle.PaddingLeft(2)
	}

	baseStyle = baseStyle.Align(lipgloss.Left)

	maxItemTitleLength := maxLength - len(item.SelectedOption) - 5

	padding := strings.Repeat(" ", maxLength-min(len(item.Title), maxItemTitleLength)+5-len(item.SelectedOption))

	baseItem := baseStyle.Render(item.Title[:min(len(item.Title), maxItemTitleLength)] + padding + item.SelectedOption + "\n" +
		lipgloss.NewStyle().Foreground(lipgloss.Color("#707070")).Render(item.Desc))

	if !expanded {
		return baseItem
	}

	baseItem += "\n"
	// Do we need a nested model here?
	return baseItem
}

func (s *SimpleAdapter) Append(item ...SimpleItem) {
	s.items = append(s.items, item...)
	if s.filterResult != nil {
		s.Filter(s.lastFilterPattern)
	}
}

func (s *SimpleAdapter) Insert(i int, item ...SimpleItem) {
	s.items = append(s.items[:i+len(item)], s.items[i:]...)
	for i2 := 0; i2 < len(item); i2++ {
		s.items[i+i2] = item[i2]
	}
	if s.filterResult != nil {
		s.Filter(s.lastFilterPattern)
	}
}

func (s *SimpleAdapter) Remove(i int) {
	s.items = append(s.items[:i], s.items[i+1:]...)
	if s.filterResult != nil {
		s.Filter(s.lastFilterPattern)
	}
}

func (s *SimpleAdapter) Filter(pattern string) {
	s.lastFilterPattern = pattern
	if len(pattern) > 0 {
		s.filterResult = fuzzy.FindFrom(pattern, s.items)
		if s.filterResult == nil {
			s.filterResult = make([]fuzzy.Match, 0)
		}
	} else {
		s.filterResult = nil
	}
}

func (s *SimpleAdapter) OriginalItemLen() int {
	return len(s.items)
}

func (s *SimpleAdapter) FilteredIndex(pos int) int {
	if s.filterResult == nil {
		return pos
	}
	return s.filterResult[pos].Index
}

func (s *SimpleAdapter) FilteredItemAt(pos int) SimpleItem {
	if s.filterResult == nil {
		return s.items[pos]
	}
	return s.items[s.filterResult[pos].Index]
}

func (s *SimpleAdapter) ItemAt(pos int) SimpleItem {
	return s.items[pos]
}

func (s *SimpleAdapter) SetItemAt(pos int, item SimpleItem) {
	s.items[pos] = item
	if s.filterResult != nil {
		s.Filter(s.lastFilterPattern)
	}
}

func (s *SimpleAdapter) SetItems(items SimpleItemList) {
	s.items = items
	if s.filterResult != nil {
		s.Filter(s.lastFilterPattern)
	}
}
