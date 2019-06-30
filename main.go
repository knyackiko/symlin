package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

type symbolicLink struct {
	name       string
	path       string
	actualPath string
}

var (
	listOption bool
	dirOption  string

	homeDir string
)

func init() {
	user, err := user.Current()
	if err != nil {
		println(err)
	}
	homeDir = user.HomeDir

	flag.BoolVar(&listOption, "list", false, "")
	flag.StringVar(&dirOption, "dir", "", "")
	flag.Parse()
}

func main() {
	err1 := formatPath()
	if err1 != nil {
		fmt.Println(err1)
		os.Exit(1)
	}

	err2 := execute()
	if err2 != nil {
		fmt.Println(err2)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func formatPath() error {
	var err error
	if strings.Contains(dirOption, "~") {
		dirOption = strings.Replace(dirOption, "~", homeDir, 1)
	}
	dirOption = filepath.Clean(dirOption)
	dirOption, err = filepath.Abs(dirOption)
	return err
}

func execute() error {
	var err error
	switch {
	case listOption && dirOption != "":
		err = printSymbolicLinks()
	default:
		flag.PrintDefaults()
	}
	return err
}

func printSymbolicLinks() error {
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	symbolicLinks, err := listUpSymbolicLinks()
	if symbolicLinks != nil {
		for _, symbolicLink := range symbolicLinks {
			fmt.Printf("%s -> %s\n", yellow(symbolicLink.name), cyan(symbolicLink.actualPath))
		}
	}
	return err
}

func listUpSymbolicLinks() ([]symbolicLink, error) {
	var (
		list []symbolicLink
		err  error
	)

	err = filepath.Walk(dirOption, func(path string, info os.FileInfo, err error) error {
		if path != dirOption && info.IsDir() {
			return filepath.SkipDir
		} else if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			actualPath, e := os.Readlink(path)
			if e != nil {
				return e
			}
			list = append(list, symbolicLink{filepath.Base(path), path, filepath.Join(filepath.Dir(path), actualPath)})
		}
		return nil
	})

	return list, err
}
