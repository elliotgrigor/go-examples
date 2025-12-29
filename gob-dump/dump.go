package main

import (
	"bytes"
	"encoding/gob"
	"os"
)

func GobDump(path string, a any, perm os.FileMode) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(a); err != nil {
		return err
	}
	err := os.WriteFile(path, buf.Bytes(), perm)
	if err != nil {
		return err
	}
	return nil
}

func GobDumpAtomic(path string, a any, perm os.FileMode) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(a); err != nil {
		return err
	}

	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if _, err = f.Write(buf.Bytes()); err != nil {
		f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err := f.Chmod(perm); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
