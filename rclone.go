package main

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/pkg/errors"
	_ "github.com/rclone/rclone/backend/all" // import all backends
	"github.com/rclone/rclone/cmd"
	"github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/config"
	"github.com/rclone/rclone/fs/operations"
)

// if you edit the below, you need to restart the gopls language server:
// https://github.com/golang/go/issues/35721#issuecomment-568543991
// https://github.com/golang/go/issues/25832#issuecomment-571631784

// #cgo darwin LDFLAGS: -Wl,-undefined,dynamic_lookup
// #include "rclone.h"
import "C"

func main() {}

//GoSetConfig sets an rclone configuration from a string
//export GoSetConfig
func GoSetConfig(cfg *C.char) {
	config.SetConfigFromString(C.GoString(cfg))
}

//GoListJSON gets a file listing
//export GoListJSON
func GoListJSON(path *C.char) {
	log.SetFlags(log.Lmicroseconds)

	// config.ShowRemotes()

	var args = []string{"gcrypt:"}
	myfs := cmd.NewFsSrc(args)

	log.Println("fs is: ", myfs)

	var opt operations.ListJSONOpt
	var items []map[string]interface{}

	err := operations.ListJSON(context.Background(), myfs, C.GoString(path), &opt, func(item *operations.ListJSONItem) error {

		data, err := json.Marshal(item)
		if err != nil {
			return errors.Wrap(err, "failed to marshal list object")
		}

		var m map[string]interface{}
		err = json.Unmarshal(data, &m)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal list object")
		}

		m["origin_parent_path"] = C.GoString(path)

		items = append(items, m)

		// newData, err := json.Marshal(m)

		// if err != nil {
		// 	return errors.Wrap(err, "failed to marshal list object")
		// }

		// log.Println()
		// _, err = os.Stdout.Write(out)
		// log.Println()

		//func C.CString(goString string) *C.char
		// log.Println(string(newData))

		// C.rust_insert_file_from_go(C.CString(string(newData)))
		// if err != nil {
		// 	return errors.Wrap(err, "failed to write to output")
		// }

		return nil
	})
	if err != nil {
		log.Println(err)
		return
	}

	data, err := json.Marshal(items)

	// log.Println(string(data))

	C.rust_insert_files_from_go(C.CString(string(data)))
	if err != nil {
		log.Println("failed to write to output")
		return
	}
}

//GoFetchFiledata fetches some file data
//export GoFetchFiledata
func GoFetchFiledata(path *C.char, startbytepos int64, endbytepos int64) {
	log.SetFlags(log.Lmicroseconds)
	log.Println("in GoFetchFiledata!")
	log.Println("path: ", C.GoString(path))
	log.Println("startbytepos: ", startbytepos)
	log.Println("endbytepos: ", endbytepos)

	var args = []string{"gcrypt:"}

	log.Println(args)

	myfs := cmd.NewFsSrc(args)

	log.Println("fs is: ", myfs)

	log.Println("calling newobject")
	o, err := myfs.NewObject(context.Background(), C.GoString(path))
	if err != nil {
		log.Println("fs.NewObject failed: ", err)
		return
	}
	log.Println("o is: ", o)

	log.Println("calling object.open")
	in, err := o.Open(context.Background(), &fs.RangeOption{Start: startbytepos, End: endbytepos})
	if err != nil {
		log.Println("o.Open failed: ", err)
		return
	}

	log.Println("calling ReadFull")

	var numbytes = (endbytepos - startbytepos) + 1
	b := make([]byte, numbytes)

	_, err = io.ReadFull(in, b[:])

	type Filecache struct {
		Cachekey     string
		Startbytepos int64
		Endbytepos   int64
	}

	item := Filecache{Cachekey: C.GoString(path), Startbytepos: startbytepos, Endbytepos: endbytepos}

	out, err := json.Marshal(item)
	if err != nil {
		log.Println(err, "failed to marshal list object")
		return
	}
	log.Println("out is: ", string(out))

	C.rust_insert_filecache_from_go(C.CString(string(out)), (*C.uchar)(C.CBytes(b[:])), C.longlong(numbytes))
}
