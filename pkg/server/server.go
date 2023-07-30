package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ReanSn0w/smb-share-unit/pkg/utils"
)

func New(log utils.Logger, port int, cache *utils.Cache, smb utils.SMB) *Server {
	return &Server{
		srv: &http.Server{
			Addr: fmt.Sprintf(":%v", port),
		},
		cache: cache,
		smb:   smb,
		log:   log,
	}
}

type Server struct {
	srv   *http.Server
	cache *utils.Cache
	smb   utils.SMB
	log   utils.Logger
}

func (s *Server) Run() error {
	s.srv.Handler = s.handler()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		s.log.Logf("[INFO] Запуск сервера")
		if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			s.log.Logf("[WARN] Произошла ошибка при запуске сервера %s\n", err.Error())
			quit <- os.Kill
		}
	}()

	registretSignal := <-quit
	s.log.Logf("[DEBUG] Зарегистрирован системный сигнал: %s", registretSignal.String())

	s.log.Logf("[DEBUG] Производится отключение сервера")
	err := s.srv.Shutdown(context.Background())
	if err != nil {
		log.Println(err)
	}

	return nil
}

func (s *Server) handler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			file := r.URL.Path[1:]

			if fileData := s.cache.Get(file); fileData != nil {
				w.WriteHeader(http.StatusOK)
				w.Write(fileData)
				return
			}

			data, err := s.smb.Get(file)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(err.Error()))
				return
			}

			s.cache.Set(file, data)
			w.WriteHeader(http.StatusOK)
			w.Write(data)
		},
	)
}
