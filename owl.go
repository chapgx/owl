package owl

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

const (
	// Subscriber will recieve file content
	R_READ = 1 << 0
	// Subscriber will recieve file meta data
	R_META = 1 << 1
	// Subscriber will get notify something happend no data send
	R_SIGNAL = 1 << 2
)

// TODO: implement logic to move paths to different queues
var (
	high_priority_queue []string
	standard_queue      []string
)

var (
	subscribers           []Subscriber
	subscribersOnModified []Subscriber
)

var state_storage *State

var ticker *time.Ticker

// This is floor timing for folling file state
const MinInterval = time.Millisecond * 500

var (
	// critical error channerl that kills process
	stop = make(chan os.Signal, 1)
	// all snap shots are sen to this channel
	output = make(chan any, 1)
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

// processQueues starts processing the [standard_queue] and [high_priority_queue] queue
func processQueues() {
	go processQueue(high_priority_queue)
	go processQueue(standard_queue)
}

// processQueue process queue of paths to files
func processQueue(queue []string) {
	for _, p := range queue {
		snap, e := takesnap(p)
		if e != nil {
			output <- e
			continue
		}
		output <- snap
	}
}

// Subscribe returns a subscriber to the watcher, flag specifies what type of data you want in the channel. Meta data,
// file contents or just a signal. Example of suage below
//
//	sub := Subscribe(R_READ)
//	go Watch("path/to/fil/or/dir")
//	for rslt := range sub {
//		swicth v := rslt.(ReadSnap)
//		fmt.Println(v.Content)
//	}
func Subscribe(flag int) Subscriber {
	allowed := R_META | R_SIGNAL | R_READ
	if flag&^allowed != 0 {
		panic(fmt.Sprintf("flag %d is not allowd", flag))
	}

	sub := Subscriber{channel: make(chan any, 1), flag: flag}
	subscribers = append(subscribers, sub)
	return sub
}

// SubscribeOnModified returns a subscriber to the watcher, flag specifies what type of data you want in the channel. Meta data,
// file contents or just a signal. The diffrence from [Subscribe] is that it will only perform the action if the file has been
// modified.
//
//	sub := SubscribeOnModified(R_READ)
//	go Watch("path/to/fil/or/dir")
//	for rslt := range sub {
//		swicth v := rslt.(ReadSnap)
//		fmt.Println(v.Content)
//	}
func SubscribeOnModified(flag int) Subscriber {
	allowed := R_META | R_SIGNAL | R_READ
	if flag&^allowed != 0 {
		panic(fmt.Sprintf("flag %d is not allowd", flag))
	}
	sub := Subscriber{channel: make(chan any, 1), flag: flag}
	subscribersOnModified = append(subscribersOnModified, sub)
	return sub
}

// signalSubscribers sends a signal to all subscribers
func signalSubscribers() {
	for r := range output {
		go signal_subs(r)
		go signal_on_change_subs(r)
	}
}

// signal_subs processes standard subscribers
func signal_subs(r any) {
	for _, sub := range subscribers {
		if e, ok := r.(error); ok {
			sub.channel <- e
			continue
		}

		if R_SIGNAL&sub.flag != 0 {
			sub.channel <- uint(1)
		}

		if R_META&sub.flag != 0 {
			sub.channel <- r
		}

		if R_READ&sub.flag != 0 {
			snap, ok := r.(SnapShot)
			if !ok {
				sub.channel <- fmt.Errorf("expected a SnapShot type got something else")
				continue
			}

			b, e := os.ReadFile(snap.Path)
			if e != nil {
				sub.channel <- e
				continue
			}
			sub.channel <- ReadSnap{Path: snap.Path, ModTime: snap.ModTime, Content: b}
		}

	}
}

// signal_on_change_subs process subscribers to on modifird only
func signal_on_change_subs(r any) {
	for _, sub := range subscribersOnModified {
		if e, ok := r.(error); ok {
			sub.channel <- e
			continue
		}

		snap, ok := r.(SnapShot)
		if !ok {
			sub.channel <- fmt.Errorf("expected SnapShot got something else")
			continue
		}

		prev := state_storage.get(snap.Path)
		if prev == nil {
			state_storage.set(snap)
			continue
		}

		if !prev.ModTime.Before(snap.ModTime) {
			continue
		}

		state_storage.set(snap)

		if R_SIGNAL&sub.flag != 0 {
			sub.channel <- uint(1)
		}

		if R_META&sub.flag != 0 {
			sub.channel <- snap
		}

		if R_READ&sub.flag != 0 {
			snap, ok := r.(SnapShot)
			if !ok {
				sub.channel <- fmt.Errorf("expected a SnapShot type got something else")
				continue
			}

			b, e := os.ReadFile(snap.Path)
			if e != nil {
				sub.channel <- e
				continue
			}
			sub.channel <- ReadSnap{Path: snap.Path, ModTime: snap.ModTime, Content: b}
		}

	}
}
