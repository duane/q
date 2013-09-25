package main

import (
  "flag"
  "fmt"
  "github.com/coreos/raft"
  "log"
  "math/rand"
  "net/http"
  "strings"
  "time"
)

type Info struct {
  ClientAddr string
  ServerAddr string
}

type Addr struct {
  Port int
  Host string
}

type Cluster []Addr

func (cluster *Cluster) String() string {
  machines := []string{}
  for _, addr := range []Addr(*cluster) {
    machines = append(machines, fmt.Sprintf("%s:%d", addr.Host, addr.Port))
  }
  return strings.Join(machines, ",")
}

func (cluster *Cluster) Set(value string) error {
  machines := strings.Split(value, ",")
  for _, machine := range machines {
    elements := strings.SplitN(machine, ":", 2)
    addr := Addr{Host: elements[0]}
    if len(elements) != 2 {
      return fmt.Errorf("address \"%s\" could not be parsed, use [hostname]:[port]")
    }

    _, err := fmt.Sscanf(elements[1], "%d", &addr.Port)
    if err != nil {
      return err
    }
    *cluster = Cluster(append([]Addr(*cluster), addr))
  }
  return nil
}

func main() {
  raft.SetLogLevel(raft.Debug)
  raft.RegisterCommand(&EnqueueCommand{})
  raft.RegisterCommand(&DequeueCommand{})

  // random number between 19000-19999
  lower := 19000
  upper := 19999
  delta := upper - lower
  rand.Seed(time.Now().Unix())

  client_port := int(rand.Uint32())%delta + lower
  server_port := int(rand.Uint32())%delta + lower

  client_addr := fmt.Sprintf("localhost:%d", client_port)
  server_addr := fmt.Sprintf("localhost:%d", server_port)

  cluster := Cluster{}

  flag.StringVar(&client_addr, "c", client_addr, "the hostname:port of public client interface")
  flag.StringVar(&server_addr, "s", server_addr, "the hostname:port of public server interface")
  flag.Var(&cluster, "cluster", "hostname:port of other machines in the cluster, separated by a comma")

  flag.Parse()
  prefix := "/q"

  name := server_addr

  raft_muxer := http.NewServeMux()
  raft_server := &http.Server{
    Handler:      raft_muxer,
    Addr:         server_addr,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
  }

  transporter := raft.NewHTTPTransporter("")

  s = NewServer(name, prefix, client_addr)
  var err error
  s.Raft, err = raft.NewServer(name, "./tmp", transporter, s.Queue, transporter, server_addr)

  if err != nil {
    log.Fatal(err)
  }

  transporter.Install(s.Raft, raft_muxer)

  s.Raft.Start()

  if len(cluster) == 0 {
    response, err := s.Raft.Do(&raft.DefaultJoinCommand{Name: s.Name, ConnectionString: server_addr})
    if err != nil {
      log.Fatalf("Attempt at self join: %+v", err)
    }
    log.Printf("Join response: %+v\n", response)
  } else {
    for _, addr := range []Addr(cluster) {
      loc := fmt.Sprintf("%s:%d", addr.Host, addr.Port)
      log.Println(loc)

      response, err := s.Raft.Do(&raft.DefaultJoinCommand{Name: s.Name, ConnectionString: loc})
      if err != nil {
        log.Fatalf("Add peer failed: %+v\n", err)
      }
      log.Printf("Join response: %+v\n", response)
    }
  }

  log.Printf("Serving client at %s, raft at %s, name is %s.\n", client_addr, server_addr, name)
  log.Printf("%+v\n", s.Raft)

  go func() { log.Fatalf("raft: %+v\n", raft_server.ListenAndServe()) }()
  log.Fatalf("q: %+v\n", s.ListenAndServe())
}
