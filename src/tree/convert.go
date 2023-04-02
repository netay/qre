package tree

import (
    "fmt"
    "errors"
)

// Converts qre text to tree node
func Convert (i interface{}) (Node, error) {
    switch v := i.(type) {
    case string:
        return Leaf {value: v}, nil
    case []interface{}:
        var ch = make([]Node, 0, len(v))
        for _, x := range v {
            y, e := Convert(x)
            if e != nil {
                return nil, e
            }
            ch = append(ch, y)
        }
        return LNode {ch: ch}, nil
    case map[string]interface{}:
        var ch = make(map[string]Node)
        for k, x := range v {
            y, e := Convert(x)
            if e != nil {
                return nil, e
            }
            ch[k] = y
        }
        return MNode {ch: ch}, nil
    case map[interface{}]interface{}:
        var ch = make(map[string]Node)
        for k, x := range v {
            y, e := Convert(x)
            if e != nil {
                return nil, e
            }
            sk := fmt.Sprintf("%v", k)
            ch[sk] = y
        }
        return MNode {ch: ch}, nil
    default:
        fmt.Printf("Wrong data: '%v'\n", v)
        //return nil, errors.New("Parsing error!")
    }
    return nil, errors.New("Parsing error!")
}

