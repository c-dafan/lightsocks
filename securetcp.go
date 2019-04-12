package lightsocks

import (
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
	log.Println(time.Now(), "目标地址：", dst.Address, "源地址：", secureSocket.Address)
	return secureSocket.EncodeCopy(dst)
}

// 从src中源源不断的读取原数据加密后写入到dst，直到src中没有数据可以再读取
func (secureSocket *SecureTCPConn) EncodeCopy(dst io.ReadWriteCloser) error {
	buf := make([]byte, bufSize)
	for {
		readCount, errRead := secureSocket.Read(buf)
		if errRead != nil {
			if errRead != io.EOF {
				return errRead
			} else {
				return nil
			}
		}
		if readCount > 0 {
			writeCount, errWrite := (&SecureTCPConn{
				ReadWriteCloser: dst,
				Cipher:          secureSocket.Cipher,
			}).EncodeWrite(buf[0:readCount])
			if errWrite != nil {
				return errWrite
			}
			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
	}
}

// 从src中源源不断的读取加密后的数据解密后写入到dst，直到src中没有数据可以再读取
func (secureSocket *SecureTCPConn) DecodeCopy(dst io.Writer) error {
	buf := make([]byte, bufSize)
	for {
		readCount, errRead := secureSocket.DecodeRead(buf)
		if errRead != nil {
			if errRead != io.EOF {
				return errRead
			} else {
				return nil
			}
		}
		if readCount > 0 {
			writeCount, errWrite := dst.Write(buf[0:readCount])
			if errWrite != nil {
				return errWrite
			}
			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
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
