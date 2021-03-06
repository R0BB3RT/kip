package project

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Chart is used for describing a chart
type chart interface {
	Name() string
	Path() string
	Deploy(args []string) error
}

type Chart struct {
	name string
	path string
	prefix string
}


func (c Chart) Name() string {
	if c.Prefix() != "" {
		return strings.Join([]string{c.prefix, c.name}, "-")
	}
	return c.name
}

func (c Chart) Path() string { 
	return c.path
}

func (c Chart) Prefix() string {
	return c.prefix
}

func (c Chart) Deploy(args []string) error {
	fmt.Println(c.Path())

	cmdArgs := []string{"upgrade", c.Name(), ".", "--install"}
	cmdArgs = append(cmdArgs, args...)

	fmt.Printf("helm %s\n", strings.Join(cmdArgs, " "))

	cmd := exec.Command("helm", cmdArgs...)
	cmd.Dir = c.Path()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
		return err;
	}

	return nil
}

func getCharts(path string, prefix string) []Chart {
	charts := []Chart{}

	files, err := ioutil.ReadDir(path)
			if err != nil {
					log.Fatal(err)
			}
			

			for _, f := range files {
				if f.IsDir() {
					chartFolder := f.Name()
					var c Chart
					c = Chart{name: chartFolder, path: filepath.Join(path, chartFolder), prefix: prefix}
					charts = append(charts, c)
				}
			}

	return charts
}

func createChart(name string, path string, args []string) (string, error) {
	cmdArgs := []string{"create", name}
	cmd := exec.Command("helm", append(cmdArgs, args...)...)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return path, err
}