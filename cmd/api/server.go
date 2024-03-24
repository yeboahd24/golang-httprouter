package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
  "context"
  "errors"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

  shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		app.logg.PrintInfo("caught signal", map[string]string{
			"signal": s.String(),
		})

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

    defer cancel()

    err := srv.Shutdown(ctx)

    if err != nil {
      shutdownError <- err
    }

    app.logg.PrintInfo("completing background tasks", map[string]string{
      "addr": srv.Addr,
    })

    app.wg.Wait()

    shutdownError <- nil

	}()

	app.logg.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})

  err := srv.ListenAndServe()

  if !errors.Is(err, http.ErrServerClosed) {
    return err
  }

  err = <-shutdownError

  if err != nil {
    return err
  }

  app.logg.PrintInfo("stopped server", map[string]string{
    "addr": srv.Addr,
  })



  return nil

}

