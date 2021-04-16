package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"

	"strings"
	"sync"
)

// Server 구조체
type Server struct {
	Addr     string
	PEMFile  string
	Listener *net.TCPListener
	Insecure bool

	ctx    context.Context
	cancel context.CancelFunc

	Peers sync.Map
}

func NewServer(addr, pem string, insecure bool) (*Server, error) {
	s := Server{}

	if !IsFile(pem) {
		return &s, fmt.Errorf("PEM파일 찾을수없음")
	}
	s.PEMFile = pem
	s.Insecure = insecure

	if len(addr) > 0 {
		s.Addr = addr
	} else {
		s.Addr = "0.0.0.0.1234"
	}

	address, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return &s, err
	}

	s.Listener, err = net.ListenTCP("tcp", address)
	if err != nil {
		return &s, err
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())
	return &s, nil
}

func (s *Server) Close() {
	s.Listener.Close()
	s.cancel()
}

// 새로운 클라이언트가 접속될때마다 핸들을 얼어줍니다.
type ServerReadHandleFunc func(string, *Client)

// 핸들함수를 등록하고 무한루핑
func (s *Server) Handle(h ServerReadHandleFunc) {
	for {
		select {
		case <-s.ctx.Done():
			return

		default:
			conn, err := s.Listener.AcceptTCP()
			if err != nil {
				return
			}

			conn.SetKeepAlive(true)
			conn.SetKeepAlivePeriod(time.Duration(15 * time.Second))

			tlsconn, err := getTLSConn(conn, s.PEMFile, s.Insecure)
			if err != nil {
				tlsconn.Close()
				conn.Close()
			}

			peer := Client{
				LocalAddr: getaddr(tlsconn.RemoteAddr().String()),
				conn:      tlsconn,
			}
			s.Peers.Store(peer.LocalAddr, &peer)
			go h(peer.LocalAddr, &peer)
		}
	}
}

/*
	포트를 제외한 주소값만 반환합니다. 예를들어
	10.10.10.10:1234 이면 10.10.10.10 를 반환
	02:42:e8:73:3d:f6:1234 이면 02:42:e8:73:3d:f6 를 반환
*/
func getaddr(addr string) string {
	ws := strings.Split(addr, ":")
	if len(ws) > 0 {
		return strings.Join(ws[:len(ws)-1], ":")
	}
	return addr
}

// 클라이언트 구조체
type Client struct {
	PrivateKeyFile string
	PublicKeyFile  string
	KeyPass        string
	ServerAddr     string
	LocalAddr      string

	mutex     sync.Mutex
	conn      *tls.Conn
	remainder []byte
}

// 새로운 클라이언트 생성
func NewClient(paddr, pem, keypass string) *Client {
	c := Client{}

	if len(paddr) > 0 {
		c.ServerAddr = paddr
	} else {
		c.ServerAddr = "127.0.0.1:1234"
	}

	c.PrivateKeyFile = pem
	c.PublicKeyFile = pem
	c.KeyPass = keypass

	return &c
}

// 클라이언트를 서버에 연결합니다.
func (c *Client) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	var err error

	if c.conn != nil {
		c.Close()
		c.conn = nil
	}

	privateKey, err := readPEM("RSA PRIVATE KEY", c.PrivateKeyFile, c.KeyPass)
	if err != nil {
		return err
	}

	publicCert, err := readPEM("CERTIFICATE", c.PublicKeyFile, "")
	if err != nil {
		return err
	}

	cert, err := tls.X509KeyPair(publicCert, privateKey)
	if err != nil {
		return err
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(privateKey)
	certPool.AppendCertsFromPEM(publicCert)

	sslConfig := tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            certPool,
		InsecureSkipVerify: false,
		CipherSuites: []uint16{
			//tls.TLS_RSA_WITH_RC4_128_SHA,
			//tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
			//tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			//tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			//tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
			//tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			//tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			//tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
			//tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
			//tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			//tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		},
	}

	c.conn, err = tls.Dial("tcp", c.ServerAddr, &sslConfig)
	if err != nil {
		return err
	}

	c.LocalAddr = getaddr(c.conn.LocalAddr().String())

	return nil
}

// 연결을 끊습니다.
func (c *Client) Close() error {
	err := c.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Read() ([][]byte, error) {
	var rs [][]byte

	if c.conn == nil {
		return rs, fmt.Errorf("conn null")
	}

	buf := make([]byte, 1000)
	sz, err := c.conn.Read(buf)
	if err != nil {
		return rs, err
	}
	if sz == 0 {
		return rs, nil
	}

	var msgs [][]byte
	msgs, c.remainder = splitMessages(buf)
	for _, msg := range msgs {
		if len(msg) > 0 {
			rs = append(rs, msg)
		}
	}

	return rs, nil
}

// 쓰기
func (c *Client) Write(data []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var msg []byte
	msg = append(msg, data...)
	msg = append(msg, 0)

	if c.conn == nil {
		return fmt.Errorf("disconnect")
	}
	if _, err := c.conn.Write(msg); err != nil {
		return err
	}
	return nil
}

// 인증서를 읽어옵니다.
func readPEM(name string, file string, password string) ([]byte, error) {
	var err error
	var b []byte

	var cert pem.Block

	b, err = ioutil.ReadFile(file)
	if err != nil {
		return []byte{}, fmt.Errorf("read %s", err.Error())
	}

	var d *pem.Block

	for {
		d, b = pem.Decode(b)

		if strings.Contains(strings.ToUpper(d.Type), strings.ToUpper(name)) {
			cert = *d
		}

		if len(b) <= 0 {
			break
		}
	}

	if len(cert.Bytes) > 0 && x509.IsEncryptedPEMBlock(&cert) {
		pdec, err := x509.DecryptPEMBlock(&cert, []byte(password))
		if err != nil {
			return []byte{}, fmt.Errorf("DecryptPEMBlock %s %s", name, err.Error())
		}

		// rsakey의 bytes를 decrypt된 것으로 바꿉니다.
		cert.Bytes = pdec
	}

	return pem.EncodeToMemory(&cert), nil
}

func splitMessages(b []byte) ([][]byte, []byte) {
	bs := bytes.Split(b, []byte{0})
	return bs[:len(bs)-1], bs[len(bs)-1]
}

func getTLSConfig(file string, insecure bool) (*tls.Config, error) {
	var sslConfig tls.Config

	certCACert, err := ioutil.ReadFile(file)
	if err != nil {
		return &sslConfig, err
	}
	clientCertPool := x509.NewCertPool()
	if ok := clientCertPool.AppendCertsFromPEM(certCACert); !ok {
		return &sslConfig, fmt.Errorf("AppendCertsFromPEM error")
	}

	cert, _ := tls.LoadX509KeyPair(file, file)
	sslConfig = tls.Config{
		ClientAuth:         tls.RequireAndVerifyClientCert,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: insecure,
		ClientCAs:          clientCertPool,
	}
	return &sslConfig, nil
}

func getTLSConn(conn *net.TCPConn, file string, insecure bool) (*tls.Conn, error) {
	var tlsconn *tls.Conn

	sslConfig, err := getTLSConfig(file, insecure)
	if err != nil {
		return tlsconn, err
	}
	tlsconn = tls.Server(conn, sslConfig)
	err = tlsconn.Handshake()
	if err != nil {
		return tlsconn, err
	}
	return tlsconn, nil
}

func IsFile(path string) bool {
	if st, err := os.Stat(path); err == nil {
		return !st.IsDir()
	}
	return false
}
