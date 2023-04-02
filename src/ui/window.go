package ui

import (
    "log"
    "qre/tree"
    tui "github.com/gizak/termui/v3"
    "github.com/gizak/termui/v3/widgets"
)

type MainWindow struct {
    root tree.Node
}

func NewMainWindow(
    n tree.Node,
) MainWindow {
    return MainWindow {
        root: n,
    }
}

func (mw MainWindow) Run() {
    menu := tree.LsNode(mw.root)

    if err := tui.Init(); err != nil {
        log.Fatalf("Failde to initialize termui: %v", err)
    }
    defer tui.Close()

    l := widgets.NewList()
    l.Title = "Files"
    l.Rows = menu
    l.TextStyle = tui.NewStyle(tui.ColorYellow)
    l.WrapText = false
    l.BorderStyle = tui.NewStyle(tui.ColorGreen)

    var ss SearchStack
    ss.Push(*l, mw.root)

    if ss.Length() == 0 {
        panic("Stack must have root and choice from it.")
    }
    view := widgets.NewParagraph()
    view.Text = "Move by arrows (left, right)\nExit by q"
    view.Title = "References"

    termWidth, termHeight := tui.TerminalDimensions()
    grid := tui.NewGrid()
    grid.SetRect(0, 0, termWidth, termHeight)

    grid.Set(
        tui.NewRow(0.5, ss.navRow()...),
        tui.NewRow(0.5, view),
    )

    tui.Render(grid)
    previousKey := ""
    uiEvents := tui.PollEvents()
    for {
        e := <-uiEvents
        switch e.ID {
        case "q", "C-c", "<Escape>":
            return
        case "<Up>", "<MouseWheelUp>":
            ss.topl().ScrollUp()
        case "<Down>", "<MouseWheelDown>":
            ss.topl().ScrollDown()
        case "<PageUp>":
            for i := 0; i < 5; i++ {
               ss.topl().ScrollUp()
            }
        case "<PageDown>":
            for i := 0; i < 5; i++ {
                ss.topl().ScrollDown()
            }
        case "<Right>", "<MouseRight>":
            if tree.CanMoveRight(ss.topn()) {
                ss.forward(ss.topl().SelectedRow)
            } else {
                continue
            }
        case "<Left>", "<MouseLeft>":
            if ss.Length() > 1 {
                ss.drop()
                ss.topl().BorderStyle = tui.NewStyle(tui.ColorGreen)
            } else {
                continue
            }
        case "<Home>":
            ss.topl().ScrollTop()
        case "g":
            if previousKey == "g" {
                ss.topl().ScrollTop()
            }
        case "G", "<End>":
            ss.topl().ScrollBottom()
        case "<Resize>":
            termWidth, termHeight = tui.TerminalDimensions()
        default:
            continue
        }

        grid = tui.NewGrid()
        grid.SetRect(0, 0, termWidth, termHeight)
        grid.Set(
            tui.NewRow(0.5, ss.navRow()...),
            tui.NewRow(0.5, view),
        )
        pt := ss.posttop()
        if pt != nil {
            view.Text = tree.PrintNode(pt)
        }

        if previousKey == "g" {
            previousKey = ""
        } else {
            previousKey = e.ID
        }

        tui.Render(grid)
    }
}
