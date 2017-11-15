// Code generated by go-bindata.
// sources:
// templates/commandline.tmpl
// templates/home.tmpl
// DO NOT EDIT!

package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _templatesCommandlineTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x56\xcd\x72\xdb\x36\x10\x3e\xc7\x4f\xb1\xc3\x66\x26\xc9\x54\x10\xda\x24\xcd\x4c\x55\x4a\x87\xda\x3d\x78\xea\xc4\x33\xf9\xeb\xa5\x87\x2c\x81\x15\x89\x0a\xc4\x72\x00\x50\xb4\xea\xfa\xdd\x3b\x20\x4d\x8b\x54\x9a\xe4\x50\x5d\xc8\xdd\xc5\x7e\xfb\xf3\xed\x82\xca\xab\x58\xdb\xcd\x19\x00\x40\x5e\x11\xea\xe1\xb5\x17\x6b\x8a\x08\x0e\x6b\x5a\x67\x7b\x43\x5d\xc3\x3e\x66\xa0\xd8\x45\x72\x71\x9d\x75\x46\xc7\x6a\xad\x69\x6f\x14\x89\x5e\x58\x80\x71\x26\x1a\xb4\x22\x28\xb4\xb4\xfe\x31\x9b\x80\x59\xe3\x76\x50\x79\xda\xae\x33\x29\x6b\xbc\x51\xda\x2d\x0b\xe6\x18\xa2\xc7\x26\x09\x8a\x6b\xf9\xa0\x90\x2f\x96\x2f\x96\xaf\xa4\x0a\xe1\xa8\x5b\xd6\xc6\x2d\x55\x08\x19\x78\xb2\xeb\x2c\xc4\x83\xa5\x50\x11\xc5\x69\x98\x5e\x7b\x94\xd3\xaf\xf1\xb4\x2c\x30\x54\x70\x3b\x53\x03\x14\xa8\x76\xa5\xe7\xd6\x69\xa1\xd8\xb2\x5f\x41\x61\x51\xed\x7e\x39\x39\x76\x6f\xfb\xee\xe5\xcf\xdb\xe2\xc5\x4f\xa7\xd6\x2d\xbb\x28\x82\xf9\x9b\x56\x10\x6a\xb4\xf6\x3f\xed\x5b\xac\x8d\x3d\xac\xe0\x9c\x5d\x60\x8b\x61\xf1\x9a\x1d\x2a\x5e\x5c\xb5\xca\x68\xbc\x57\xd3\xe2\xca\x14\xe4\x31\x1a\x76\xf0\x9a\x1d\x2f\x2e\xe8\x2f\xfc\xd8\xc2\x3b\x74\x61\x50\xfc\x6a\x52\x2b\x08\x6b\xf8\x48\x1e\x27\x86\x73\x6e\xbd\x21\x0f\x6f\xa8\x5b\x40\xcd\x8e\x43\x83\x8a\x4e\x73\xe1\x3d\xf9\xad\xe5\x6e\x05\xd8\x46\x3e\xb5\x76\xec\xb5\xe8\x3c\x36\x2b\x70\xec\x6b\xfc\xac\x94\xae\x32\x91\x44\x8f\xbc\x4a\x5d\x9d\xdb\xef\x8e\x24\xc8\x09\x0b\xb9\x3c\x4e\x55\x5e\xb0\x3e\x4c\xc8\x72\xb8\x07\x65\x31\x84\x75\xe6\x70\x5f\xa0\x87\xe1\x21\x22\x97\xa5\x25\x2c\x2c\x89\x5a\x8f\x4a\xe3\xf6\xe4\x03\x41\x51\x3e\xbc\x6e\xcd\x0d\x69\x11\xb9\xc9\xa6\x94\xe7\x38\x47\x15\x85\x47\xa7\xb3\x71\xfa\xb2\x4d\x89\xae\xec\xf0\x90\x4b\x9c\xb9\x69\xf3\x90\x8e\x62\x6b\xb1\x09\x34\xc6\x7e\x90\x9b\xd6\x5a\xe1\x4d\x59\xc5\x0c\x8c\x1e\x23\xbc\x6b\x9b\xb4\x1d\xa4\xcf\x87\xed\xc8\x36\x67\x8f\x7a\xc8\xd6\x4e\x52\x19\xc1\xd2\x6b\xed\x45\xe2\x20\xdb\x9c\xf4\x38\xb7\x66\xe2\x21\x4c\xa4\xfa\xb3\x33\x27\x15\x8a\xb4\x5a\x0f\xd5\x59\x2e\xb9\x8d\xd9\xe6\xaa\x7f\x9e\x94\x38\x10\x62\xcd\x5c\x97\xcb\xd6\xce\x1a\x21\xb5\xd9\x4f\x68\x92\x0e\xa7\xe2\xac\x4d\x2e\xa2\x71\xe4\xe7\x39\xe6\xd5\xf3\xcd\x1f\x64\x15\xd7\x04\xb7\xb7\xb0\xfc\x10\xc8\xa7\x7b\x04\xee\xee\x96\xb9\xac\x9e\x0f\xdd\xc9\x9b\x11\xc6\x12\xea\x39\xc2\xa5\x03\xf6\x9a\x3c\x44\x86\x92\x22\x28\xae\x6b\x74\x3a\x95\x4a\x80\x4a\x51\x08\xc9\x14\xab\x21\xc0\xb9\x6d\x43\x24\xff\x66\x88\x01\xbf\xb7\x05\x79\x47\x91\x02\xa8\xc1\xb2\x80\x03\xb7\xd0\x19\x6b\xc1\x11\xe9\xe4\xab\xd8\x6d\x4d\xd9\x7a\x82\xeb\x86\xdc\xe5\x45\xda\x42\x47\x2a\xc2\xd3\xeb\xcb\x8b\xf3\x67\x69\x45\x2a\x72\x46\x0d\xfb\xb8\x65\x9f\x20\x3c\x28\x6b\xc8\xc5\xe5\x59\x5f\xc2\xfb\x8a\x66\xc1\xa6\x59\xb6\xd1\x58\x13\x0f\x0b\xd8\xb5\x05\xa9\x68\x17\x50\xe3\x01\x0a\x02\xe3\x42\x44\x6b\x49\x83\x35\x3b\x82\xc0\xab\xa1\x1d\xb2\x99\x37\xb1\xf1\x34\x36\x28\x5d\x5f\xd9\xe6\xec\x31\xa8\xd6\x5b\x10\x57\xd7\x50\xc5\xd8\x84\x95\x94\x21\xb2\xc7\x92\x96\x25\x73\xda\x98\xc6\x84\xfe\x22\xdd\x3d\x24\x25\x3c\x59\xc2\x40\x72\x7c\x7e\x1a\x30\xc2\xff\x80\x08\x31\xad\xe6\x32\xde\xc4\x4f\xb2\x30\x4e\x3e\x7e\xda\xf6\xf4\xfe\x03\xd8\xed\xe0\xc9\x6d\xe3\x8d\x8b\x10\xd9\x72\x47\xfe\xe9\xe3\x1f\x9e\xdd\x3d\x79\x26\xb1\xd6\xaf\x5e\xca\xfb\x66\xa4\x52\xaa\x9a\x35\x7c\x7f\x03\xcb\x89\x32\xb4\x9a\xa1\xde\x1f\x75\x20\xdb\xe0\xa5\x65\x85\xb6\x0f\x35\x1e\x9d\x8d\x6f\xe3\xe9\xa4\x75\x27\x93\xf5\xe8\xde\x70\xed\x14\x8d\x7c\x80\x09\x47\x2a\x86\xf9\x48\x04\xd1\x0d\xa9\x36\x52\x3f\x5a\x5b\xb6\x96\x3b\xe3\xca\xd5\x49\xbc\x6f\x13\x35\x06\x19\xa6\x0c\x02\x45\xa1\x3c\x69\x72\xe9\xfb\x18\xfa\xa1\xfd\xad\x46\x63\xd3\xb8\xc2\x9f\x3d\x9c\x48\x37\x42\x25\x1a\xcf\x7b\xa3\xc9\xaf\xd9\x68\xf5\x05\x9b\x40\x5f\xae\x8d\x6e\x84\x09\xa1\x25\x2f\x5a\x6f\xd7\x09\xf2\xb2\x17\x3f\xbc\xbd\xfa\x32\x6c\xef\x3a\x0c\xb1\x30\x7a\x3d\x6c\x4f\x92\x2e\x2f\xbe\xe1\xe4\x69\xeb\x29\x54\x22\xf2\x8e\x5c\xef\xf8\x76\xd0\xbc\x4f\x8a\xe4\xfc\xb5\x5c\x27\x6e\x97\x17\xa3\xc7\x57\x68\x9c\x5c\x42\xb9\x1c\xbe\x1b\xb9\x1c\xfe\xab\xfc\x1b\x00\x00\xff\xff\xb8\xf6\x8e\x4f\xb3\x08\x00\x00")

func templatesCommandlineTmplBytes() ([]byte, error) {
	return bindataRead(
		_templatesCommandlineTmpl,
		"templates/commandline.tmpl",
	)
}

func templatesCommandlineTmpl() (*asset, error) {
	bytes, err := templatesCommandlineTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/commandline.tmpl", size: 2227, mode: os.FileMode(436), modTime: time.Unix(1510775170, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesHomeTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x54\x41\x8f\xf3\x44\x0c\x3d\xf3\xfd\x0a\x2b\x5c\x40\x6a\x12\xf1\x2d\x20\x51\xd2\x4a\xd0\xbd\xac\xf8\x96\x3d\x2c\xda\xbb\x33\x71\x13\xd3\x89\x1d\xcd\x38\xc9\x16\xc4\x7f\x47\x49\xb7\x6d\x5a\xc8\x65\xec\xf7\x12\xdb\xf1\x7b\x9a\xa2\xb1\xd6\x6f\x3f\x01\x00\x14\x0d\x61\x75\x0a\xe7\xb4\x25\x43\x10\x6c\x69\x93\x0c\x4c\x63\xa7\xc1\x12\x70\x2a\x46\x62\x9b\x64\xe4\xca\x9a\x4d\x45\x03\x3b\x4a\xe7\x64\x05\x2c\x6c\x8c\x3e\x8d\x0e\x3d\x6d\xbe\x4b\x16\xc5\x3c\xcb\x01\x9a\x40\xfb\x4d\x92\xe7\x2d\xbe\xbb\x4a\xb2\x52\xd5\xa2\x05\xec\xa6\xc4\x69\x9b\x5f\x80\xfc\x21\x7b\xc8\x7e\xcc\x5d\x8c\x57\x2c\x6b\x59\x32\x17\x63\x02\x81\xfc\x26\x89\x76\xf4\x14\x1b\x22\x5b\xb6\x99\xd1\x6b\x3e\x3d\x5d\xa0\xac\xc4\xd8\xc0\xdf\x37\x30\x40\x89\xee\x50\x07\xed\xa5\x4a\x9d\x7a\x0d\x6b\x28\x3d\xba\xc3\xcf\x77\xaf\x7d\x70\x5f\x7f\xff\xd3\xbe\x7c\xf8\xe1\x9e\xdd\xab\x58\x1a\xf9\x2f\x5a\x43\x6c\xd1\xfb\xff\xe5\xf7\xd8\xb2\x3f\xae\x61\xa7\x12\xd5\x63\x5c\x3d\xab\xa0\xd3\xd5\x97\xde\x71\x85\x1f\x30\xad\xbe\x70\x49\x01\x8d\x55\xe0\x59\x45\x57\x8f\xf4\x27\xbe\xf5\xf0\x8a\x12\x4f\xc0\xaf\x3c\xad\x82\xb0\x85\x37\x0a\xb8\x20\x76\xda\x07\xa6\x00\xbf\xd3\xb8\x82\x56\x45\x63\x87\x8e\xee\x67\xd1\x81\xc2\xde\xeb\xb8\x06\xec\x4d\xef\xd9\x51\x43\x95\x8e\x01\xbb\x35\x88\x86\x16\xff\xf3\x2b\x63\xc3\x46\xe9\x5c\x79\x3d\x6d\xf5\x96\xff\xe7\x2a\x42\xbe\x50\xa1\xc8\xaf\xae\x2a\x4a\xad\x8e\x0b\xb1\x04\x07\x70\x1e\x63\xdc\x24\x82\x43\x89\x01\x4e\x47\x6a\x5a\xd7\x9e\xb0\xf4\x94\xb6\xd5\x19\x64\x19\x28\x44\x82\xb2\xbe\x84\x7b\x7e\xa7\x2a\x35\xed\x92\xa5\xe4\x05\xde\x56\x4d\xcb\x80\x52\x25\x67\xf7\x25\xdb\x1a\xa5\x1e\xf1\x58\xe4\xb8\x18\x26\x17\x1c\x16\x69\xc5\x97\xd9\x26\xcb\x23\x0b\x85\x9b\x2e\x50\x34\x9f\xb7\x3b\x95\x3d\xd7\x7d\x20\xf8\xad\x2f\x29\x08\x19\x45\xf8\xa5\xb7\x86\xc4\xd8\xcd\x52\x16\x79\xf3\xf9\xf6\xbb\xee\x5c\xd8\x13\x56\xc9\xf6\xd3\x57\x00\xf0\x47\xc3\x11\x7a\x63\xcf\x76\x84\x91\xbd\x87\x86\x7c\x07\x47\xed\x27\xad\xce\xe5\x08\x46\xb6\x66\x42\xc3\xb2\xa1\xf3\x7d\x34\x0a\x27\x12\x05\x5e\x3a\x92\xa7\xc7\xc9\x56\x42\xce\xe0\x9b\x97\xa7\xc7\xdd\xb7\x30\x29\x9f\xc1\x2b\xd7\x02\x2c\x60\x0a\x35\x19\x44\xc3\x60\x54\x65\xf3\x10\x45\xde\xdd\x4e\x7a\x59\x64\x69\x02\xa5\x49\xda\x05\x6e\x31\x1c\xe7\xd8\xd7\xf3\x51\x7a\x75\x87\xcb\x72\xbd\xd6\x2c\xc9\x76\xee\xf2\x24\x77\x1b\xae\x78\x38\x9b\xe2\xe4\x84\x22\x3f\xdd\x3e\xff\x06\x00\x00\xff\xff\xab\x73\xc1\x66\x85\x04\x00\x00")

func templatesHomeTmplBytes() ([]byte, error) {
	return bindataRead(
		_templatesHomeTmpl,
		"templates/home.tmpl",
	)
}

func templatesHomeTmpl() (*asset, error) {
	bytes, err := templatesHomeTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/home.tmpl", size: 1157, mode: os.FileMode(436), modTime: time.Unix(1510775170, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"templates/commandline.tmpl": templatesCommandlineTmpl,
	"templates/home.tmpl": templatesHomeTmpl,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"templates": &bintree{nil, map[string]*bintree{
		"commandline.tmpl": &bintree{templatesCommandlineTmpl, map[string]*bintree{}},
		"home.tmpl": &bintree{templatesHomeTmpl, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

