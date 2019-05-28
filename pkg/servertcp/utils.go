package servertcp

// import (
// 	"bufio"
// 	"bytes"
// 	"io"
// 	"net"
// )

// // Read  -
// func Read(conn net.Conn) (bytes.Buffer, error) {
// 	reader := bufio.NewReader(conn)
// 	var buffer bytes.Buffer
// 	for {
// 		ba, err := reader.Read()
// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			}
// 			return buffer, err
// 		}
// 		buffer.Write(ba)
// 		if !isPrefix {
// 			break
// 		}
// 	}
// 	return buffer, nil
// }

// //Write -
// func Write(conn net.Conn, content string) (int, error) {
// 	writer := bufio.NewWriter(conn)
// 	number, err := writer.WriteString(content)
// 	if err == nil {
// 		err = writer.Flush()
// 	}
// 	return number, err
// }
