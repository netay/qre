package qre

import (
    "errors"
    "gopkg.in/yaml.v2"
    "encoding/json"
    "github.com/pelletier/go-toml"
    "qre/tree"
)

func parseTreeMap(input []byte, ftype int) (tree.Node, error) {
    var p map[string]interface{}
    var e error
    switch ftype {
    case qre_json:
        e = json.Unmarshal(input, &p)
    case qre_yaml:
        e = yaml.Unmarshal(input, &p)
    case qre_toml:
        e = toml.Unmarshal(input, &p)
    default:
        return nil, errors.New("not implemented")
    }
    if e != nil {
        return nil, e
    }
    m := make(map[string]tree.Node)
    for k, x := range p {
        y, e := tree.Convert(x)
        if e != nil {
            return nil, e
        }
        m[k] = y
    }
    node := tree.NewMNode(m)
    return node, nil
}

func parseTreeList(input []byte, ftype int) (tree.Node, error) {
    var l []interface{}
    var e error
    switch ftype {
    case qre_json:
        e = json.Unmarshal(input, &l)
    case qre_yaml:
        e = yaml.Unmarshal(input, &l)
    case qre_toml:
        e = toml.Unmarshal(input, &l)
    default:
        return nil, errors.New("not implemented")
    }
    if e != nil {
        return nil, e
    }
    m := make([]tree.Node, 0, len(l))
    for _, x := range l {
        y, e := tree.Convert(x)
        if e != nil {
            return nil, e
        }
        m = append(m, y)
    }
    node := tree.NewLNode(m)
    return node, nil
}

func parseTree(input []byte, ftype int) (tree.Node, error) {
    var node tree.Node = nil
    var err error = nil
    node, err = parseTreeMap(input, ftype)
    if err == nil {
        return node, nil
    }
    node, err =parseTreeList(input, ftype)
    if err == nil {
        return node, nil
    }
    return nil, err
}
