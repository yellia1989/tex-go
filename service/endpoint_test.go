package tex

import (
    "testing"
    "time"
)

func TestEndpoint(t *testing.T) {
    ep := &Endpoint{}
    ep.Proto = "tcp"
    ep.IP = "127.0.0.1"
    ep.Port = 8080
    ep.Idletimeout = 3600000 * time.Millisecond
    m := map[string]*Endpoint{
        "tcp -h 127.0.0.1 -p 8080 -t 3600000":ep,
        "tcp  -h 127.0.0.1 -p 8080 -t 3600000":ep,
    }

    for k, v := range m {
        real, err := NewEndpoint(k)
        if err != nil {
            t.Fatalf("%s\n", err.Error())
        }
        if v.String() != real.String() {
            t.Fatalf("real:%s, expect=%s", real, v)
        }
    }
}
