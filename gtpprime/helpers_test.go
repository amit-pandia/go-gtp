package gtpprime_test

import (
	"fmt"
	"math"
	"net"
	"strconv"
	"testing"

	"github.com/amit-pandia/go-gtp/gtpprime"
)

var testConn *gtpprime.Conn
var sessions []*gtpprime.Session
var dummyAddr net.Addr = &net.UDPAddr{IP: net.IP{0x00, 0x00, 0x00, 0x00}, Port: 2123}

func init() {
	testConn = gtpprime.NewConn(dummyAddr, gtpprime.IFTypeS11MMEGTPC, 0)
	sessions = []*gtpprime.Session{
		gtpprime.NewSession(dummyAddr, &gtpprime.Subscriber{IMSI: "001011234567891"}),
		gtpprime.NewSession(dummyAddr, &gtpprime.Subscriber{IMSI: "001011234567892"}),
		gtpprime.NewSession(dummyAddr, &gtpprime.Subscriber{IMSI: "001011234567893"}),
		gtpprime.NewSession(dummyAddr, &gtpprime.Subscriber{IMSI: "001011234567894"}),
	}

	for i, sess := range sessions {
		_ = sess.Activate()
		sess.AddTEID(gtpprime.IFTypeS11MMEGTPC, uint32(i+1))
		testConn.RegisterSession(uint32(i+1), sess)
	}
}

func TestSessionCount(t *testing.T) {
	if want, got := testConn.SessionCount(), len(sessions); want != got {
		t.Errorf("SessionCount is invalid. want: %d, got: %d", want, got)
	}
}

func TestGetSessionByIMSI_GetTEID(t *testing.T) {
	for i := 1; i <= testConn.SessionCount(); i++ {
		lastDigit := strconv.Itoa(i)
		sess, err := testConn.GetSessionByIMSI("00101123456789" + lastDigit)
		if err != nil {
			t.Fatal(err)
		}

		teid, err := sess.GetTEID(gtpprime.IFTypeS11MMEGTPC)
		if err != nil {
			t.Fatal(err)
		}

		if teid != uint32(i) {
			t.Errorf("Got wrong TEID at %d, %d, %s", i, teid, sess.IMSI)
		}
	}
}

func BenchmarkAddSession(b *testing.B) {
	for k := 0.; k < 6; k++ {
		existingSessions := int(math.Pow(10, k))
		benchConn := gtpprime.NewConn(dummyAddr, gtpprime.IFTypeS11MMEGTPC, 0)
		for i := 0; i < existingSessions; i++ {
			imsi := fmt.Sprintf("%015d", i)
			benchConn.RegisterSession(0, gtpprime.NewSession(dummyAddr, &gtpprime.Subscriber{IMSI: imsi}))
		}
		b.Run(fmt.Sprintf("%d", existingSessions), func(b *testing.B) {
			for i := 1; i <= b.N; i++ {
				benchConn.RegisterSession(0, gtpprime.NewSession(dummyAddr, &gtpprime.Subscriber{IMSI: "001011234567891"}))
			}
		})
	}
}

func TestGetSessionByTEID(t *testing.T) {
	for i := 1; i <= testConn.SessionCount(); i++ {
		sess, err := testConn.GetSessionByTEID(uint32(i), dummyAddr)
		if err != nil {
			t.Fatal(err)
		}

		lastDigit := strconv.Itoa(i)
		if string(sess.IMSI[14]) != lastDigit {
			t.Errorf("Got wrong session at %d, %s", i, sess.IMSI)
		}
	}
}

func TestGetIMSIByTEID(t *testing.T) {
	for i := 1; i <= testConn.SessionCount(); i++ {
		imsi, err := testConn.GetIMSIByTEID(uint32(i), dummyAddr)
		if err != nil {
			t.Fatal(err)
		}

		lastDigit := strconv.Itoa(i)
		if string(imsi[14]) != lastDigit {
			t.Errorf("Got wrong IMSI at %d, %s", i, imsi)
		}
	}
}

func TestRemoveSession(t *testing.T) {
	testConn.RemoveSession(sessions[0])

	if testConn.SessionCount() != len(sessions)-1 {
		t.Errorf("Session not removed expectedly: %d, %v", testConn.SessionCount(), testConn.Sessions())
	}

	for i := 2; i <= testConn.SessionCount(); i++ {
		sess, err := testConn.GetSessionByTEID(uint32(i), dummyAddr)
		if err != nil {
			t.Fatal(err)
		}

		lastDigit := strconv.Itoa(i)
		if string(sess.IMSI[14]) != lastDigit {
			t.Errorf("Got wrong session at %d, %s", i, sess.IMSI)
		}
	}

	// add the session again
	s := gtpprime.NewSession(dummyAddr, &gtpprime.Subscriber{IMSI: "001011234567891"})
	_ = s.Activate()
	s.AddTEID(gtpprime.IFTypeS11MMEGTPC, uint32(0))
	testConn.RegisterSession(0, s)
}

func TestRemoveSessionByIMSI(t *testing.T) {
	testConn.RemoveSessionByIMSI("001011234567891")

	if testConn.SessionCount() != len(sessions)-1 {
		t.Errorf("Session not removed expectedly: %d, %v", testConn.SessionCount(), testConn.Sessions())
	}

	for i := 2; i <= testConn.SessionCount(); i++ {
		sess, err := testConn.GetSessionByTEID(uint32(i), dummyAddr)
		if err != nil {
			t.Fatal(err)
		}

		lastDigit := strconv.Itoa(i)
		if string(sess.IMSI[14]) != lastDigit {
			t.Errorf("Got wrong session at %d, %s", i, sess.IMSI)
		}
	}

	// add the session again
	s := gtpprime.NewSession(dummyAddr, &gtpprime.Subscriber{IMSI: "001011234567891"})
	_ = s.Activate()
	s.AddTEID(gtpprime.IFTypeS11MMEGTPC, uint32(0))
	testConn.RegisterSession(0, s)
}
