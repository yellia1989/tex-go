package util

import (
    "container/list"
)

type Stack struct {
    l *list.List
}

func (s *Stack) Push(v interface{}) {
    s.l.PushBack(v)
}

func (s *Stack) Pop() {
    if s.l.Len() == 0 {
        panic("stack empty")
    }
    s.l.Remove(s.l.Back())
}

func (s *Stack) Top() interface{} {
    ele := s.l.Back()
    if ele != nil {
        return ele.Value
    }
    return nil
}

func (s *Stack) Len() int {
    return s.l.Len()
}

func NewStack() *Stack {
    s := &Stack{l:list.New()}
    return s
}
