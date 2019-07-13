package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

type symbolicLink struct {
	name       string
	path       string
	actualPath string
}

const (
	listCommand   string = "list"
	createCommand string = "create"
	unlinkCommand string = "unlink"
)

var (
	isRelative bool
	currentDir string

	// color
	yellow = color.New(color.FgYellow).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
)

func init() {
	var err error
	currentDir, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "Symlin"
	app.Version = "1.0.0"
	app.Author = "kyklades"
	app.Usage = "Symbolic links manager."
	app.Description = "This manages symbolic links."
	app.Commands = getCommands()
	sort.Sort(cli.CommandsByName(app.Commands))
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getCommands() []cli.Command {
	return []cli.Command{
		{
			Name:        listCommand,
			Aliases:     []string{"l"},
			Usage:       "List up symbolic links.",
			ArgsUsage:   "[target dir]",
			Description: "List up symbolic links in a target directory. Default target directory is the current directory.",
			Action:      doListUp,
		},
		{
			Name:        createCommand,
			Aliases:     []string{"c"},
			Usage:       "Create a symbolic link.",
			ArgsUsage:   "[target path] [new path]",
			Description: "Create a symbolic link.",
			Action:      doCreate,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "relative, r",
					Usage:       "if true, the type of symbolic link is relative symbolic link.",
					Destination: &isRelative,
				},
			},
		},
		{
			Name:        unlinkCommand,
			Aliases:     []string{"u"},
			Usage:       "Unlink a target symbolic link.",
			ArgsUsage:   "[target symbolic link]",
			Description: "Unlink a target symbolic link.",
			Action:      doUnlink,
		},
	}
}

func doListUp(c *cli.Context) (err error) {
	var symbolicLinks []symbolicLink
	switch len(c.Args()) {
	case 0:
		symbolicLinks, err = listUp(currentDir)
	case 1:
		var targetDirPath string
		targetDirPath, err = formatPath(c.Args().Get(0))
		if err != nil {
			return nil
		}
		symbolicLinks, err = listUp(targetDirPath)
	default:
		return
	}

	if symbolicLinks != nil {
		for _, symbolicLink := range symbolicLinks {
			fmt.Printf("%s -> %s\n", yellow(symbolicLink.name), cyan(symbolicLink.actualPath))
		}
	}
	return err
}

func listUp(dir string) ([]symbolicLink, error) {
	var (
		list []symbolicLink
		err  error
	)

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if path != dir && info.IsDir() {
			return filepath.SkipDir
		} else if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			actualPath, e := os.Readlink(path)
			if e != nil {
				return e
			}
			list = append(list, symbolicLink{filepath.Base(path), path, actualPath})
		}
		return nil
	})

	return list, err
}

func doCreate(c *cli.Context) (err error) {
	if len(c.Args()) != 2 {
		cli.ShowCommandHelp(c, createCommand)
		return nil
	}

	var paths []string
	paths, err = formatPaths(c.Args())
	if err != nil {
		return err
	}

	err = os.Symlink(paths[0], paths[1])
	if err == nil {
		fmt.Printf("New symbolic link has been created!\n%s -> %s\n", yellow(filepath.Base(paths[1])), cyan(paths[0]))
	}
	return err
}

func doUnlink(c *cli.Context) (err error) {
	if len(c.Args()) != 1 {
		cli.ShowCommandHelp(c, unlinkCommand)
		return nil
	}

	var target string
	target, err = formatPath(c.Args().Get(0))
	err = unlink(target)
	if err == nil {
		fmt.Printf("%s has been successfully unlinked!\n", yellow(target))
	}
	return err
}

func unlink(target string) error {
	_, err := os.Lstat(target)
	if err != nil {
		return err
	}
	return os.Remove(target)
}

func formatPath(path string) (formattedPath string, err error) {
	formattedPath = filepath.Clean(path)
	if isRelative {
		return path, nil
	}
	return filepath.Abs(path)
}

func formatPaths(paths []string) (formatPaths []string, err error) {
	for _, v := range paths {
		path, err := formatPath(v)
		if err != nil {
			return nil, err
		}
		formatPaths = append(formatPaths, path)
	}
	return formatPaths, nil
}
