package ui

import (
    "qre/tree"
    tui "github.com/gizak/termui/v3"
    "github.com/gizak/termui/v3/widgets"
)
func getNodeMenu(n tree.Node, title string) widgets.List {
    newList := widgets.NewList()
    newList.Title = title
    newList.Rows = tree.LsNode(n)
    newList.TextStyle = tui.NewStyle(tui.ColorYellow)
    newList.WrapText = false
    newList.BorderStyle = tui.NewStyle(tui.ColorGreen)
    return *newList
}
type DirItem struct {
    l widgets.List
    n tree.Node
}
type SearchStack struct {
    stack []DirItem
}
func (ss *SearchStack) drop() {
    ss.stack = ss.stack[:len(ss.stack) - 1]
}
func (ss SearchStack) Length() int {
    return len(ss.stack)
}
func (ss *SearchStack) Push(li widgets.List, no tree.Node) {
    di := DirItem {l: li, n: no}
    ss.stack = append(ss.stack, di)
}
func (ss SearchStack) topl() *widgets.List {
    return &ss.stack[len(ss.stack) - 1].l
}
func (ss SearchStack) pretopl() *widgets.List {
    return &ss.stack[len(ss.stack) - 2].l
}
func (ss SearchStack) topn() tree.Node {
    return ss.stack[len(ss.stack) - 1].n
}
func (ss SearchStack) posttop() tree.Node {
    return tree.FindChild(
        ss.topn(),
        ss.topl().Rows,
        ss.topl().SelectedRow,
    )
}
func (ss *SearchStack) updaten() {
    child := tree.FindChild(
        ss.topn(),
        ss.pretopl().Rows,
        ss.pretopl().SelectedRow,
    )
    if child != nil {
        ss.stack[ss.Length() - 1].n = child
    }
}
func (ss *SearchStack) forward(i int) {
    pt := ss.posttop()
    var title string = ss.topl().Rows[ss.topl().SelectedRow]
    lst := getNodeMenu(pt, title)
    ss.topl().BorderStyle = tui.NewStyle(tui.ColorWhite)
    ss.Push(lst, pt)
}

func (ss *SearchStack) navRow() []interface{} {
    nav := make([]interface{}, 0, ss.Length())
    d := ss.stack[ss.Length() - 1]
    col := tui.NewCol(1.0, &d.l)
    nav = append(nav, col)
    return nav
}
