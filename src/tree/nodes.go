package tree

import (
    "io"
    "sort"
    "strings"
)

type Node interface {
    app(
        io.StringWriter, // where to write
        io.StringWriter, // where to write
        int,              // margin, depth of the node
        func(
            io.StringWriter,
            io.StringWriter,
            int,
            interface{},
            interface{},
        ) error,
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
    return MNode {ch: ch}
}
func NewLNode(ch []Node) LNode {
    return LNode {ch: ch}
}

func (l Leaf) app(
    w1 io.StringWriter,
    w2 io.StringWriter,
    m int,
    f func(
        io.StringWriter,
        io.StringWriter,
        int,
        interface{},
        interface{},
    ) error,
) {
    f(w1, w2, m, l.value + "\n", "\n")
}
func (n MNode) app(
    w1 io.StringWriter,
    w2 io.StringWriter,
    m int,
    f func(
        io.StringWriter,
        io.StringWriter,
        int,
        interface{},
        interface{},
    ) error,
) {
    for key, value := range n.ch {
        switch v := value.(type) {
        case Leaf:
            f(w1, w2, m, key + "\n", (v.value + "\n"))
        default:
            f(w1, w2, m, key + ":\n", "\n")
            value.app(w1, w2, m + 1, f)
        }
    }
}
func (l LNode) app(
    w1 io.StringWriter,
    w2 io.StringWriter,
    m int,
    f func(
        io.StringWriter,
        io.StringWriter,
        int,
        interface{},
        interface{},
    )error,
) {
    for _, n := range l.ch {
        n.app(w1, w2, m + 1, f)
    }
}
func (s SNode) app(
    w1 io.StringWriter,
    w2 io.StringWriter,
    m int,
    f func(
        io.StringWriter,
        io.StringWriter,
        int,
        interface{},
        interface{},
    ) error,
) {
    f(w1, w2, m, s.name + "\n", "\n")
    s.ch.app(w1, w2, m, f)
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

func PrintNode(n Node) (string, string) {
    var bleft strings.Builder
    var bright strings.Builder
    var fu =
    func(
        w1   io.StringWriter,
        w2   io.StringWriter,
        m    int,
        lline interface{},
        rline interface{},
    ) error {
        for i := 0 ; i != m ; i++ {
            w1.WriteString(" ")
        }
        if w1 != nil {
            w1.WriteString(lline.(string))
        }
        if w2 != nil {
            w2.WriteString(rline.(string))
        }
        return nil
    }
    n.app(&bleft, &bright, 0, fu)
    return bleft.String(), bright.String()
}


