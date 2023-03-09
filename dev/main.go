
package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/kmcsr/go-logger"
	"github.com/kmcsr/go-logger/golog"
	"github.com/kmcsr/PluginWebPoint/api"
)

var loger logger.Logger = getLogger()

func getLogger()(loger logger.Logger){
	loger = golog.Logger
	golog.Unwrap(loger).Logger.SetTimeFormat("2006-01-02 15:04:05.000:")
	return
}

var (
	DEBUG bool = false
	Target string = "./vue-project/dist"
	host string = "127.0.0.1"
	port int = 3080
)

func main(){
	loger.SetLevel(logger.TraceLevel)
	loger.Trace("Debug mode on")
	api.SetLogger(loger)

	username := os.Getenv("DB_USER")
	if username == "" {
		username = "root"
	}
	passwd := os.Getenv("DB_PASSWD")
	address := os.Getenv("DB_ADDR")
	database := os.Getenv("DB_NAME")
	if database == "" {
		database = "databasename"
	}

	api.Ins = api.NewMySqlAPI(username, passwd, address, database)

	http.Handle("/dev/", http.StripPrefix("/dev", api.GetDevAPIHandler()))
	http.Handle("/assets/", http.StripPrefix("/assets",
		http.FileServer(http.Dir(filepath.Join(Target, "assets")))))
	http.Handle("/", (http.HandlerFunc)(func(rw http.ResponseWriter, req *http.Request){
		http.ServeFile(rw, req, filepath.Join(Target, "index.html"))
	}))

	server := &http.Server{
		Addr: net.JoinHostPort(host, strconv.Itoa(port)),
		Handler: http.DefaultServeMux,
		ReadTimeout: time.Second * 30,
		WriteTimeout: time.Second * 60,
	}

	done := make(chan struct{}, 0)
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)

	go func(){
		defer close(done)
		loger.Infof("Server start at %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			loger.Fatal(err)
		}
	}()

	select {
	case sig := <-sigch:
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 3)
		_ = sig // TODO: reload config
		server.Shutdown(ctx)
		cancel()
	case <-done:
	}
}
