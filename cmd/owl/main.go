package main

import (
	"fmt"

	"github.com/chapgx/owl"
	_ "github.com/chapgx/owl/cmd/owl/cmd"
	"github.com/racg0092/rhombifer"
	"github.com/racg0092/rhombifer/pkg/builtin"
	"github.com/racg0092/rhombifer/pkg/models"
)

func main() {
	if e := rhombifer.Start(); e != nil {
		panic(e)
	}
}

func init() {
	cfg := rhombifer.GetConfig()
	cfg.RunHelpIfNoInput = true
	cfg.AllowFlagsInRoot = true

	help := builtin.HelpCommand(nil, nil)

	root := rhombifer.Root()
	root.Name = "OWL 󰏒 "
	root.ShortDesc = "File watcher"
	root.LongDesc = `
Is a file watcher cli tool and library.
`
	root.Run = rootrun

	root.AddSub(&help)
	root.AddFlags(src_flag)
}

func rootrun(args ...string) error {
	src, e := rhombifer.FindFlag("src")
	if e != nil {
		return e
	}

	path := src.Values[0]
	//TODO: need to do som evalidations for security

	sub := owl.SubscribeToOnModified()

	go owl.WatchWithMinInterval(path)

	for r := range sub {
		if r.Error != nil {
			fmt.Println(r.Error)
			continue
		}

		fmt.Printf("%+v\n", r.Snap)
	}

	return nil
}

var src_flag = &models.Flag{
	Name:     "src",
	Short:    "source to watch can be directory or file",
	Required: true,
}
