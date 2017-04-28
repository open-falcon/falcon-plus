// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package falcon

import "errors"

var (
	ErrUnsupported = errors.New("unsupported")
	ErrExist       = errors.New("entry exists")
	ErrNoent       = errors.New("entry not exists")
	ErrParam       = errors.New("param error")
	ErrEmpty       = errors.New("empty items")
	EPERM          = errors.New("Operation not permitted")
	ENOENT         = errors.New("No such file or directory")
	ESRCH          = errors.New("No such process")
	EINTR          = errors.New("Interrupted system call")
	EIO            = errors.New("I/O error")
	ENXIO          = errors.New("No such device or address")
	E2BIG          = errors.New("Argument list too long")
	ENOEXEC        = errors.New("Exec format error")
	EBADF          = errors.New("Bad file number")
	ECHILD         = errors.New("No child processes")
	EAGAIN         = errors.New("Try again")
	ENOMEM         = errors.New("Out of memory")
	EACCES         = errors.New("Permission denied")
	EFAULT         = errors.New("Bad address")
	ENOTBLK        = errors.New("Block device required")
	EBUSY          = errors.New("Device or resource busy")
	EEXIST         = errors.New("File exists")
	EXDEV          = errors.New("Cross-device link")
	ENODEV         = errors.New("No such device")
	ENOTDIR        = errors.New("Not a directory")
	EISDIR         = errors.New("Is a directory")
	EINVAL         = errors.New("Invalid argument")
	ENFILE         = errors.New("File table overflow")
	EMFILE         = errors.New("Too many open files")
	ENOTTY         = errors.New("Not a typewriter")
	ETXTBSY        = errors.New("Text file busy")
	EFBIG          = errors.New("File too large")
	ENOSPC         = errors.New("No space left on device")
	ESPIPE         = errors.New("Illegal seek")
	EROFS          = errors.New("Read-only file system")
	EMLINK         = errors.New("Too many links")
	EPIPE          = errors.New("Broken pipe")
	EDOM           = errors.New("Math argument out of domain of func")
	ERANGE         = errors.New("Math result not representable")
	EFMT           = errors.New("Invalid format") // custom
	EALLOC         = errors.New("Allocation Failure")
)
