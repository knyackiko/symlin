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
	listOption          bool
	dirOption           string
	unlinkOption        bool
	targetSymlinkOption string

	homeDir string

	// color
	yellow = color.New(color.FgYellow).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
)

func init() {
	user, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}
	homeDir = user.HomeDir

	flag.BoolVar(&listOption, "list", false, "List up symbolic links in a target directory.")
	flag.StringVar(&dirOption, "dir", homeDir, "Set a target directory path.")
	flag.BoolVar(&unlinkOption, "unlink", false, "Unlink a target symbolic link.")
	flag.StringVar(&targetSymlinkOption, "target", "", "Set a target symbolic link to unlink.")
	flag.Parse()
}

func main() {
	var err error
	err = formatPaths()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func execute() error {
	var err error
	switch {
	case listOption && dirOption != "":
		err = printSymbolicLinks()
	case unlinkOption && targetSymlinkOption != "":
		err = unlink()
		if err == nil {
			fmt.Printf("%s has been successfully unlinked!\n", yellow(targetSymlinkOption))
		}
	default:
		flag.PrintDefaults()
	}
	return err
}

func printSymbolicLinks() error {
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

func unlink() error {
	_, err := os.Lstat(targetSymlinkOption)
	if err != nil {
		return err
	}
	return os.Remove(targetSymlinkOption)
}

func formatPaths() error {
	var err error
	dirOption, err = formatPath(dirOption)
	targetSymlinkOption, err = formatPath(targetSymlinkOption)
	return err
}

func formatPath(path string) (formattedPath string, err error) {
	if strings.Contains(path, "~") {
		formattedPath = strings.Replace(path, "~", homeDir, 1)
	}
	formattedPath = filepath.Clean(path)
	return filepath.Abs(path)
}
