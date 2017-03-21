package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/blackchip-org/vt128/ansi"
	"github.com/blackchip-org/vt128/d71"
)

const (
	prog = "d71"
)

var (
	disk     string
	commands = map[string]func([]string){
		"create": create,
		"dir":    dir,
	}
)

func init() {
	flag.StringVar(&disk, "d", "disk.d71", "disk image to use")
}

func w(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func usage() {
	w("\nusage: %v [options] command <args...>\n", prog)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		w("no command provided\n")
		usage()
		os.Exit(1)
	}
	cmdName := flag.Arg(0)
	args := flag.Args()[1:]

	cmd, ok := commands[cmdName]
	if !ok {
		w("no such command: %v\n", cmdName)
		usage()
		os.Exit(1)
	}
	cmd(args)
}

func create(args []string) {
	var (
		force bool
		name  string
		id    string
	)

	fs := flag.NewFlagSet("create", flag.ExitOnError)
	fs.BoolVar(&force, "f", false, "create disk if file already exists")
	fs.StringVar(&name, "n", "", "name of the disk")
	fs.StringVar(&id, "i", "", "disk id")
	fs.Parse(args)

	_, err := os.Stat(disk)
	if err == nil && !force {
		w("disk file already exists: %v\n", disk)
		os.Exit(1)
	}

	d, err := d71.NewDisk(name, id)
	if err != nil {
		w("unable to create image: %v\n", err)
		os.Exit(1)
	}
	err = d.Save(disk)
	if err != nil {
		w("unable to save image: %v\n", err)
		os.Exit(1)
	}
}

func dir(args []string) {
	d, err := d71.Load(disk)
	if err != nil {
		w("unable to load disk: %v\n", err)
	}
	info := d.Info()
	fmt.Printf("0 %v\"%-16v\" %2v %2v%v\n", ansi.Reverse, info.Name, info.ID,
		info.DosType, ansi.Normal)
	list := d.List()
	for _, file := range list {
		fmt.Printf("%-4d %-18v %v\n", file.Size, "\""+file.Name+"\"", file.Type)
	}
	fmt.Printf("%-4d BLOCKS FREE\n", info.Free)
}
