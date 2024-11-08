package main

import (
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"

    "github.com/gin-gonic/gin"
)

// Конфигурация маршрутов к микросервисам
var services = map[string]string{
    "/api/service1": "http://localhost:8081",
    "/api/service2": "http://localhost:8082",
}

func main() {
    router := gin.Default()

    // Middleware для логирования и восстановления после паники
    router.Use(gin.Logger())
    router.Use(gin.Recovery())

    // Обработчик для всех маршрутов
    router.Any("/*proxyPath", proxyHandler)

    log.Println("Gateway запущен на порту 8000...")
    router.Run(":8000")
}

func proxyHandler(c *gin.Context) {
    path := c.Request.URL.Path

    // Определяем, к какому сервису направить запрос
    for prefix, target := range services {
        if matchPrefix(path, prefix) {
            remoteURL, err := url.Parse(target)
            if err != nil {
                c.String(http.StatusBadGateway, "Bad Gateway")
                return
            }

            // Создаем прокси
            proxy := httputil.NewSingleHostReverseProxy(remoteURL)

            // Переписываем путь запроса
            c.Request.URL.Path = singleJoiningSlash(remoteURL.Path, path[len(prefix):])

            // Изменяем хост запроса
            c.Request.Host = remoteURL.Host
            c.Request.URL.Host = remoteURL.Host
            c.Request.URL.Scheme = remoteURL.Scheme

            // Проксируем запрос
            proxy.ServeHTTP(c.Writer, c.Request)
            return
        }
    }
    c.String(http.StatusNotFound, "404 not found")
}

func matchPrefix(path, prefix string) bool {
    if path == prefix {
        return true
    }
    if len(path) > len(prefix) && path[:len(prefix)+1] == prefix+"/" {
        return true
    }
    return false
}

func singleJoiningSlash(a, b string) string {
    aslash := a[len(a)-1] == '/'
    bslash := b[0] == '/'
    switch {
    case aslash && bslash:
        return a + b[1:]
    case !aslash && !bslash:
        return a + "/" + b
    }
    return a + b
}
