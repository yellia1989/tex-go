package tex

import (
    "strings"
    "strconv"
    "fmt"
    "time"
    "github.com/yellia1989/tex-go/tools/util"
)

type Endpoint struct {
    Proto string
    IP string
    Port int
    Idletimeout time.Duration
}

func (ep *Endpoint) String() string {
    return ep.Proto + " -h " + ep.IP + " -p " + strconv.Itoa(ep.Port) + " -t " + strconv.FormatInt(ep.Idletimeout.Milliseconds(), 10)
}

func (ep *Endpoint) Address() string {
    return ep.IP + ":" + strconv.Itoa(ep.Port)
}

// tcp -h 127.0.0.1 -p 8080 -t 3600000
func NewEndpoint(endpoint string) (*Endpoint,error) {
    tokens := strings.Fields(endpoint)

    if len(tokens) != 7 {
        return nil, fmt.Errorf("invalid endpoint")
    }

    var err error
    ep := &Endpoint{}
    ep.Proto = tokens[0]
    ep.IP = tokens[2]
    if ep.Port,err = strconv.Atoi(tokens[4]); err != nil {
        panic("invalid endpoint")
    }
    ep.Idletimeout = util.AtoDuration(tokens[6]+"ms")
    return ep,nil
}
