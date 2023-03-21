
package main

import (
	"bufio"
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/kmcsr/go-logger"
	"github.com/kmcsr/go-logger/golog"
)

var loger logger.Logger = getLogger()

func getLogger()(loger logger.Logger){
	loger = golog.Logger
	loger.SetLevel(logger.TraceLevel)
	loger.Debug("Debug mode on")
	golog.Unwrap(loger).Logger.SetTimeFormat("2006-01-02 15:04:05.000:")
	return
}

var (
	Target string = "./vue-project/dist"
	host string = "127.0.0.1"
	port int = 3080
)

func main(){
	username := os.Getenv("DB_USER")
	if username == "" {
		os.Setenv("DB_USER", "root")
	}
	// passwd := os.Getenv("DB_PASSWD")
	// address := os.Getenv("DB_ADDR")
	database := os.Getenv("DB_NAME")
	if database == "" {
		os.Setenv("DB_NAME", "databasename")
	}

	apihostportch := make(chan string, 1)

	cmdctx, cancel := context.WithCancel(context.Background())
	apicmd := exec.CommandContext(cmdctx, "go", "run", "./handlers/dev", "127.0.0.1:")
	defer func(cancel context.CancelFunc){
		cancel()
		if apicmd.Process != nil {
			apicmd.Wait()
		}
	}(cancel)

	go func(){
		apiout, err := apicmd.StdoutPipe()
		if err != nil {
			loger.Fatalf("Cannot make stdout pipe for api: %v", err)
		}
		apicmd.Stderr = apicmd.Stdout
		if err = apicmd.Start(); err != nil {
			loger.Fatalf("Cannot start api: %v", err)
		}
		sc := bufio.NewScanner(apiout)
		waitingflag := true
		for sc.Scan() {
			txt := sc.Text()
			os.Stdout.Write(([]byte)(txt + "\n"))
			if waitingflag {
				const serverListeneingAtPrefix = "API]: Server listening at "
				if i := strings.Index(txt, serverListeneingAtPrefix); i >= 0 {
					waitingflag = false
					hostport := txt[i + len(serverListeneingAtPrefix):]
					apihostportch <- hostport
				}
			}
		}
		if err = apicmd.Wait(); err != nil {
			loger.Fatalf("Api process exited with error: %v", err)
		}
	}()

	hostport := <-apihostportch
	loger.Infof("Detected api host-port: <%s>", hostport)
	apiurl, err := url.Parse("http://" + hostport)
	if err != nil {
		loger.Fatalf("Cannot parse api host-port: %v", err)
	}

	http.Handle("/dev/", http.StripPrefix("/dev", httputil.NewSingleHostReverseProxy(apiurl)))
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
