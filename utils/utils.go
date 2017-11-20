package utils

import (
    "bytes"
)

func Convert( b []byte ) string {
    buf := bytes.NewBuffer(b)
    return buf.String()
}