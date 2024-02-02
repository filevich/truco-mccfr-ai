package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func Append(x interface{}, file string, identation bool) {

	precheck_dir_existance(file)

	var (
		bs  []byte
		err error
	)

	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if identation {
		bs, err = json.MarshalIndent(x, "", "  ")
	} else {
		bs, err = json.Marshal(x)
	}

	if err != nil {
		panic(err)
	}

	if _, err := f.Write(bs); err != nil {
		panic(err)
	}

	if _, err := f.Write([]byte("\n")); err != nil {
		panic(err)
	}
}

func Write(x interface{}, file string, identation bool) {

	precheck_dir_existance(file)

	var (
		bs  []byte
		err error
	)

	if identation {
		bs, err = json.MarshalIndent(x, "", "  ")
	} else {
		bs, err = json.Marshal(x)
	}

	if err != nil {
		panic(err)
	}
	// ejemplo:
	// si `pwd` retorna `/home/jp/Workspace/facu/pdg/truco-cfr`
	// y corro `go run cmd/playground/writer/*.go` desde ahi
	// entonces lo va a gaurdar en `pwd`
	f, err := os.Create(file)

	if err != nil {
		panic(err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}()

	if _, err := f.Write(bs); err != nil {
		f.Close()
		panic(err)
	}
}

func Touch(file string) *os.File {

	precheck_dir_existance(file)

	// ejemplo:
	// si `pwd` retorna `/home/jp/Workspace/facu/pdg/truco-cfr`
	// y corro `go run cmd/playground/writer/*.go` desde ahi
	// entonces lo va a gaurdar en `pwd`
	f, err := os.Create(file)

	if err != nil {
		panic(err)
	}

	return f
}

func FastAppend(x interface{}, f *os.File) {
	bs, _ := json.Marshal(x)

	if _, err := f.Write(bs); err != nil {
		panic(err)
	}

	if _, err := f.Write([]byte("\n")); err != nil {
		panic(err)
	}
}

func precheck_dir_existance(path string) {
	if r := strings.Split(path, "/"); len(r) > 1 {
		// assure dir existance
		subdirs := r[:len(r)-1]
		dirs := strings.Join(subdirs, "/")
		loc := filepath.FromSlash(dirs)
		// similar to `mkdir -p dirs`
		newpath := filepath.Join(loc)
		if err := os.MkdirAll(newpath, os.ModePerm); err != nil {
			panic(err)
		}
	}
}
