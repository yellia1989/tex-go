package util

import (
    "os"
    "io"
    "bytes"
)

func LoadFromFile(file string) ([]byte, error) {
    f, err := os.Open(file)
    if err != nil {
        return nil, err
    }

    defer f.Close()

    buffer := bytes.Buffer{}
    tmp := make([]byte, 1024)
    for {
        n, err := f.Read(tmp)
        if err != nil {
            if err == io.EOF {
                break
            }
            return nil, err
        }
        buffer.Write(tmp[:n])
    }

    return buffer.Bytes(), nil
}

func SaveToFile(file string, content []byte, append bool) error {
    mode := os.O_CREATE|os.O_WRONLY
    if append {
        mode |= os.O_APPEND
    } else {
        mode |= os.O_TRUNC
    }
    f, err := os.OpenFile(file, mode, 0644)
    if err != nil {
        return err
    }

    defer f.Close()

    _, err = f.Write(content)
    if err != nil {
        return err
    }
    return nil
}

func Mkdir(path string) error {
    if err := os.MkdirAll(path, 0775); err != nil {
        return err
    }

    return nil
}
