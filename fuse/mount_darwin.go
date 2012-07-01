// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TODO: Rewrite using package syscall not cgo

package fuse

/*

// Adapted from Plan 9 from User Space's src/cmd/9pfuse/fuse.c,
// which carries this notice:
//
// The files in this directory are subject to the following license.
// 
// The author of this software is Russ Cox.
// 
//         Copyright (c) 2006 Russ Cox
// 
// Permission to use, copy, modify, and distribute this software for any
// purpose without fee is hereby granted, provided that this entire notice
// is included in all copies of any software which is or includes a copy
// or modification of this software and in all copies of the supporting
// documentation for such software.
// 
// THIS SOFTWARE IS BEING PROVIDED "AS IS", WITHOUT ANY EXPRESS OR IMPLIED
// WARRANTY.  IN PARTICULAR, THE AUTHOR MAKES NO REPRESENTATION OR WARRANTY
// OF ANY KIND CONCERNING THE MERCHANTABILITY OF THIS SOFTWARE OR ITS
// FITNESS FOR ANY PARTICULAR PURPOSE.

#include <stdlib.h>
#include <sys/param.h>
#include <sys/mount.h>
#include <unistd.h>
#include <string.h>
#include <stdio.h>
#include <errno.h>
#include <fcntl.h>

#define nil ((void*)0)

static int
load_osxfusefs(char *mtpt, char **err)
{
	int i, fd, r;
	char buf[200];
	struct vfsconf vfs;
	char *f;

	if(getvfsbyname("fusefs", &vfs) < 0){
		if(access(f="/Library/Filesystems/osxfusefs.fs"
			"/Support/load_osxfusefs", 0) < 0){
		         *err = strdup("cannot find load_fusefs");
		   	return -1;
		}
		if((r=system(f)) < 0){
			snprintf(buf, sizeof buf, "%s: %s", f, strerror(errno));
			*err = strdup(buf);
			return -1;
		}
		if(r != 0){
			snprintf(buf, sizeof buf, "load_fusefs failed: exit %d", r);
			*err = strdup(buf);
			return -1;
		}
		if(getvfsbyname("osxfusefs", &vfs) < 0){
			snprintf(buf, sizeof buf, "getvfsbyname osxfusefs: %s", strerror(errno));
			*err = strdup(buf);
			return -1;
		}
	}

	// Look for available FUSE device.
	for(i=0;; i++){
		snprintf(buf, sizeof buf, "/dev/osxfuse%d", i);
		if(access(buf, 0) < 0){
			*err = strdup("no available fuse devices");
			return -1;
		}
		if((fd = open(buf, O_RDWR)) >= 0)
			break;
	}
	return fd;
}
*/
import "C"

import "unsafe"

import (
	"syscall"
	"strconv"
	"os/exec"
	"fmt"
)

func mount_osxfusefs(fd int, mnt string) error {
	var err error
	if err = syscall.Setenv("MOUNT_FUSEFS_CALL_BY_LIB", ""); err != nil {
		return err
	}
	if err = syscall.Setenv("MOUNT_FUSEFS_DAEMON_PATH", "/Library/Filesystems/osxfusefs.fs/Support/mount_osxfusefs"); err != nil {
		return err
	}
	cmd := exec.Command("/Library/Filesystems/osxfusefs.fs/Support/mount_osxfusefs", "-o", "iosize=4096", strconv.Itoa(fd), mnt)
	/*cmd.SysProcAttr = &syscall.SysProcAttr{
		//Setsid:  true,
		Setpgid: true,
	}*/
	out, err := cmd.CombinedOutput()
	if len(out) > 0 || err != nil {
		fmt.Errorf("mount_osxfusefs: %q, %v", out, err)
	}
	return nil
}

func mount(dir string) (int, string) {
	errp := (**C.char)(C.malloc(16))
	*errp = nil
	defer C.free(unsafe.Pointer(errp))
	cdir := C.CString(dir)
	defer C.free(unsafe.Pointer(cdir))
	fd := C.load_osxfusefs(cdir, errp)
	var err string
	if *errp != nil {
		err = C.GoString(*errp)
		return int(fd), err
	}
	go mount_osxfusefs(int(fd), dir)
	//if err0 := fork_mount_osxfusefs(int(fd), dir); err0 != nil {
	//	return int(fd), err0.Error()
	//}
	return int(fd), ""
}

/*
func mount(dir string) (int, string) {
	errp := (**C.char)(C.malloc(16))
	*errp = nil
	defer C.free(unsafe.Pointer(errp))
	cdir := C.CString(dir)
	defer C.free(unsafe.Pointer(cdir))
	fd := C.mountfuse(cdir, errp)
	var err string
	if *errp != nil {
		err = C.GoString(*errp)
	}
	return int(fd), err
}
*/
