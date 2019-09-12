package main

import (
	stdContext "context"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"github.com/sevlyar/go-daemon"

	"github.com/jinil-ha/blind-msg/handler"
	"github.com/jinil-ha/blind-msg/service/slack"
	"github.com/jinil-ha/blind-msg/utils/config"
)

var logf *os.File
var app *iris.Application
var daemonFlag bool

func openLogFile() {
	if !daemonFlag {
		return
	}

	if logf != nil {
		logf.Close()
		logf = nil
	}

	filename := config.GetString("log_file")
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		golog.Fatalf("Cannot open log file(%s)", filename)
		return
	}

	golog.SetOutput(f)
	app.Logger().SetOutput(f)
	logf = f
}

func printUsage() {
	fmt.Printf("Usage: %s COMMAND\n", os.Args[0])
	fmt.Printf("COMMAND:\n")
	fmt.Printf("\tstart\tstart as daemon.\n")
	fmt.Printf("\tstop\tstop daemon.\n")
	fmt.Printf("\thup\tsend SIGHUP signal to daemon(reopen log file).\n")
	fmt.Printf("\trun\tstart as console program.\n")
	os.Exit(1)
}

func termHandler(sig os.Signal) error {
	golog.Error("Server shutdown")

	timeout := 5 * time.Second
	ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
	defer cancel()
	app.Shutdown(ctx)

	return daemon.ErrStop
}

func hupHandler(sig os.Signal) error {
	golog.Warn("SIGHUP received")

	openLogFile()
	return nil
}

func worker() {
	// set log level
	golog.SetTimeFormat("2006/01/02_15:04:05")

	level := config.GetString("log_level")
	golog.SetLevel(level)
	app.Logger().SetLevel(level)

	golog.Error("Server starts")

	template := iris.HTML("./resource/html", ".html")
	template.Reload(true)
	app.RegisterView(template)

	// home page
	app.Get("/", handler.TopHandler)

	// auth
	app.Get("/goauth", handler.GoAuthHandler)
	app.Get("/auth", handler.AuthHandler)
	app.Get("/logout", handler.LogoutHandler)

	// send page
	app.Get("/send", handler.GetSendHandler)
	app.Post("/send", handler.PostSendHandler)

	listenAddr := config.GetString("listen_addr")
	slack.Start()

	app.Run(iris.Addr(listenAddr), iris.WithoutServerError(iris.ErrServerClosed))
}

func sendSignal(ctx *daemon.Context, sig os.Signal) {
	d, err := ctx.Search()
	if err != nil {
		golog.Fatalf("Cannot search daemon process: %s", err.Error())
		return
	}

	if err = d.Signal(sig); err != nil {
		golog.Fatalf("Fail to send %s to the daemon: %s", sig, err.Error())
	}
}

func main() {
	// check argument
	if len(os.Args) < 2 {
		printUsage()
	}

	// set daemon context
	ctx := &daemon.Context{
		PidFileName: config.GetString("pid_file"),
		PidFilePerm: 0644,
		LogFileName: "",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        nil,
	}

	switch os.Args[1] {
	case "start":
		child, err := ctx.Reborn()
		if err != nil {
			panic(err)
		}

		if child != nil {
			// post parent
		} else {
			defer ctx.Release()
			daemonFlag = true

			// set signal handler
			daemon.SetSigHandler(termHandler, syscall.SIGTERM)
			daemon.SetSigHandler(hupHandler, syscall.SIGHUP)

			// start server logic
			app = iris.Default()
			openLogFile()

			go worker()

			// process signals
			err = daemon.ServeSignals()
			if err != nil {
				golog.Error(err)
			}

			golog.Error("Server terminated")
			if logf != nil {
				logf.Close()
			}
		}

	case "stop":
		sendSignal(ctx, syscall.SIGTERM)

	case "hup":
		sendSignal(ctx, syscall.SIGHUP)

	case "run":
		daemonFlag = false
		app = iris.Default()
		worker()

	default:
		printUsage()
	}
}
