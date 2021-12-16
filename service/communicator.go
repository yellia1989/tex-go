package tex

import (
    "sync"
    "fmt"
)

type Communicator struct {
    mu sync.Mutex
    mPrx map[string]*servicePrxImpl
    sLocator string
}

func NewCommunicator(locator string) *Communicator {
    comm := &Communicator{}
    comm.mPrx = make(map[string]*servicePrxImpl)
    comm.sLocator = locator

    return comm
}

func (comm *Communicator) StringToProxy(name string, prx ServicePrx) error {
    if name == "" {
        return fmt.Errorf("service obj name required")
    }
    comm.mu.Lock()

    if impl, ok := comm.mPrx[name]; ok {
        comm.mu.Unlock()
        prx.SetPrxImpl(impl)
        return nil
    }

    comm.mu.Unlock()
    impl := newPrxImpl(name, comm)
    comm.mu.Lock()
    if impl, ok := comm.mPrx[name]; ok {
        comm.mu.Unlock()
        prx.SetPrxImpl(impl)
        return nil
    }
    comm.mPrx[name] = impl
    comm.mu.Unlock()
    prx.SetPrxImpl(impl)
    return nil
}

func (comm *Communicator) Close() {
    comm.mu.Lock()
    defer comm.mu.Unlock()

    for _, v := range comm.mPrx {
        v.close()
    }
}
