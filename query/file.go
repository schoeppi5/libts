package query

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"

	"github.com/schoeppi5/libts"
)

// FileInfo returns a File for the specified parameters
// This method is only supported by ServerQuery
func (a Agent) FileInfo(sid int, cid int, cpw string, name string) (*File, error) {
	if reflect.TypeOf(a.Query).String() != "serverquery.ServerQuery" {
		return nil, errors.New("FileInfo is not supported by this query")
	}
	f := File{}
	err := a.Query.Do(
		libts.Request{
			ServerID: sid,
			Command:  "ftgetfileinfo",
			Args: map[string]interface{}{
				"cid":  cid,
				"cpw":  cpw,
				"name": name,
			},
		}, &f)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// FileList lists all files in the given path
// This function is mainly used for debug purposes, but can easily be used, if file transfer should ever be in the scope of the api
// This function is only supported by ServerQuery
func (a Agent) FileList(sid int, cid int, cpw string, path string) ([]File, error) {
	if reflect.TypeOf(a.Query).String() != "*serverquery.ServerQuery" {
		return nil, errors.New("FileInfo is not supported by this query")
	}
	// test if path is file. If true, return file
	f, err := a.FileInfo(sid, cid, cpw, path)
	if err == nil {
		f.IsFile = true // against contrary belief (docs) ftgetfileinfo **does not** return the type **and can not** be used on directories (returns error 1538 invalid parameter) *sighn*
		return []File{*f}, nil
	}
	files := []File{}
	err = a.Query.Do(
		libts.Request{
			ServerID: sid,
			Command:  "ftgetfilelist",
			Args: map[string]interface{}{
				"cid":  cid,
				"cpw":  cpw,
				"path": path,
			},
		}, &files)
	if err != nil {
		return nil, err
	}
	return files, err
}

// InitDownload initializes the download on the TeamSpeak server
// This is only supported by ServerQuery
func (a Agent) InitDownload(sid int, f *File, cpw string) (*FileTransfer, error) {
	if reflect.TypeOf(a.Query).String() != "*serverquery.ServerQuery" {
		return nil, errors.New("FileInfo is not supported by this query")
	}
	clientftfid := rand.Int31n(1234)
	ft := FileTransfer{}
	err := a.Query.Do(
		libts.Request{
			ServerID: sid,
			Command:  "ftinitdownload",
			Args: map[string]interface{}{
				"clientftfid": clientftfid,
				"name":        f.Name,
				"cid":         f.ChannelID,
				"cpw":         cpw,
				"seekpos":     0,
			},
		}, &ft)
	if err != nil {
		return nil, err
	}
	return &ft, nil
}

// InitUpload initialized the upload on the teamspeak server
// This is only supported by ServerQuery
func (a Agent) InitUpload(sid, cid int, name string, cpw string, overwrite bool, size int) (*FileTransfer, error) {
	if reflect.TypeOf(a.Query).String() != "*serverquery.ServerQuery" {
		return nil, errors.New("FileInfo is not supported by this query")
	}
	clientftfid := rand.Int31n(1234)
	ft := FileTransfer{}
	err := a.Query.Do(
		libts.Request{
			ServerID: sid,
			Command:  "ftinitupload",
			Args: map[string]interface{}{
				"clientftfid": clientftfid,
				"name":        name,
				"cid":         cid,
				"cpw":         cpw,
				"size":        size,
				"overwrite":   overwrite,
			},
		}, &ft)
	if err != nil {
		return nil, err
	}
	return &ft, nil
}

// DownloadAvatar for a specific client
func (a Agent) DownloadAvatar(sid int, clientBase64Hash string) ([]byte, error) {
	f, err := a.FileInfo(sid, 0, "", fmt.Sprintf("/avatar_%s", clientBase64Hash))
	if err != nil {
		return nil, err
	}
	ft, err := a.InitDownload(sid, f, "")
	if err != nil {
		return nil, err
	}
	avatar, err := a.DownloadFile(ft)
	if err != nil {
		return nil, err
	}
	return avatar, nil
}

// DownloadIcon for specific iconID
// iconID is normally an int in the structs. It is supposed to be an uint32, but the query screws up sometimes
// Read all about it here: https://community.teamspeak.com/t/bug-query-sends-wrong-icon-id-in-response/15054
func (a Agent) DownloadIcon(sid int, iconID uint32) ([]byte, error) {
	if iconID == 0 {
		return []byte(""), nil
	}
	f, err := a.FileInfo(sid, 0, "", fmt.Sprintf("/icon_%d", iconID))
	if err != nil {
		return nil, err
	}
	ft, err := a.InitDownload(sid, f, "")
	if err != nil {
		return nil, err
	}
	icon, err := a.DownloadFile(ft)
	if err != nil {
		return nil, err
	}
	return icon, nil
}
