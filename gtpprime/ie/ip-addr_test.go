package ie_test

import (
	ie2 "github.com/amit-pandia/go-gtp/gtpprime/ie"
	"io"
	"net"
	"testing"

	"github.com/amit-pandia/go-gtp/gtpprime"
	"github.com/amit-pandia/go-gtp/gtpprime/ie"
	"github.com/google/go-cmp/cmp"
)

func TestPDNAddressAllocationIP(t *testing.T) {
	cases := []struct {
		description string
		paa         *ie2.IE
		pdnType     uint8
		ipv4        net.IP
		ipv6        net.IP
	}{
		{
			"PDNType IPv4",
			ie.NewPDNAddressAllocation("1.2.3.4"),
			gtpprime.PDNTypeIPv4,
			net.ParseIP("1.2.3.4"),
			nil,
		},
		{
			"PDNType IPv6",
			ie.NewPDNAddressAllocation("::1"),
			gtpprime.PDNTypeIPv6,
			nil,
			net.ParseIP("::1"),
		},
		{
			"PDNType IPv4v6",
			ie.NewPDNAddressAllocationDual("1.2.3.4", "::1", 64),
			gtpprime.PDNTypeIPv4v6,
			net.ParseIP("1.2.3.4"),
			net.ParseIP("::1"),
		},
		{
			"PDNType NonIP",
			ie.NewPDNAddressAllocation(""),
			gtpprime.PDNTypeNonIP,
			nil,
			nil,
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			pdnType := c.paa.MustPDNType()
			if diff := cmp.Diff(pdnType, c.pdnType); diff != "" {
				t.Error(diff)
			}

			ipv4, _ := c.paa.IPv4()
			if diff := cmp.Diff(ipv4, c.ipv4); diff != "" {
				t.Error(diff)
			}

			ipv6, _ := c.paa.IPv6()
			if diff := cmp.Diff(ipv6, c.ipv6); diff != "" {
				t.Error(diff)
			}

			ip, err := c.paa.IP()
			if err == nil {
				v := ipv4
				if pdnType == gtpprime.PDNTypeIPv6 {
					v = ipv6
				}
				if diff := cmp.Diff(ip, v); diff != "" {
					t.Error(diff)
				}
			} else if err != io.ErrUnexpectedEOF {
				t.Error(err)
			}
		})
	}
}
