package tex

import (
    "sync"
)

type Communicator struct {
    mu sync.Mutex
    mPrx map[string]*servicePrxImpl
}

func NewCommunicator() *Communicator {
    comm := &Communicator{}
    comm.mPrx = make(map[string]*servicePrxImpl)

    return comm
}

func (comm *Communicator) StringToProxy(name string, prx ServicePrx) error {
    comm.mu.Lock()

    if impl, ok := comm.mPrx[name]; ok {
        comm.mu.Unlock()
        prx.SetPrxImpl(impl)
        return nil
    }

    comm.mu.Unlock()
    impl, err := newPrxImpl(name, comm)
    if err != nil {
        return err
    }
    comm.mu.Lock()
    comm.mPrx[name] = impl
    comm.mu.Unlock()
    impl.SetTimeout(cliCfg.invokeTimeout)
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
