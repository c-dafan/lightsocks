package lightsocks

import (
	l4g "github.com/alecthomas/log4go"
	"io"
	"log"
	"net"
	"time"
)

const (
	bufSize = 1024
)

type Addr interface {
	Network() string // name of the network (for example, "tcp", "udp")
	String() string  // string form of address (for example, "192.0.2.1:25", "[2001:db8::1]:80")
}

// 加密传输的 TCP Socket
type SecureTCPConn struct {
	io.ReadWriteCloser
	Address Addr
	Cipher  *cipher
}

// 从输入流里读取加密过的数据，解密后把原数据放到bs里
func (secureSocket *SecureTCPConn) DecodeRead(bs []byte) (n int, err error) {
	n, err = secureSocket.Read(bs)
	if err != nil {
		return
	}
	secureSocket.Cipher.decode(bs[:n])
	return
}

// 把放在bs里的数据加密后立即全部写入输出流
func (secureSocket *SecureTCPConn) EncodeWrite(bs []byte) (int, error) {
	secureSocket.Cipher.encode(bs)
	return secureSocket.Write(bs)
}

func (secureSocket *SecureTCPConn) EncodeCopyServer(dst *SecureTCPConn) error {
	byteSize, _, err := secureSocket.EncodeCopy(dst)
	//tagCoder := mahonia.NewDecoder("utf-8")
	//_, cdata, _ := tagCoder.Translate([]byte(content), true)
	//result := string(cdata)
	//log.Println(time.Now(), "目标地址：", dst.Address, "源地址：", secureSocket.Address, "大小:", byteSize,
	//	"内容:")
	l4g.Info("%s,%s,%s,%d", time.Now().String()[:27], dst.Address.String(),
		secureSocket.Address.String(), byteSize)
	return err
}

// 从src中源源不断的读取原数据加密后写入到dst，直到src中没有数据可以再读取
func (secureSocket *SecureTCPConn) EncodeCopy(dst io.ReadWriteCloser) (int, string, error) {
	buf := make([]byte, bufSize)
	lastCount := 0
	byteSize := 0
	str := ""
	for {
		readCount, errRead := secureSocket.Read(buf)
		if errRead != nil {
			if errRead != io.EOF {
				return 0, "", errRead
			} else {
				return (byteSize-1)*bufSize + lastCount, str, nil
			}
		}
		if readCount > 0 {

			str += string(buf[0:readCount])
			writeCount, errWrite := (&SecureTCPConn{
				ReadWriteCloser: dst,
				Cipher:          secureSocket.Cipher,
			}).EncodeWrite(buf[0:readCount])
			if errWrite != nil {
				return 0, "", errWrite
			}
			if readCount != writeCount {
				return 0, "", io.ErrShortWrite
			}
		}
		byteSize++
		lastCount = readCount
	}
}

// 从src中源源不断的读取加密后的数据解密后写入到dst，直到src中没有数据可以再读取
func (secureSocket *SecureTCPConn) DecodeCopy(dst io.Writer) (int, string, error) {
	buf := make([]byte, bufSize)
	str := ""
	byteSize := 0
	lastRead := 0
	for {
		readCount, errRead := secureSocket.DecodeRead(buf)
		str += string(buf[:readCount])
		if errRead != nil {
			if errRead != io.EOF {
				return 0, "", errRead
			} else {
				return (byteSize-1)*bufSize + lastRead, str, nil
			}
		}
		if readCount > 0 {
			writeCount, errWrite := dst.Write(buf[0:readCount])
			if errWrite != nil {
				return 0, "", errWrite
			}
			if readCount != writeCount {
				return 0, "", io.ErrShortWrite
			}
		}
		byteSize++
		lastRead = readCount
	}
}

// see net.DialTCP
func DialTCPSecure(raddr *net.TCPAddr, cipher *cipher) (*SecureTCPConn, error) {
	remoteConn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		return nil, err
	}
	return &SecureTCPConn{
		ReadWriteCloser: remoteConn,
		Cipher:          cipher,
	}, nil
}

// see net.ListenTCP
func ListenSecureTCP(laddr *net.TCPAddr, cipher *cipher, handleConn func(localConn *SecureTCPConn), didListen func(listenAddr net.Addr)) error {
	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return err
	}

	defer listener.Close()

	if didListen != nil {
		didListen(listener.Addr())
	}

	for {
		localConn, err := listener.AcceptTCP()

		//log.Println("remote: ",localConn.RemoteAddr())
		//log.Println("local: ",localConn.LocalAddr())

		if err != nil {
			log.Println(err)
			continue
		}
		// localConn被关闭时直接清除所有数据 不管没有发送的数据
		localConn.SetLinger(0)
		go handleConn(&SecureTCPConn{
			ReadWriteCloser: localConn,
			Address:         localConn.RemoteAddr(),
			Cipher:          cipher,
		})
	}
}
