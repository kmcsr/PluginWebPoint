
package main

import (
	"context"
	"crypto/sha256"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"sync"
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

//go:embed dist
var dist embed.FS

var startTime = time.Now()

var (
	DEBUG bool = false
	host string = ""
	port int = 80
)

func init(){
	golog.Unwrap(loger).Logger.SetTimeFormat("2006-01-02 15:04:05.000:")
}

func parseArgs(){
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-debug":
			DEBUG = true
		}
	}
}

func main(){
	parseArgs()

	if DEBUG {
		loger.SetLevel(logger.TraceLevel)
	}else{
		loger.SetLevel(logger.InfoLevel)
	}
	loger.Debug("Debug mode on")
	api.SetLogger(loger)

	username := os.Getenv("DB_USER")
	passwd := os.Getenv("DB_PASSWD")
	address := os.Getenv("DB_ADDR")
	database := os.Getenv("DB_NAME")

	api.Ins = api.NewMySqlAPI(username, passwd, address, database)

	var (
		err error
		assetsFS fs.FS
		indexFile io.ReadSeeker
	)
	if assetsFS, err = fs.Sub(dist, "dist/assets"); err != nil {
		loger.Fatalf("Couldn't load dist/assets: %v", err)
	}
	if fd, err := dist.Open("dist/index.html"); err != nil {
		loger.Fatalf("Couldn't load dist/index.html: %v", err)
	}else{
		indexFile = fd.(io.ReadSeeker)
	}

	http.Handle("/dev/", http.StripPrefix("/dev", api.GetDevAPIHandler()))
	http.Handle("/assets/", http.StripPrefix("/assets", NewConstFileServer(assetsFS, time.Time{})))
	http.Handle("/", HandleConstData(indexFile, time.Time{}, "index.html"))

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

func HandleConstData(data io.ReadSeeker, modTime time.Time, name string)(http.Handler){
	h := sha256.New()
	_, err := data.Seek(0, io.SeekStart)
	if err != nil {
		loger.Panic(err)
	}
	if _, err = io.Copy(h, data); err != nil {
		loger.Panic(err)
	}
	etag := fmt.Sprintf(`"sha256:%x"`, h.Sum(nil))
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request){
		loger.Debugf("Conn from %s;%s;%s", req.RemoteAddr, req.Method, req.URL.RawPath)
		if req.Method != "GET" {
			http.Error(rw, "Request method not allowed, only allows GET", http.StatusMethodNotAllowed)
			return
		}
		rw.Header().Set("Etag", etag)
		http.ServeContent(rw, req, name, modTime, data)
	})
}

type ConstFileServer struct{
	root fs.FS
	modTime time.Time
	mux sync.RWMutex
	etags map[string]string
}

func NewConstFileServer(root fs.FS, modTime time.Time)(*ConstFileServer){
	return &ConstFileServer{
		root: root,
		modTime: modTime,
		etags: make(map[string]string),
	}
}

func (s *ConstFileServer)ServeHTTP(rw http.ResponseWriter, req *http.Request){
	loger.Debugf("Conn from %s;%s;%s", req.RemoteAddr, req.Method, req.URL.Path)
	if req.Method != "GET" {
		http.Error(rw, "Request method not allowed, only allows GET", http.StatusMethodNotAllowed)
		return
	}
	name := req.URL.Path
	if !strings.HasPrefix(name, "/") {
		name = "/" + name
		req.URL.Path = name
	}
	name = path.Clean(name)
	fd, err := s.root.Open(name[1:])
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, fs.ErrPermission) {
			code = http.StatusForbidden
		}else if errors.Is(err, fs.ErrNotExist) {
			code = http.StatusNotFound
		}else{
			loger.Errorf("%s:%d: %v", req.URL.Path, code, err)
		}
		http.Error(rw, fmt.Sprintf("%d: %s", code, http.StatusText(code)), code)
		return
	}
	s.mux.RLock()
	etag, ok := s.etags[name]
	s.mux.RUnlock()
	if !ok {
		h := sha256.New()
		if _, err := io.Copy(h, fd); err == nil {
			etag = fmt.Sprintf(`"sha256:%x"`, h.Sum(nil))
			s.mux.Lock()
			s.etags[name] = etag
			s.mux.Unlock()
		}else{
			loger.Errorf("%s: Cannot calc etag: %v", req.URL.Path, err)
		}
	}
	if len(etag) > 0 {
		rw.Header().Set("Etag", etag)
	}
	http.ServeContent(rw, req, name, s.modTime, fd.(io.ReadSeeker))
}
