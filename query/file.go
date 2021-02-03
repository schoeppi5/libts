package query

import (
	"fmt"
	"math/rand"

	"github.com/schoeppi5/libts"
)

// File represents a file on a virtual server
type File struct {
	ChannelID int    `mapstructure:"cid"`
	Name      string `mapstructure:"name"`
	Size      int64  `mapstructure:"size"`
	Timestamp int64  `mapstructure:"datetime"`
	IsFile    bool   `mapstructure:"type"`
}

// FileTransfer represents an ongoing filetransfer to or from teamspeak
type FileTransfer struct {
	ClientFTFID int    `mapstructure:"clientftfid"`
	ServerFTFID int    `mapstructure:"serverftfid"`
	FTKey       string `mapstructure:"ftkey"`
	Port        int    `mapstructure:"port"`
	Size        int64  `mapstructure:"size"`
	Host        string
}

// FileInfo returns a File on server sid in channel cid with channelpassword cpw and name name
// This method is only supported by Telnet and SSH query
// sid - required
// cid - required
// cpw - optional
// name - required
func (a Agent) FileInfo(sid int, cid int, cpw string, name ...string) ([]File, error) {
	f := []File{}
	req := libts.Request{
		Command:  "ftgetfileinfo",
		ServerID: sid,
		Args: map[string]interface{}{
			"cid":  cid,
			"cpw":  cpw,
			"name": name,
		},
	}
	err := a.Query.Do(req, &f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// FileList lists all files on server sid in channel cid with channelpassword cpw in path path
// This function is only supported by telnet and SSH query
// sid - required
// cid - required
// cpw - optional
// path - required
func (a Agent) FileList(sid int, cid int, cpw string, path string) ([]File, error) {
	// test if path is file. If true, return file
	f, err := a.FileInfo(sid, cid, cpw, path)
	if err == nil {
		f[0].IsFile = true // against contrary belief (docs) ftgetfileinfo **does not** return the type **and can not** be used on directories (returns error 1538 invalid parameter) *sighn*
		return f, nil
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

// CreateDirectory on server sid in channel cid with channelpassword cpw and name name
// This function is only supported by telnet and SSH query
// sid - required
// cid - required
// cpw - optional
// name - required
func (a Agent) CreateDirectory(sid int, cid int, cpw string, name string) error {
	req := libts.Request{
		Command:  "ftcreatedir",
		ServerID: sid,
		Args: map[string]interface{}{
			"cid":  cid,
			"cpw":  cpw,
			"name": name,
		},
	}
	return a.Query.Do(req, nil)
}

// DeleteFile deletes one or more files on server sid in channel cid with channelpassword cpw and names name
// sid - required
// cid - required
// cpw - optional
// name - required
func (a Agent) DeleteFile(sid int, cid int, cpw string, name ...string) error {
	req := libts.Request{
		Command:  "ftdeletefile",
		ServerID: sid,
		Args: map[string]interface{}{
			"cid":  cid,
			"cpw":  cpw,
			"name": name,
		},
	}
	return a.Query.Do(req, nil)
}

// MoveFile on server sid from channel cid with channelpassword cpw to channel tcid with channelpassword tcpw and rename from name to newname
// sid - required
// cid - required
// cpw - optional
// tcid - required - can be same as cid (then the file will only be renamed)
// tcpw - optional - see above
// name - required
// newname - required - can be same as name (then the file will only be moved)
func (a Agent) MoveFile(sid int, cid int, cpw string, tcid int, tcpw string, name string, newname string) error {
	req := libts.Request{
		Command:  "ftrenamefile",
		ServerID: sid,
		Args: map[string]interface{}{
			"cid":     cid,
			"cpw":     cpw,
			"oldname": name,
			"newname": newname,
		},
	}
	if tcid != cid {
		req.Args["tcid"] = tcid
		req.Args["tcpw"] = tcpw
	}
	return a.Query.Do(req, nil)
}

// TODO: look into ftstop

// InitDownload initializes the download for file f with channelpassword cpw on server sid
// This is only supported by telnet and ssh query
// sid - required
// f - required
// cpw - optional
func (a Agent) InitDownload(sid int, f File, cpw string) (*FileTransfer, error) {
	clientftfid := rand.Int31n(1234) // random, unique id
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

// InitUpload initialized the upload on server sid into channel cid for a file with name name and size size
// This is only supported by telnet and ssh query
// sid - required
// cid - required
// cpw - optional
// name - required
// overwrite - optional - default false
// size - required
func (a Agent) InitUpload(sid, cid int, cpw string, name string, overwrite bool, size int) (*FileTransfer, error) {
	clientftfid := rand.Int31n(1234) // random, unique id
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

// TODO: ftlist

// DownloadAvatar for a specific client
// This function is only supported by telnet and SSH query
func (a Agent) DownloadAvatar(sid int, clientBase64Hash string) ([]byte, error) {
	f, err := a.FileInfo(sid, 0, "", fmt.Sprintf("/avatar_%s", clientBase64Hash))
	if err != nil {
		return nil, err
	}
	ft, err := a.InitDownload(sid, f[0], "")
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
// This function is only supported by telnet and SSH query
// iconID is normally an int in the structs. It is supposed to be an uint32, but the query screws up sometimes
// Read all about it here: https://community.teamspeak.com/t/bug-query-sends-wrong-icon-id-in-response/15054
// so just do uint32(iconID) to get the correct iconID
func (a Agent) DownloadIcon(sid int, iconID uint32) ([]byte, error) {
	if iconID == 0 {
		return []byte(""), nil
	}
	f, err := a.FileInfo(sid, 0, "", fmt.Sprintf("/icon_%d", iconID))
	if err != nil {
		return nil, err
	}
	ft, err := a.InitDownload(sid, f[0], "")
	if err != nil {
		return nil, err
	}
	icon, err := a.DownloadFile(ft)
	if err != nil {
		return nil, err
	}
	return icon, nil
}
