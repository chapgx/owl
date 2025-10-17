# OWL <img src="/assets/logo.png" width="50" height="50">

Is simple file watcher `CLI` application and library fully writing in [GO](https://go.dev). The goal is to have an intuitive way to watch changes of assets in a file system (files, directories) and perform an action based on the change. The `CLI` is meant to provide you with the most commons our of the box use cases for this. The library is for usage when the use cases cover by the CLI are not enough or for integration with a larger project.

The `OWL` CLI uses the [Rhombifer](https://github.com/racg0092/rhombifer) framework. Another side project I have been working one, check it out if you are interested. Keep in mind the interface will change drastically in the upcoming changes to `Rhombifer`


## Quick Start Library

The interface to use the library is simple. Call the watch command and subscribe to the event stream to get notifications on changes.

Subscribe to all changes
```go 
package main

import "github.com/chapgx/owl"


func main() {
  // create a subscriber
  sub := owl.Subscribe()

  // this sets the path and the interval 
  //note: there is a minimal interval allowed anything lower will panic
  go owl.Watch("/path/to/file/or/dir", owl.MinInterval)


  for result := range sub {
    if result.Error != nil {
      continue
    }
    fmt.Println(resutl.Snap)
  }
}

```

If you wanted to subscribe to only modifications in the assets your using do the following
```go 
package main

import "github.com/chapgx/owl"


func main() {
  // create a subscriber that listens only for modifications in any asset
  sub := owl.SubscribeToOnModified()

  // this sets the path and the interval 
  //note: there is a minimal interval allowed anything lower will panic
  go owl.Watch("/path/to/file/or/dir", owl.MinInterval)


  for result := range sub {
    if result.Error != nil {
      continue
    }
    fmt.Println(resutl.Snap)
  }
}

```
