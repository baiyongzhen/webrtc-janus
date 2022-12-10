package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	_"os/signal"
	"time"

	"example.com/webrtc-game/pkg/config"
	"example.com/webrtc-game/pkg/util/logging"
	"example.com/webrtc-game/pkg/worker"
	"github.com/eiannone/keyboard"
	"github.com/faiface/mainthread"
	"github.com/golang/glog"
	"github.com/spf13/pflag"
)

func run() {
	rand.Seed(time.Now().UTC().UnixNano())

	cfg := config.NewDefaultConfig()
	cfg.AddFlags(pflag.CommandLine)

	logging.Init()
	defer logging.Flush()

	ctx, cancelCtx := context.WithCancel(context.Background())
	o := worker.New(ctx, cfg)
	if err := o.Run(); err != nil {
		glog.Errorf("Failed to run worker, reason %v", err)
		os.Exit(1)
	}

	/*
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	select {
	case <-stop:
		glog.Infoln("Received SIGTERM, Quiting Worker")
		o.Shutdown()
		cancelCtx()
	}
	*/

	//https://docs.libretro.com/guides/input-and-controls/
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {
		fmt.Println("Press ESC to quit")
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		//fmt.Printf("You pressed: rune %q, key %X\r\n", char, key)
		fmt.Printf("rune %d, key %d\r\n", char, key)
        if key == keyboard.KeyEsc {
			glog.Infoln("Received SIGTERM, Quiting Worker")
			o.Shutdown()
			cancelCtx()
			break
		}
		o.InputKeyboard()
	}

}


func main() {
	// enables mainthread package and runs run in a separate goroutine
	mainthread.Run(run)
}