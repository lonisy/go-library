package gin

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "go-library/str"
    "golang.org/x/net/context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "reflect"
    "strings"
    "syscall"
    "time"
)

type GinHttpServer struct {
    AppName string
    State   chan os.Signal
    Actions map[string]interface{}
}

func (s *GinHttpServer) Start(srv *http.Server) {
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen%s %s\n", srv.Addr, err)
        }
    }()
    s.State = make(chan os.Signal)
    signal.Notify(s.State, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
    <-s.State
    log.Println("Shutdown Server ...")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server Shutdown: ", err)
    }
    log.Println("Server exiting")
}

func (s *GinHttpServer) Stop() {
    s.State <- syscall.SIGQUIT
}

func (s *GinHttpServer) Register(module interface{}) {
    if len(s.Actions) == 0 {
        s.Actions = make(map[string]interface{}, 0)
    }
    s.Actions[reflect.TypeOf(module).String()] = module
}

func (s *GinHttpServer) Routing(c *gin.Context) {
    if len(c.Params) >= 1 {
        cc := ""
        for action, _ := range s.Actions {
            cc = strings.ToLower(action)
            cc = strings.Replace(cc, "api.", "", 1)
            cc = strings.Replace(cc, "controler", "", 1)
            cc = strings.Replace(cc, "struct", "", 1)
            cc = strings.Replace(cc, "*", "", 1)
            cv := c.Params[0].Value
            cv = strings.Replace(cv, "_", "", 1)
            cv = strings.Replace(cv, "-", "", 1)
            if cc == cv {
                cc = action
                MethodName := "Index"
                if len(c.Params) >= 2 {
                    MethodName = str.StringToCamel(strings.ReplaceAll(c.Params[1].Value, "-", "_"))
                }
                _, ok := reflect.TypeOf(s.Actions[cc]).MethodByName(MethodName)
                if ok {
                    method := reflect.ValueOf(s.Actions[cc]).MethodByName(MethodName)
                    params := make([]reflect.Value, 1)
                    params[0] = reflect.ValueOf(c)
                    method.Call(params)
                    fmt.Println("auto router called")
                    return
                } else {
                    fmt.Println("router miss")
                    fmt.Println("router default")
                    method := reflect.ValueOf(s.Actions[cc]).MethodByName("Index")
                    params := make([]reflect.Value, 1)
                    params[0] = reflect.ValueOf(c)
                    if method.IsValid() {
                        method.Call(params)
                    }
                    return
                }
            }
        }
    }
    c.JSON(http.StatusNotFound, gin.H{
        "code":    -1,
        "message": "NotFound",
    })
    return
}
