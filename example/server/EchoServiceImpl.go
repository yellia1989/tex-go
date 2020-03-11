package main

type EchoServiceImpl struct {
}

func (s *EchoServiceImpl) Hello(req string, resp *string) error {
    *resp = req

    return nil
}
