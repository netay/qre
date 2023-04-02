package qre

import (
    "errors"
    "gopkg.in/yaml.v2"
    "encoding/json"
    "github.com/pelletier/go-toml"
    "qre/tree"
)

func parseTree(input []byte, ft int) (tree.Node, error) {
    var p map[string]interface{}
    var err error
    switch ft {
    case qre_json:
        err = json.Unmarshal(input, &p)
    case qre_yaml:
        err = yaml.Unmarshal(input, &p)
    case qre_toml:
        err = toml.Unmarshal(input, &p)
    default:
        return nil, errors.New("not implemented")
    }
    if err != nil {
        return nil, err
    }
    m := make(map[string]tree.Node)
    for k, x := range p {
        y, e := tree.Convert(x)
        if e != nil {
            return nil, e
        }
        m[k] = y
    }
    var n = tree.NewMNode(m)
    return n, nil
}
