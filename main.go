package main

import (
    "net/http"
    "log"
    "github.com/gorilla/mux"
    "github.com/xuanbo/FileServer/controller"
)

func main()  {
    // 初始化Router
    router := newRouter()

    // 注册HandleFunc
    registerRoutes(router)

    // 运行
    run(router)
}

func newRouter() *mux.Router {
    return mux.NewRouter()
}

func registerRoutes(router *mux.Router)  {
    // static
    router.PathPrefix("/static/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(""))))

    // json
    router.HandleFunc("/public", controller.PublicHandler).
        Methods("GET").Headers("Content-Type", "application/json")
    router.PathPrefix("/public").
        Subrouter().HandleFunc("/{relativePath:[^\\s]*}", controller.FileWalkerHandler).
        Methods("GET").Headers("Content-Type", "application/json")

    // html
    router.HandleFunc("/public", controller.PublicHandlerHtml).
        Methods("GET")
    router.PathPrefix("/public").
        Subrouter().HandleFunc("/{relativePath:[^\\s]*}", controller.FileWalkerHandlerHtml).
        Methods("GET")
    router.PathPrefix("/download").
        Subrouter().HandleFunc("/{relativePath:[^\\s]*}", controller.DownloadFile).
        Methods("GET")
}

func run(router *mux.Router) {
    // Routes consist of a path and a handler function.
    http.Handle("/", router)
    // Bind to a port and pass our router in
    log.Println("服务已启动...")
    log.Println("监听端口8080...")
    log.Fatal(http.ListenAndServe(":8080", router))
}