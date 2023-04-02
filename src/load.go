package qre

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "os/user"
    "strings"
    "qre/tree"
)

// Source file variants
const (
    not_qre = iota
    qre_json
    qre_yaml
    qre_toml
)

func LoadQres() (tree.Node, error) {
    username, err := user.Current()
    if err != nil {
        log.Fatalf(err.Error())
    }
    qrePath := "/home/" + username.Username + "/.qre/"
    if _, err := os.Stat(qrePath); os.IsNotExist(err) {
        fmt.Printf("no qres found in %v\n", qrePath)
        err := os.Mkdir(qrePath, 0775)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error creating ~/.qre\n")
            return nil, err
        }
        log.Println("Creating ~/.qre/")
    }

    var ch = make(map[string]tree.Node)
    walk(qrePath, ".", &ch)
    return tree.NewMNode(ch), nil
}

// walk directory recursively and
// look for qres, parse and construct tree
func walk(
    pathsrc string,
    nodeprefix string,
    mapdst * map[string]tree.Node,
) error {
    files, err := ioutil.ReadDir(pathsrc)
    if err != nil {
        log.Fatal(err)
    }

    for _, f := range files {
        if f.Mode() & os.ModeSymlink != 0 {
            li, link_err := os.Readlink(pathsrc + "/" + f.Name())
            if link_err != nil {
                fmt.Fprintf(os.Stderr, "Error reading link '%v': %v\n", f, link_err)
                continue
            }
            var link_stat_err error = nil
            if !strings.HasPrefix(li, "/") && !strings.HasPrefix(li, "~") {
                li = pathsrc + "/" + li
            }
            f, link_stat_err = os.Lstat(li)
            if link_stat_err != nil {
                fmt.Fprintf(os.Stderr, "Error following link '%v': %v\n", f, link_stat_err)
                continue
            }
        }
        if f.IsDir() {
            var subdir_map = make(map[string]tree.Node)
            walk(
                pathsrc + "/" + f.Name() + "/",
                nodeprefix + "/" + f.Name(),
                &subdir_map,
            )
            (*mapdst)[nodeprefix + "/" + f.Name() + "/"] = tree.NewMNode(subdir_map)
        } else {
            if ft := file_type(f.Name()); ft != not_qre {
                tree, qre_loading_err := load_qre(pathsrc, f, ft)
                if qre_loading_err != nil {
                    fmt.Fprintf(
                        os.Stderr,
                        "Error loading `%s`: %v\n",
                        f.Name(),
                        qre_loading_err,
                    )
                } else {
                    // extensions .json and .yaml obly
                    nodeName := f.Name()[:len(f.Name()) - 5]
                    (*mapdst)[nodeName] = tree
                }
            } else {
                fmt.Printf("unknown file in ~/.qre: `%s`\n", f.Name())
            }
        }
    }
    return nil
}

func load_qre(dirsrc string, f os.FileInfo, ft int) (tree.Node, error) {
    qre_text, e_read := ioutil.ReadFile(
        dirsrc + "/" + f.Name(),
    )

    if e_read != nil {
        fmt.Fprintf(
            os.Stderr,
            "Error reading `%s`: %v",
            f.Name(),
            e_read,
        )
        return nil, e_read
    }

    tree, e_parse := parseTree(qre_text, ft)
    if e_parse != nil {
        fmt.Fprintf(
            os.Stderr,
            "Error parsing `%s`: %v\n",
            f.Name(),
            e_parse,
        )
        return nil, e_parse
    }

    return tree, nil
}

func file_type(f string) int {
    if strings.HasSuffix(f, ".json") {
        return qre_json
    }
    if strings.HasSuffix(f, ".yaml") {
        return qre_yaml
    }
    if strings.HasSuffix(f, ".toml") {
        return qre_toml
    }
    return not_qre
}

