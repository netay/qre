package qre

import (
    "qre/ui"
)


func Main() {
    root, _ := LoadQres()
    var mw ui.MainWindow = ui.NewMainWindow(root)
    mw.Run()
}
