package gin

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
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

type HttpServer struct {
	*gin.Engine
	AppName string
	Actions map[string]interface{}
	State   chan os.Signal
}

func NewHttpServer(appName string) *HttpServer {
	return &HttpServer{
		AppName: appName,
		Actions: make(map[string]interface{}, 0),
		State:   make(chan os.Signal),
	}
}

func (r *HttpServer) MvcRouter() *HttpServer {
	r.Static("/assets", "./dist/assets")
	r.StaticFile("/favicon.png", "./dist/favicon.png")
	r.LoadHTMLGlob("./dist/*.html")
	return r
}

func (r *HttpServer) AnyApiRouter() *HttpServer {
	apiGroups := r.Group("/api/v1")
	apiGroups.Any("/:arg1/:arg2/:arg3/:arg4", r.RouteMatch)
	apiGroups.Any("/:arg1/:arg2/:arg3", r.RouteMatch)
	apiGroups.Any("/:arg1/:arg2", r.RouteMatch)
	apiGroups.Any("/:arg1", r.RouteMatch)
	return r
}

func (r *HttpServer) UseGzip() *HttpServer {
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	return r
}

func (r *HttpServer) UseCORS() *HttpServer {
	r.Use(CORSMiddleware())
	return r
}

func (r *HttpServer) UseApiNotFound() *HttpServer {
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    -1,
			"message": "NotFound",
		})
	})
	return r
}

func (s *HttpServer) Start(srv *http.Server) {
	go func() {
		log.Println("Start Server ...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen%s %s\n", srv.Addr, err)
		}
	}()
	s.State = make(chan os.Signal)
	signal.Notify(s.State, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-s.State
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}
	log.Println("Server exiting...")
}

func (s *HttpServer) Stop() {
	s.State <- syscall.SIGQUIT
}

func (s *HttpServer) Register(module interface{}) {
	if len(s.Actions) == 0 {
		s.Actions = make(map[string]interface{}, 0)
	}
	s.Actions[reflect.TypeOf(module).String()] = module
}

func (s *HttpServer) RouteMatch(c *gin.Context) {
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
					MethodName = stringToCamel(strings.ReplaceAll(c.Params[1].Value, "-", "_"))
				}
				_, ok := reflect.TypeOf(s.Actions[cc]).MethodByName(MethodName)
				if ok {
					method := reflect.ValueOf(s.Actions[cc]).MethodByName(MethodName)
					params := make([]reflect.Value, 1)
					params[0] = reflect.ValueOf(c)
					method.Call(params)
					log.Printf("router %s.%s called", cc, MethodName)
					return
				} else {
					log.Printf("router %s.%s miss", cc, MethodName)
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

func stringToCamel(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}
