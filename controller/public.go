package controller

import (
    "net/http"
    "io/ioutil"
    "log"
    "path"
    "io"
    "os"
    "errors"
    "strings"
    "encoding/json"
    "html/template"
    "github.com/xuanbo/FileServer/model"
    "github.com/gorilla/mux"
    "github.com/xuanbo/FileServer/conf"
)



// 显示根目录下面的文件
func PublicHandler(w http.ResponseWriter, r *http.Request) {
    fileInfos, err := listDir(conf.ROOT_DIR)
    if err != nil {
        log.Printf("错误信息[%s]\n", err)
        io.WriteString(w, err.Error())
        return
    }
    data, err := json.Marshal(fileInfos)
    if err != nil {
        log.Printf("FileInfo转json错误，错误信息[%s]\n", err)
        io.WriteString(w, err.Error())
        return
    }
    w.Write(data)
}

// 显示根目录下面的文件
func PublicHandlerHtml(w http.ResponseWriter, r *http.Request) {
    fileInfos, err := listDir(conf.ROOT_DIR)
    if err != nil {
        log.Printf("错误信息[%s]\n", err)
        return
    }
    // 解析模板
    tpl, err := template.ParseFiles("template/fileWalker.html")
    if err != nil {
        log.Printf("错误信息[%s]\n", err)
        return
    }
    // 渲染
    data := &struct {
        CurrentPath string
        RelativePath string
        FileInfos []*model.FileInfo
    }{
        CurrentPath: conf.ROOT_DIR,
        RelativePath: "",
        FileInfos: fileInfos,
    }
    err = tpl.Execute(w, data)
    if err != nil {
        log.Printf("错误信息[%s]\n", err)
    }
}

// 根据url路径遍历目录下的文件
func FileWalkerHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    relativePath := strings.TrimSuffix(vars["relativePath"], "/")
    log.Printf("relativePath: %s\n", relativePath)
    fileInfos, err := listDir(conf.ROOT_DIR + "/" + relativePath)
    if err != nil {
        log.Printf("错误信息[%s]\n", err)
        io.WriteString(w, err.Error())
        return
    }
    data, err := json.Marshal(fileInfos)
    if err != nil {
        log.Printf("FileInfo转json错误，错误信息[%s]\n", err)
        io.WriteString(w, err.Error())
        return
    }
    w.Write(data)
}

func FileWalkerHandlerHtml(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    relativePath := strings.TrimSuffix(vars["relativePath"], "/")
    log.Printf("relativePath: %s\n", relativePath)
    currentPath := conf.ROOT_DIR + "/" + relativePath
    fileInfos, err := listDir(currentPath)
    if err != nil {
        log.Printf("错误信息[%s]\n", err)
        return
    }
    // 解析模板
    tpl, err := template.ParseFiles("template/fileWalker.html")
    if err != nil {
        log.Printf("错误信息[%s]\n", err)
        return
    }
    // 渲染
    data := &struct {
        CurrentPath string
        RelativePath string
        FileInfos []*model.FileInfo
    }{
        CurrentPath: currentPath,
        RelativePath: relativePath,
        FileInfos: fileInfos,
    }
    err = tpl.Execute(w, data)
    if err != nil {
        log.Printf("错误信息[%s]\n", err)
    }
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    relativePath := strings.TrimSuffix(vars["relativePath"], "/")
    log.Printf("relativePath: %s\n", relativePath)
    currentPath := conf.ROOT_DIR + "/" + relativePath
    fileInfo, err := os.Stat(currentPath)
    if err != nil {
        log.Printf("错误信息[%s]\n", err)
        return
    }
    // 目录
    if fileInfo.IsDir() {
        io.WriteString(w, "不是文件")
        return
    }
    // 文件
    content, err := ioutil.ReadFile(currentPath)
    if err != nil {
        log.Printf("错误信息[%s]\n", err)
        return
    }
    w.Write(content)
}

// 遍历目录下的所有文件
func listDir(dir string) ([]*model.FileInfo, error) {
    fileInfo, err := os.Stat(dir)
    if err != nil {
        log.Printf("错误信息[%s]\n", err)
        return nil, err
    }
    // 目录
    if fileInfo.IsDir() {
        return processDir(dir)
    }
    // 文件
    return nil, errors.New("不是目录")
}

func processDir(dir string) ([]*model.FileInfo, error) {
    list, err := ioutil.ReadDir(dir)
    if err != nil {
        log.Printf("错误信息[%s]\n", err)
        return nil, err
    }
    fileInfos := make([]*model.FileInfo, 0)
    for _, file := range list {
        fileName := file.Name()
        var fileType string
        if file.IsDir() {
            fileType = "文件夹"
        } else {
            fileType = strings.TrimPrefix(path.Ext(file.Name()), ".")
        }
        fileInfo := &model.FileInfo{
            Name: fileName,
            Size: file.Size(),
            Type: fileType,
            ModifyDate: file.ModTime(),
        }
        fileInfos = append(fileInfos, fileInfo)
    }
    return fileInfos, nil
}