// package main

// import (
// 	//"golang.org/x/sync/errgroup"

// 	"goduckduckgo/pkg/config"

// 	"errors"
// 	"fmt"
// 	"log"
// 	"os"
// 	"os/signal"
// 	"syscall"

// 	"github.com/oklog/run"
// )

// func main() {

// 	cfg := config.Init()

// 	var g run.Group

// 	runQuery(&g, cfg)
// 	//runStore(&g)

// 	// Listen for termination signals.
// 	{
// 		cancel := make(chan struct{})
// 		g.Add(func() error {
// 			return interrupt(cancel)
// 		}, func(error) {
// 			close(cancel)
// 		})
// 	}

// 	if err := g.Run(); err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// }

// func interrupt(cancel <-chan struct{}) error {
// 	c := make(chan os.Signal, 1)
// 	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
// 	select {
// 	case s := <-c:
// 		log.Printf("Caught signal %v", s)
// 		return nil
// 	case <-cancel:
// 		return errors.New("canceled")
// 	}
// }

package main

func main() {
	execute()
}
