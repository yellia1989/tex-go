package main

import (
    "context"
)

type EchoServiceImpl struct {
}

func (s *EchoServiceImpl) Hello(ctx context.Context, req string, resp *string) error {
    *resp = req

    return nil
}
