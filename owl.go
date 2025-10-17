package owl

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

type Result struct {
	Snap  SnapShot
	Error error
}

var (
	high_priority_queue []string
	standard_queue      []string
)

var (
	subscribers           []chan Result
	subscribersOnModified []chan Result
)

var state_storage *State

var ticker *time.Ticker

const MinInterval = time.Millisecond * 500

var (
	stop   = make(chan os.Signal, 1)
	output = make(chan Result, 1)
)

// WatchWithMinInterval starts the watcher with minimum interval allowed
func WatchWithMinInterval(path string) {
	Watch(path, MinInterval)
}

// Watch starts the watcher with the specified path and interval. Will panic if interval is less than 500 milliseconds
func Watch(path string, interval time.Duration) {
	info, e := os.Stat(path)
	if e != nil {
		panic(e)
	}

	if info.IsDir() {
		recorddir(path)
	} else {
		standard_queue = append(standard_queue, path)
	}

	if interval < MinInterval {
		panic(fmt.Sprintf("minimum interval allow is %d but got %d", MinInterval, interval))
	}

	ticker = time.NewTicker(interval)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	state_storage = &State{store: make(map[string]SnapShot)}

	for {
		select {
		case <-ticker.C:
			processQueues()
			go signalSubscribers()
		case <-stop:
			fmt.Println("exiting from os signal Interrupt")
			os.Exit(1)
			return
		}
	}
}

// recorddir record paths in the standard priority queue
func recorddir(path string) {
	entries, e := os.ReadDir(path)
	if e != nil {
		// output <- e
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			go recorddir(filepath.Join(path, entry.Name()))
			continue
		}
		standard_queue = append(standard_queue, filepath.Join(path, entry.Name()))
	}
}

func processQueues() {
	go processQueue(high_priority_queue)
	go processQueue(standard_queue)
}

func processQueue(queue []string) {
	for _, p := range queue {
		snap, e := takesnap(p)
		if e != nil {
			output <- Result{Error: e}
			continue
		}
		output <- Result{Snap: snap}
	}
}

func Subscribe() chan Result {
	sub := make(chan Result, 1)
	subscribers = append(subscribers, sub)
	return sub
}

func SubscribeToOnModified() chan Result {
	sub := make(chan Result, 1)
	subscribersOnModified = append(subscribersOnModified, sub)
	return sub
}

func signalSubscribers() {
	for r := range output {
		go signal_on_any_subs(r)
		go signal_on_change_subs(r)
	}
}

func signal_on_any_subs(r Result) {
	for _, sub := range subscribers {
		sub <- r
	}
}

func signal_on_change_subs(r Result) {
	for _, sub := range subscribersOnModified {
		if r.Error != nil {
			continue
		}

		prev := state_storage.get(r.Snap.Path)
		if prev == nil {
			state_storage.set(r.Snap)
			continue
		}

		if prev.ModTime.Before(r.Snap.ModTime) {
			state_storage.set(r.Snap)
			sub <- r
		}
	}
}
