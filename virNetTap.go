package virNetTap

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// VirNetTap virNetTap
type VirNetTap struct{}

// Stub Stub stats
type Stub struct {
	Bytes uint64
	Pkts  uint64
	Errs  uint64
	Drops uint64
}

// InterfaceStats Interface stats
type InterfaceStats struct {
	VIF string
	IN  Stub
	OUT Stub
}

// GetAllVifStats Returns stats for all active virtual interfaces
func (v *VirNetTap) GetAllVifStats() (ret map[string]InterfaceStats, err error) {
	var fp *os.File
	var scanner *bufio.Scanner

	ret = make(map[string]InterfaceStats)

	if fp, err = os.Open("/proc/net/dev"); err != nil {
		return
	}

	scanner = bufio.NewScanner(fp)

	var i int
	var ifname string
	var bytesIn uint64
	var bytesOut uint64
	var pktsIn uint64
	var pktsOut uint64
	var errsIn uint64
	var errsOut uint64
	var dropsIn uint64
	var dropsOut uint64
	var dummy uint64

	for scanner.Scan() {
		i, err = fmt.Sscanf(scanner.Text(),
			"%s %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d",
			&ifname,
			&bytesOut, &pktsOut, &errsOut, &dropsOut,
			&dummy, &dummy, &dummy, &dummy,
			&bytesIn, &pktsIn, &errsIn, &dropsIn,
			&dummy, &dummy, &dummy, &dummy)

		if i < 17 || err != nil {
			continue
		}

		ifname = strings.TrimRight(ifname, ":")

		ret[ifname] = InterfaceStats{
			VIF: ifname,
			IN: Stub{
				Bytes: bytesIn,
				Pkts:  pktsIn,
				Errs:  errsIn,
				Drops: dropsIn,
			},
			OUT: Stub{
				Bytes: bytesOut,
				Pkts:  pktsOut,
				Errs:  errsOut,
				Drops: dropsOut,
			},
		}

	}

	if err = scanner.Err(); err != nil {
		err = fmt.Errorf("Scanner error: %s", err.Error())
		return
	}

	return
}
