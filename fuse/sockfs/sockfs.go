// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Hellofs implements a simple "hello world" file system.
package main

import (
	"fmt"
	"log"
	"os"
	"github.com/petar/rscfuse/fuse"
)

// TODO:
//	* When program that created Unix socket is closed, socket file does not disappear
//		* Furthermore manual rm on the file gives permission denied
//      * Is it correct to get file Inode from MknodRequest.Header.Node?

func main() {
	println("Mounting ...")
	c, err := fuse.Mount("/mnt/sockfs")
	if err != nil {
		log.Fatal(err)
	}
	println("Mounted.")

	c.Serve(FS{})
}

// FS implements the hello world file system.
type FS struct{}

func (FS) Root() (fuse.Node, fuse.Error) {
	return rootDir, nil
}

var rootDir = makeDir()

type Dir struct{
	mknods map[string]*File
}

func makeDir() *Dir {
	return &Dir{
		mknods: make(map[string]*File),
	}
}

func (Dir) Attr() fuse.Attr {
	return fuse.Attr{Mode: os.ModeDir | 0555}
}

func (d *Dir) Lookup(name string, intr fuse.Intr) (fuse.Node, fuse.Error) {
	nod, ok := d.mknods[name]
	if !ok {
		return nil, fuse.ENOENT
	}
	return nod, nil
}

func (d *Dir) Mknod(req *fuse.MknodRequest, intr fuse.Intr) (fuse.Node, fuse.Error) {
	println("dir.Mknod")
	file := MakeFile(req)
	d.mknods[req.Name] = file
	return file, nil
}

func (d *Dir) Forget() {
	println("dir.Forget")
}

func (d *Dir) Remove(req *fuse.RemoveRequest, intr fuse.Intr) fuse.Error {
	println("dir.Remove")
	delete(d.mknods, req.Name)
	return nil
}

func (d *Dir) ReadDir(intr fuse.Intr) ([]fuse.Dirent, fuse.Error) {
	dirs := make([]fuse.Dirent, 0, len(d.mknods))
	for name, file := range d.mknods {
		dirs = append(dirs, fuse.Dirent{
			Inode: uint64(file.req.Header.Node),
			Name:  name, 
			Type:  12,  // DT_SOCK = 12, see `man dirent`
		})
	}
	return dirs, nil
}

// File implements both Node and Handle for the hello file.
type File struct {
	req   *fuse.MknodRequest
}

func MakeFile(req *fuse.MknodRequest) *File {
	fmt.Printf("%#v\n", req)
	return &File{
		req: req,
	}
}

func (f *File) Attr() fuse.Attr {
	println("file.Attr")
	return fuse.Attr{
		Inode: uint64(f.req.Header.Node),
		Mode:  f.req.Mode,
		Uid:   f.req.Header.Uid,
		Gid:   f.req.Header.Gid,
		Rdev:  f.req.Rdev,
	}
}

func (f *File) Open(req *fuse.OpenRequest, resp *fuse.OpenResponse) (fuse.Handle, fuse.Error) {
	println("file.Open")
	return f, nil
}

func (f *File) Forget() {
	println("file.Forget")
}

func (f *File) Release(req *fuse.ReleaseRequest, intr fuse.Intr) fuse.Error {
	println("file.Release")
	return nil
}
