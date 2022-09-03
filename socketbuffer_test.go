package socketbuffer_test

import (
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/danlapid/socketbuffer"
)

func TestGetReadBufferUdp(t *testing.T) {
	ip := "127.0.0.1"
	port := rand.Intn(30000) + 30000
	addr := net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	rawconn, err := conn.SyscallConn()
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		bufsize int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"test1", args{8 * 1024}, 8 * 1024, false},
		{"test2", args{100 * 1024}, 100 * 1024, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := conn.SetReadBuffer(tt.args.bufsize)
			if err != nil {
				t.Error(err)
				return
			}
			time.Sleep(300 * time.Millisecond)

			got, err := socketbuffer.GetReadBuffer(rawconn)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetReadBuffer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetReadBuffer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAvailableBytesUdp(t *testing.T) {
	ip := "127.0.0.1"
	port := rand.Intn(30000) + 30000
	addr := net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	}

	receiving_conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		t.Fatal(err)
	}
	defer receiving_conn.Close()

	sending_conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		t.Errorf("Error creating udp socket: %v", err)
		return
	}
	defer sending_conn.Close()

	rawconn, err := receiving_conn.SyscallConn()
	if err != nil {
		t.Fatal(err)
	}

	chunksize := 8192
	chunk := make([]byte, chunksize)
	for i := 0; i < 5; i++ {
		expected := (i + 1) * chunksize
		_, err := sending_conn.Write(chunk)
		if err != nil {
			t.Error(err)
			return
		}
		time.Sleep(300 * time.Millisecond)

		avail, err := socketbuffer.GetAvailableBytes(rawconn)
		if err != nil {
			t.Errorf("GetAvailableBytes() error = %v", err)
			return
		}
		if avail < expected {
			t.Errorf("GetAvailableBytes() = %v, want %v", avail, expected)
			return
		}
	}
}

func TestGetReadBufferTcp(t *testing.T) {
	ip := "127.0.0.1"
	port := rand.Intn(30000) + 30000
	addr := net.TCPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	}

	listening_conn, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		t.Fatal(err)
	}
	defer listening_conn.Close()
	var receiving_conn *net.TCPConn
	var accept_err error
	accepted_chan := make(chan int)
	go func() {
		receiving_conn, accept_err = listening_conn.AcceptTCP()
		accepted_chan <- 1
	}()

	sending_conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		t.Errorf("Error creating udp socket: %v", err)
		return
	}
	defer sending_conn.Close()

	<-accepted_chan
	if accept_err != nil {
		t.Errorf("Error creating udp socket: %v", accept_err)
		return
	}

	rawconn, err := receiving_conn.SyscallConn()
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		bufsize int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"test1", args{8 * 1024}, 8 * 1024, false},
		{"test2", args{100 * 1024}, 100 * 1024, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := receiving_conn.SetReadBuffer(tt.args.bufsize)
			if err != nil {
				t.Error(err)
				return
			}
			time.Sleep(300 * time.Millisecond)

			got, err := socketbuffer.GetReadBuffer(rawconn)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetReadBuffer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetReadBuffer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAvailableBytesTcp(t *testing.T) {
	ip := "127.0.0.1"
	port := rand.Intn(30000) + 30000
	addr := net.TCPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	}

	listening_conn, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		t.Fatal(err)
	}
	defer listening_conn.Close()
	var receiving_conn *net.TCPConn
	var accept_err error
	accepted_chan := make(chan int)
	go func() {
		receiving_conn, accept_err = listening_conn.AcceptTCP()
		accepted_chan <- 1
	}()

	sending_conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		t.Errorf("Error creating udp socket: %v", err)
		return
	}
	defer sending_conn.Close()

	<-accepted_chan
	if accept_err != nil {
		t.Errorf("Error creating udp socket: %v", accept_err)
		return
	}

	rawconn, err := receiving_conn.SyscallConn()
	if err != nil {
		t.Fatal(err)
	}

	chunksize := 8192
	chunk := make([]byte, chunksize)
	for i := 0; i < 5; i++ {
		expected := (i + 1) * chunksize
		_, err := sending_conn.Write(chunk)
		if err != nil {
			t.Error(err)
			return
		}
		time.Sleep(300 * time.Millisecond)

		avail, err := socketbuffer.GetAvailableBytes(rawconn)
		if err != nil {
			t.Errorf("GetAvailableBytes() error = %v", err)
			return
		}
		if avail < expected {
			t.Errorf("GetAvailableBytes() = %v, want %v", avail, expected)
			return
		}
	}
}
