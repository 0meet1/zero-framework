package structs

import (
	"errors"
	"io"
	"os"
	"path"
)

var Xfexists = func(srcpath string) bool {
	_, err := os.Open(srcpath)
	return !(err != nil && os.IsNotExist(err))
}

var Xfmake = func(xpath string) error {
	xdir := path.Dir(xpath)
	if !Xfexists(xdir) {
		err := os.MkdirAll(xdir, 0777)
		if err != nil {
			return err
		}
	}

	if Xfexists(xpath) {
		err := os.Remove(xpath)
		if err != nil {
			return err
		}
	}
	return nil
}

var Xfwrite = func(srcpath string, datas []byte) error {
	err := Xfmake(srcpath)
	if err != nil {
		return err
	}

	distfile, err := os.Create(srcpath)
	if err != nil {
		return err
	}
	defer distfile.Close()

	distfile.Write(datas)
	return nil
}

var Xfread = func(srcpath string) ([]byte, error) {
	if !Xfexists(srcpath) {
		return nil, errors.New("file `" + srcpath + "` not found")
	}
	file, err := os.Open(srcpath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}

var Xfmove = func(srcpath string, distpath string) error {
	srcdatas, err := Xfread(srcpath)
	if err != nil {
		return err
	}

	err = Xfwrite(distpath, srcdatas)
	if err != nil {
		return err
	}

	return os.Remove(srcpath)
}
