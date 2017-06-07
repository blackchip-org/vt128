package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/blackchip-org/vt128/ansi"
	"github.com/blackchip-org/vt128/d71"
)

const (
	prog = "d71"
)

type commandInfo struct {
	run  func([]string)
	help string
}

var (
	disk     string
	commands = map[string]commandInfo{
		"bam":    commandInfo{run: bam, help: "print block availability map"},
		"create": commandInfo{run: create, help: "create a formatted disk"},
		"dir":    commandInfo{run: dir, help: "list directory"},
	}
)

func init() {
	flag.StringVar(&disk, "d", "disk.d71", "disk image to use")
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %v [options] command <args...>\n", prog)
	fmt.Fprintf(os.Stderr, "\noptions:\n")
	flag.PrintDefaults()

	fmt.Fprintf(os.Stderr, "\ncommands:\n")
	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Fprintf(os.Stderr, "  %-10s  %v\n", name, commands[name].help)
	}
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "%v: error: no command provided\n\n", prog)
		usage()
		os.Exit(1)
	}
	cmdName := flag.Arg(0)
	args := flag.Args()[1:]

	cmdInfo, ok := commands[cmdName]
	if !ok {
		fmt.Fprintf(os.Stderr, "%v: error: no such command: %v\n", prog, cmdName)
		usage()
		os.Exit(1)
	}
	cmdInfo.run(args)
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
		fmt.Fprintf(os.Stderr, "%v: disk file already exists: %v\n", prog, disk)
		os.Exit(1)
	}

	d := d71.NewDisk(name, id)
	err = d.Export(disk)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: unable to save image: %v\n", prog, err)
		os.Exit(1)
	}
}

func dir(args []string) {
	d, err := d71.Import(disk)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: unable to load disk: %v\n", prog, err)
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

func bam(args []string) {
	d, err := d71.Import(disk)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to load disk: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("            1         2         3         4         5         6         7")
	fmt.Println("   1234567890123456789012345678901234567890123456789012345678901234567890")
	for sector := 0; sector < d71.MaxTrackLen; sector++ {
		fmt.Printf("%2d ", sector)
		for track := 1; track <= d71.MaxTrack; track++ {
			sectorN := d71.Geom[track].Sectors
			if sector >= sectorN {
				fmt.Print(" ")
			} else if d.BamRead(track, sector) {
				fmt.Print(".") // Free Sector
			} else {
				fmt.Print("*") // Used Sector
			}
		}
		fmt.Println()
	}
}
