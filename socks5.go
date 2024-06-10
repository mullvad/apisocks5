package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/mullvad/proxy"
)

var (
	ErrNoProxyAddressAvailable = errors.New("no proxy address available")
	ErrUnsupportedAddressType  = errors.New("unsupported address type")
	ErrUnsupportedCommand      = errors.New("unsupported command")
	ErrUnsupportedVersion      = errors.New("unsupported version")
)

const (
	Version  = 0x05
	Reserved = 0x00
)

const (
	ATypIPv4       = 0x01
	ATypDomainName = 0x03
	ATypIPv6       = 0x04
)

const (
	AuthNoAuth              = 0x00
	AuthGSSAPI              = 0x01
	AuthUsernamePassword    = 0x02
	AuthNoAcceptableMethods = 0xff
)

const (
	CommandConnect      = 0x01
	CommandBind         = 0x02
	CommandUDPAssociate = 0x03
)

const (
	StatusSucceeded               = 0x00
	StatusGeneralServerFailure    = 0x01
	StatusConnectionNotAllowed    = 0x02
	StatusNetworkUnreachable      = 0x03
	StatusHostUnreachable         = 0x04
	StatusConnectionRefused       = 0x05
	StatusTTLExpired              = 0x06
	StatusCommandNotSupported     = 0x07
	StatusAddressTypeNotSupported = 0x08
)

// handleSOCKS5Conn ensures that the given TCP connection is actually a SOCKS5
// connection and handles it accordignly if everything is OK.
func handleSOCKS5Conn(conn net.Conn, proxies []proxy.Proxy, verbose bool) {
	var err error
	var wg sync.WaitGroup
	var targetConn net.Conn
	var prx proxy.Proxy
	connBuff := bufio.NewReader(conn)

	defer conn.Close()

	if err = handleVersion(connBuff); err != nil {
		goto out
	}

	if err = handleAuthentication(conn, connBuff); err != nil {
		goto out
	}

	if err = handleHeader(connBuff); err != nil {
		goto out
	}

	if err = handleTargetAddress(connBuff); err != nil {
		goto out
	}

	targetConn, prx, err = newTargetConn(proxies, verbose)
	if err != nil {
		goto out
	}
	defer targetConn.Close()

	if err = handleTargetConn(conn, targetConn); err != nil {
		goto out
	}

	if verbose {
		log.Printf("Forwarding traffic from client to %s\n", targetConn.RemoteAddr())
	}

	wg.Add(2)
	go func() {
		prx.FromPeer(targetConn, conn)
		wg.Done()
	}()
	go func() {
		prx.ToPeer(conn, targetConn)
		wg.Done()
	}()
	wg.Wait()

out:
	if err != nil {
		log.Printf("error: %v\n", err)
	}
}

// handleVersion ensures that the first byte of the connection contains the
// SOCKS5 version identifier.
func handleVersion(connBuff *bufio.Reader) error {
	version, err := connBuff.ReadByte()
	if err != nil {
		return err
	}

	if version != Version {
		return ErrUnsupportedVersion
	}

	return nil
}

// handleAuthentication manages the authentication of the SOCKS5 connection.
func handleAuthentication(conn io.Writer, connBuff *bufio.Reader) error {
	numberOfAuthMethods, err := connBuff.ReadByte()
	if err != nil {
		return err
	}

	if _, err = connBuff.Discard(int(numberOfAuthMethods)); err != nil {
		return err
	}

	if _, err := conn.Write([]byte{Version, AuthNoAuth}); err != nil {
		return err
	}

	return nil
}

// handleHeader parses the three bytes that serves as the SOCKS5 header.
func handleHeader(connBuff *bufio.Reader) error {
	// The header contains three bytes.
	// 00: version
	// 01: command
	// 02: reserved
	header := make([]byte, 3)
	if _, err := io.ReadFull(connBuff, header); err != nil {
		return err
	}

	if header[0] != Version {
		return ErrUnsupportedVersion
	}

	if header[1] != CommandConnect {
		return ErrUnsupportedCommand
	}

	return nil
}

// handleTargetAddress parses the target address from the SOCKS5 payload. In
// our case we'll just discard this data since we will get the actual target
// from another source.
func handleTargetAddress(connBuff *bufio.Reader) error {
	typ, err := connBuff.ReadByte()
	if err != nil {
		return err
	}

	// The last two bytes that we read contains the port number and this is
	// the same for all three address types. So we can safely add the two
	// bytes here.
	bytesToDiscard := 2

	// Determine which address type that was sent and increase the number
	// of bytes to discard.
	switch typ {
	case ATypIPv4:
		bytesToDiscard += net.IPv4len
	case ATypIPv6:
		bytesToDiscard += net.IPv6len
	case ATypDomainName:
		l, err := connBuff.ReadByte()
		if err != nil {
			return err
		}
		bytesToDiscard += int(l)
	default:
		return ErrUnsupportedAddressType
	}

	if _, err = connBuff.Discard(bytesToDiscard); err != nil {
		return err
	}

	return nil
}

// newTargetConn iterates over the slice of proxies and tries to initiate a TCP
// connection to the address and port defined in the proxy.
func newTargetConn(proxies []proxy.Proxy, verbose bool) (net.Conn, proxy.Proxy, error) {
	for _, p := range proxies {
		if verbose {
			log.Printf("Using proxy with address %s and type %s\n", p.Address(), p.Type())
		}

		conn, err := net.DialTimeout("tcp", p.Address(), time.Second*10)
		if err != nil {
			if verbose {
				log.Printf("Target address %s did not respond in time\n", p.Address())
			}
			continue
		}

		return conn, p, nil
	}

	return nil, nil, ErrNoProxyAddressAvailable
}

// handleTargetConn initializes a new TCP connection to the given targetAddress
// and writes the source address back to the SOCKS5 client if everythings goes
// according to plan.
func handleTargetConn(conn io.Writer, targetConn net.Conn) error {
	localAddr := targetConn.LocalAddr().(*net.TCPAddr)
	localPort := make([]byte, 2)
	binary.BigEndian.PutUint16(localPort, uint16(localAddr.Port))

	if err := socksWrite(conn, StatusSucceeded, []byte{
		ATypIPv4,
		localAddr.IP[0], localAddr.IP[1], localAddr.IP[2], localAddr.IP[3],
		localPort[0], localPort[1],
	}); err != nil {
		return err
	}

	return nil
}

// socksWrite writes the given data to the conn.
func socksWrite(conn io.Writer, status byte, data []byte) error {
	payload := []byte{Version, status, Reserved}

	if data != nil {
		payload = append(payload, data...)
	}

	_, err := conn.Write(payload)
	if err != nil {
		return err
	}

	return nil
}
