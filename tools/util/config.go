package util

import (
    "strings"
    "bytes"
    "strconv"
    "time"
)

type Config struct {
    sName string    // 节点的名字
    vItem []string  // 节点中的所有子节点名字(叶子节点和子节点)
    mKeyValue map[string]string // 叶子节点
    mSubConfig map[string]*Config // 子节点
}

func NewConfig() *Config {
    c := &Config{}
    c.mKeyValue = make(map[string]string)
    c.mSubConfig = make(map[string]*Config)
    return c
}

func (c *Config) ParseString(content string) {
    lines := strings.Split(content, "\n")

    s := NewStack()
    s.Push(c)
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if len(line) == 0 || line[0] == '#' {
            continue
        }
        if line[0] == '<' {
            total := len(line)
            if len(line) < 3 || line[total-1] != '>' {
                panic("parse err:" + line)
            }

            top := s.Top().(*Config)
            if line[1] == '/' {
                // 当前节点结束了
                if line[2:(total-1)] != top.sName {
                    panic("key mismatch:" + line[2:(total-1)])
                }
                s.Pop()
            } else {
                // 子节点开始
                sName := line[1:(total-1)]
                
                if Contain(top.vItem, sName) {
                    panic("duplicate key:"+sName+", in " + top.sName)
                }

                top.vItem = append(top.vItem, sName)
                c := NewConfig()
                c.sName = sName
                top.mSubConfig[sName] = c
                s.Push(c)
            }
        } else {
            top := s.Top().(*Config)
            // 叶子节点
            tmp := strings.SplitN(line, "=", 2)
            key := strings.TrimSpace(tmp[0])
            val := strings.TrimSpace(tmp[1])
            if len(key) == 0 {
                panic("parse err:" + line)
            }

            if Contain(top.vItem, key) {
                panic("duplicate key:"+key+", in " + top.sName)
            }
            top.vItem = append(top.vItem, key)
            top.mKeyValue[key] = val
        }
    }
}

func (c *Config) ParseFile(file string) {
    content, err := LoadFromFile(file)
    if err != nil {
        panic(err.Error())
    }
    c.ParseString(string(content))
}

func (c *Config) toString(indent int) string {
    sTab := strings.Repeat(" ", indent)

    var buffer bytes.Buffer
    for _, key := range c.vItem {
        val, ok := c.mKeyValue[key]
        if ok {
            buffer.WriteString(sTab + key + " = " + val + "\n")
            continue
        } 

        c, ok := c.mSubConfig[key]
        if ok {
            buffer.WriteString(sTab+"<"+key+">\n")
            buffer.WriteString(c.toString(indent+4))
            buffer.WriteString(sTab+"</"+key+">\n")
        }
    }
    return buffer.String()
}

func (c *Config) String() string {
    return c.toString(0)
}

func (c *Config) GetCfg(path string, def string) string {
    if path != "" && path[0] == '/' {
        path = path[1:]
    }
    if len(path) == 0 {
        return def
    }

    var node *Config
    var key string
    paths := strings.Split(path, "/")
    if len(paths) == 1 {
        node = c
        key = paths[0]
    } else {
        node = c.getSubCfg(paths[:(len(paths)-1)])
        key = paths[len(paths)-1]
    }
    if node == nil {
        return def
    }
    val, ok := node.mKeyValue[key]
    if ok {
        return val
    }
    return def
}

func (c *Config) GetSubCfg(path string) *Config {
    if len(path) == 0 {
        return nil
    }

    paths := strings.Split(path, "/")
    return c.getSubCfg(paths)
}

func (c *Config) getSubCfg(paths []string) *Config {
    if len(paths) == 0 {
        return nil
    }
    c, ok := c.mSubConfig[paths[0]]
    if !ok {
        return nil
    }
    if len(paths) == 1 {
        return c
    }
    return c.getSubCfg(paths[1:])
}

func (c *Config) GetInt(key string, def int) int {
    v := c.GetCfg(key, strconv.Itoa(def))
    i,_ := strconv.Atoi(v)
    return i
}

func (c *Config) GetBool(key string, def bool) bool {
    sdef := "0"
    if def {
        sdef = "1"
    }
    v := c.GetCfg(key, sdef)
    i,_ := strconv.Atoi(v)
    return i != 0
}

func (c *Config) GetDuration(key string, def string) time.Duration {
    v := c.GetCfg(key, def)
    return AtoDuration(v)
}
