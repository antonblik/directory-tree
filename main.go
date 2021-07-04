package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Vertex interface {
	fmt.Stringer
}
type File struct {
	Name string
	Size int64
}
type Directory struct {
	Name     string
	children []Vertex
}

func (f File) String() string {
	if f.Size == 0 {
		return f.Name + "(empty)"
	} else {
		return f.Name + "(" + strconv.FormatInt(f.Size, 10) + "b)"
	}
}
func (d Directory) String() string {
	return d.Name
}

func rDir(path string, vertexes []Vertex, containFiles bool) (error, []Vertex) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	files, err := file.Readdir(0)
	if err != nil {
		panic(err)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, children := range files {
		if !(children.IsDir() || containFiles) {
			continue
		}
		var newVertex Vertex

		if children.IsDir() {
			_, Checker := rDir(filepath.Join(path, children.Name()), []Vertex{}, containFiles)
			newVertex = Directory{
				Name:     children.Name(),
				children: Checker,
			}
		} else {
			newVertex = File{
				Name: children.Name() + " ",
				Size: children.Size(),
			}
		}
		vertexes = append(vertexes, newVertex)
	}
	return err, vertexes
}
func wDir(out io.Writer, vertexes []Vertex, symbols []string) {
	if len(vertexes) == 0 {
		return
	}
	fmt.Fprintf(out, "%s", strings.Join(symbols, ""))

	vertex := vertexes[0]

	if len(vertexes) == 1 {
		fmt.Fprintf(out, "%s%s\n", "└───", vertex)
		if dir, ok := vertex.(Directory); ok {
			wDir(out, dir.children, append(symbols, "\t"))
		}
		return
	}
	fmt.Fprintf(out, "%s%s\n", "├───", vertex)
	if dir, ok := vertex.(Directory); ok {
		wDir(out, dir.children, append(symbols, "│\t"))
	}
	wDir(out, vertexes[1:], symbols)
}
func dirTree(out io.Writer, path string, containFiles bool) error {
	err, vertexes := rDir(path, []Vertex{}, containFiles)
	if err != nil {
		panic(err)
	}
	wDir(out, vertexes, []string{})
	return err
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
