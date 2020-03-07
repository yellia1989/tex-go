package util

import (
    "testing"
)

func TestConfig(t *testing.T) {
    defer func() {
        if err := recover(); err != nil {
            t.Fatalf("load cfg file err:%s", err)
        }
    }()

    c := NewConfig()
    c.ParseFile("./testdata/test.cfg")
    t.Logf("%s", c)

    t.Run("GetCfg", func (t *testing.T) {
        mkey := map[string]string{
            "test1": "test1",
            "mfw/application/client/locator": "tex.mfwregistry.QueryObj@TEMPLATE_LOCATOR",
            "mfw/application/server/Service_1/endpoint": "TEMPLATE_ENDPOINT_GameServiceObj",
        }
        for k, v := range mkey {
            v2 := c.GetCfg(k, "default")
            if v != v2 {
                t.Fatalf("%s != %s, real:%s", k, v, v2)
            }
        }
    })
}

func TestAtoMB(t *testing.T) {
    msize := map[string]uint64{
        "1024KB": 1,
        "100MB": 100,
        "100M": 100,
        "1GB": 1024,
        "1TB": 1024 * 1024,
        "1PB": 1024 * 1024 * 1024,
        "1EB": 1024 * 1024 * 1024 * 1024,
    }
    for k,v := range msize {
        size := AtoMB(k)
        if size != v {
            t.Fatalf("expect:%d, real:%d", v, size)
        }
    }
}
