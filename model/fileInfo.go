package model

import (
    "time"
)

// 文件信息
type FileInfo struct {
    Name string `json:"name"`
    Size int64 `json:"size"`
    Type string `json:"type"`
    ModifyDate time.Time `json:"modify_date"`
}