package tree

import (
    "io"
    "sort"
    "strings"
)

type Node interface {
    app(
        io.StringWriter, // where to write
        int,              // margin, depth of the node
        func(io.StringWriter, int, ...interface{})error,
    )
}
type Leaf struct {
    value string
}
type SNode struct { // objects in lists
    name string
    ch   Node
}
type LNode struct { // list-node
    ch   []Node
}
type MNode struct { // map-node
    ch map[string]Node
}

func NewMNode(ch map[string]Node) MNode {
    var node = MNode {ch: ch}
    return node
}

func (l Leaf) app(
    w io.StringWriter,
    m int,
    f func(io.StringWriter, int, ...interface{})error,
) {
    f(w, m, l.value + "\n")
}
func (n MNode) app(
    w io.StringWriter,
    m int,
    f func(io.StringWriter, int, ...interface{})error,
) {
    for key, value := range n.ch {
        switch v := value.(type) {
        case Leaf:
            f(w, m, key + "[ // ", v.value + "](fg:yellow,mod:bold)\n")
        default:
            f(w, m, key + ":\n")
            value.app(w, m + 1, f)
        }
    }
}
func (l LNode) app(
    w io.StringWriter,
    m int,
    f func(io.StringWriter, int, ...interface{})error,
) {
    for _, n := range l.ch {
        n.app(w, m + 1, f)
    }
}
func (s SNode) app(
    w io.StringWriter,
    m int,
    f func(io.StringWriter, int, ...interface{})error,
) {
    f(w, m, s.name)
    s.ch.app(w, m, f)
}

func CanMoveRight(n Node) bool {
    var ch = FindChildren(n)
    if ch == nil {
        return false
    }
    for _, n := range ch {
        switch n.(type) {
        case Leaf, SNode: {
            return false
        }
        default:
            continue
        }
    }
    return true
}

func LsNode(n Node) []string {
    switch v := n.(type) {
    case LNode:
        ch := make([]string, 0)
        for _, v := range v.ch {
            ch = append(ch, LsNode(v)...)
        }
        return ch
    case MNode:
        ch := make([]string, 0)
        for k, _ := range v.ch {
            ch = append(ch, k)
        }
        sort.Strings(ch)
        return ch
    case Leaf:
        return []string{v.value}
    default:
        return []string{}
    }
}

func FindChildren(n Node) []Node {
    var li []Node
    switch n := n.(type) {
    case LNode:
        for _, v := range n.ch {
            switch y := v.(type) {
            case MNode:
                for _, u := range y.ch {
                    li = append(li, u)
                }
            default:
                li = append(li, v)
            }
        }
    case MNode:
        for _, v := range n.ch {
            li = append(li, v)
        }
    default:
        return []Node{}
    }
    return li
}

func FindChild(n Node, names []string, nodeIndex int) Node {
    var child Node
    switch n := n.(type) {
    case MNode:
        child = n.ch[names[nodeIndex]]
    case LNode:
        child = n.ch[nodeIndex]
    default:
        child = nil
    }
    return child
}

func PrintNode(n Node) string {
    var b strings.Builder
    var fu =
    func(
        w    io.StringWriter,
        m    int,
        args ...interface{},
    ) error {
        for i :=0 ; i != m ; i++ {
            w.WriteString(" ")
        }
        for _, a := range args {
            w.WriteString(a.(string))
        }
        return nil
    }
    n.app(&b, 0, fu)
    return b.String()
}


