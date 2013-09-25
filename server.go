package main

import (
  "github.com/coreos/raft"
  "io/ioutil"
  "log"
  "net/http"
  "path"
  "time"
)

type Server struct {
  HTTPServer *http.Server
  Name       string
  Prefix     string
  Raft       *raft.Server
  Queue      *Queue
  NameMap    map[string]string
}

var s *Server

func dispatch(c raft.Command, w http.ResponseWriter, req *http.Request) (err error) {
  r := s.Raft
  state := r.State()
  if state == raft.Leader {
    if body, err := s.Raft.Do(c); err != nil {
      return err
    } else {
      body, _ := body.([]byte)
      w.WriteHeader(http.StatusOK)
      w.Write(body)
      return err
    }
  } else {
    leader := s.Raft.Leader()
    log.Printf("%s\n", leader)

    // tell the client where is the leader
    path := req.URL.Path
    var url string
    addr, ok := s.Raft.Peers()[leader]
    if !ok {
      panic("No name for leader")
    }

    url = addr.ConnectionString + path

    http.Redirect(w, req, url, http.StatusTemporaryRedirect)
    return
  }
  w.WriteHeader(300)
  _, err = w.Write([]byte{})
  return
}

func EnqueueHandler(w http.ResponseWriter, r *http.Request) {
  bytes, err := ioutil.ReadAll(r.Body)
  if err != nil {
    log.Fatalf("enqueue: %+v\n", err)
  }
  cmd := &EnqueueCommand{Message: bytes}
  err = dispatch(cmd, w, r)
  if err != nil {
    log.Fatalf("enqueue dispatch: %+v\n", err)
  }
}

func DequeueHandler(w http.ResponseWriter, r *http.Request) {
  err := dispatch(&DequeueCommand{}, w, r)
  if err != nil {
    log.Fatal(err)
  }
}

func NewServer(name string, prefix string, address string) *Server {
  mux := http.NewServeMux()

  mux.HandleFunc(path.Join(prefix, "enqueue"), EnqueueHandler)
  log.Printf("Mounting %s", path.Join(prefix, "enqueue"))
  mux.HandleFunc(path.Join(prefix, "dequeue"), DequeueHandler)
  log.Printf("Mounting %s", path.Join(prefix, "dequeue"))

  return &Server{
    HTTPServer: &http.Server{
      Handler:      mux,
      Addr:         address,
      ReadTimeout:  10 * time.Second,
      WriteTimeout: 10 * time.Second,
    },
    Name:   name,
    Prefix: prefix,
    Raft:   nil,
    Queue:  NewQueue(),
  }
}

func (server *Server) ListenAndServe() error {
  return server.HTTPServer.ListenAndServe()
}
