package R06

import (
	"path/filepath"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var tests04 = `
#[cfg(test)]
mod shortinette_tests_rust0604 {
    use super::*;
    use libc;
    use std::ffi::CString;

    #[test]
    fn test_errno_last() {
        unsafe {
            *libc::__errno_location() = libc::EINTR;
        }
        let errno = Errno::last();
        assert_eq!(errno.0, libc::EINTR);
    }

    #[test]
    fn test_errno_make_last() {
        let errno = Errno(libc::ENOMEM);
        errno.make_last();
        let new_errno = unsafe { *libc::__errno_location() };
        assert_eq!(new_errno, libc::ENOMEM);
    }

    #[test]
    fn test_errno_description() {
        let errno = Errno(libc::EIO);
        let description = errno.description();
        assert_eq!(description, "Input/output error");

        let errno = Errno(libc::ENOSPC);
        let description = errno.description();
        assert_eq!(description, "No space left on device");
    }

    #[test]
    fn test_fd_open_success() {
        let file_name = CString::new("/tmp/test_fd_open_success").expect("CString::new failed");
        let fd = Fd::create(&file_name).expect("Failed to create file for test");
        assert!(fd.0 >= 0);
        fd.close().expect("Failed to close file");

        let fd_opened = Fd::open(&file_name).expect("Failed to open file");
        assert!(fd_opened.0 >= 0);
        fd_opened.close().expect("Failed to close file");

        unsafe {
            libc::unlink(file_name.as_ptr());
        }
    }

    #[test]
    fn test_fd_open_fail() {
        let file_name = CString::new("/tmp/nonexistent_file").expect("CString::new failed");
        let result = Fd::open(&file_name);
        assert!(result.is_err());
        if let Err(e) = result {
            println!("{}", e);
            assert_eq!(e.description(), "No such file or directory");
        }
    }

    #[test]
    fn test_fd_create_success() {
        let file_name = CString::new("/tmp/test_fd_create_success").expect("CString::new failed");
        let fd = Fd::create(&file_name).expect("Failed to create file");
        assert!(fd.0 >= 0);
        fd.close().expect("Failed to close file");

        unsafe {
            libc::unlink(file_name.as_ptr());
        }
    }

    #[test]
    fn test_fd_write_success() {
        let file_name = CString::new("/tmp/test_fd_write_success").expect("CString::new failed");
        let fd = Fd::create(&file_name).expect("Failed to create file");
        let data = b"Test data";
        let bytes_written = fd.write(data).expect("Failed to write to file");
        assert_eq!(bytes_written, data.len());

        fd.close().expect("Failed to close file");

        unsafe {
            libc::unlink(file_name.as_ptr());
        }
    }

    #[test]
    fn test_fd_write_fail() {
        let fd = Fd(-1);
        let data = b"Test data";
        let result = fd.write(data);
        assert!(result.is_err());
        if let Err(e) = result {
            println!("{}", e);
        }
    }

    #[test]
    fn test_fd_read_success() {
        let file_name = CString::new("/tmp/test_fd_read_success").expect("CString::new failed");
        let write_fd = Fd::create(&file_name).expect("Failed to create file");
        let data = b"Test data";
        write_fd.write(data).expect("Failed to write to file");

        let read_fd = Fd::open(&file_name).unwrap();
        let mut buffer = [0u8; 9];
        let bytes_read = read_fd.read(&mut buffer).expect("Failed to read from file");
        assert_eq!(bytes_read, data.len());
        assert_eq!(&buffer, data);

        read_fd.close().expect("Failed to close file");

        unsafe {
            libc::unlink(file_name.as_ptr());
        }
    }

    #[test]
    fn test_fd_read_fail() {
        let fd = Fd(-1);
        let mut buffer = [0u8; 9];
        let result = fd.read(&mut buffer);
        assert!(result.is_err());
        if let Err(e) = result {
            println!("{}", e);
        }
    }

    #[test]
    fn test_fd_close_success() {
        let file_name = CString::new("/tmp/test_fd_close_success").expect("CString::new failed");
        let fd = Fd::create(&file_name).expect("Failed to create file");
        assert!(fd.close().is_ok());

        unsafe {
            libc::unlink(file_name.as_ptr());
        }
    }

    #[test]
    fn test_file_open_success() {
        let file_name = CString::new("/tmp/test_file_open_success").expect("CString::new failed");
        let _ = Fd::create(&file_name).expect("Failed to create file");

        let file = File::open(&file_name).expect("Failed to open file");
        file.0.close().expect("Failed to close file");

        unsafe {
            libc::unlink(file_name.as_ptr());
        }
    }

    #[test]
    fn test_file_open_fail() {
        let file_name = CString::new("/tmp/nonexistent_file").expect("CString::new failed");
        let result = File::open(&file_name);
        assert!(result.is_err());
        if let Err(e) = result {
            println!("{}", e);
            assert_eq!(e.description(), "No such file or directory");
        }
    }

    #[test]
    fn test_file_create_success() {
        let file_name = CString::new("/tmp/test_file_create_success").expect("CString::new failed");
        let file = File::create(&file_name).expect("Failed to create file");
        assert!(file.0 .0 >= 0);

        file.0.close().expect("Failed to close file");

        unsafe {
            libc::unlink(file_name.as_ptr());
        }
    }

    #[test]
    fn test_file_write_and_read_success() {
        let file_name =
            CString::new("/tmp/test_file_write_and_read_success").expect("CString::new failed");
        let write_file = File::create(&file_name).expect("Failed to create file for writing");
        let data = b"Hello, File!";

        {
            let bytes_written = write_file.write(data).expect("Failed to write to file");
            assert_eq!(bytes_written, data.len());
        }
        std::thread::sleep(std::time::Duration::from_millis(10));
        let read_file = File::open(&file_name).expect("Failed to open file for reading");
        let mut buffer = [0u8; 12];
        read_file
            .read(&mut buffer)
            .expect("Failed to read from file");
        assert_eq!(&buffer, data);

        read_file.0.close().expect("Failed to close file");

        unsafe {
            libc::unlink(file_name.as_ptr());
        }
    }

    #[test]
    fn test_file_leak() {
        let file_name = CString::new("/tmp/test_file_leak").expect("CString::new failed");
        let file = File::create(&file_name).expect("Failed to create file");
        let fd = file.leak();
        assert!(fd.0 >= 0);

        unsafe {
            libc::unlink(file_name.as_ptr());
        }
    }

    #[test]
    fn test_file_drop() {
        let file_name = CString::new("/tmp/test_file_drop").expect("CString::new failed");
        {
            let file = File::create(&file_name).expect("Failed to create file");
            let _ = file.write(b"Testing drop").expect("Failed to drop file");
        }

        unsafe {
            libc::unlink(file_name.as_ptr());
        }
    }
}
`

func ex04Test(exercise *Exercise.Exercise) Exercise.Result {
	workingDirectory := filepath.Join(exercise.RepoDirectory, exercise.TurnInDirectory)
	if err := testutils.AppendStringToFile(tests04, exercise.TurnInFiles[0]); err != nil {
		logger.Exercise.Printf("could not write to %s: %v", exercise.TurnInFiles[0], err)
		return Exercise.InternalError(err.Error())
	}
	_, err := testutils.RunCommandLine(workingDirectory, "cargo", []string{"valgrind", "test"})
	if err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	return Exercise.Passed("OK")
}

func ex04() Exercise.Exercise {
	return Exercise.NewExercise("04", "studentcode", "ex04", []string{"src/lib.rs", "Cargo.toml"}, []string{"std::copy::Copy", "std::clone::Clone", "std::str::from_utf8_unchecked", "libc::__errno_location", "libc::strerror", "libc::{write, read, open, close}", "cstr::cstr", "std::cmp::{PartialEq, Eq, PartialOrd, Ord}", "std::fmt::{Debug, Display}"}, map[string]int{"unsafe": 100}, 20, ex04Test)
}