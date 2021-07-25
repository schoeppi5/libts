package query

import (
	"fmt"
	"net"
)

// DownloadFile returns the File for ft
func (a Agent) DownloadFile(ft *FileTransfer) ([]byte, error) {
	// open connection to teamspeak file port
	conn, err := net.Dial("tcp", net.JoinHostPort(ft.Host, fmt.Sprint(ft.Port)))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	// send key to server
	fmt.Fprint(conn, ft.FTKey)
	buffer := make([]byte, ft.Size)
	_, err = conn.Read(buffer)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

// UploadFile uploads the given file to the server
func (a Agent) UploadFile(ft *FileTransfer, file []byte) error {
	// open connection to teamspeak file port
	conn, err := net.Dial("tcp", net.JoinHostPort(ft.Host, fmt.Sprint(ft.Port)))
	if err != nil {
		return err
	}
	defer conn.Close()
	// send key to server
	fmt.Fprint(conn, ft.FTKey)
	size, err := conn.Write(file)
	if err != nil {
		return err
	}
	if size != len(file) {
		return fmt.Errorf("failed to upload file correctly! Uploaded %d bytes of %d. %d%%", size, len(file), (size/len(file))*100)
	}
	return nil
}
