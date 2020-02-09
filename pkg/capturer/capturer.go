package capturer

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func CheckPerms(eth string) error {
	err := checkInterface(eth)
	if err != nil {
		return err
	}
	var (
		deviceName  string        = eth
		snapshotLen int32         = 1024
		promiscuous bool          = false
		timeout     time.Duration = -1 * time.Second
		handle      *pcap.Handle
	)
	// Open the device for capturing
	handle, err = pcap.OpenLive(deviceName, snapshotLen, promiscuous, timeout)
	if err != nil {
		return err
	}
	defer handle.Close()
	return nil
}

func OpenBrowser(link, eth, proxyString string, timeout int, headless bool) (string, error) {
	port, err := pickUnusedPort()
	if err != nil {
		return "", err
	}
	var opts []selenium.ServiceOption
	service, err := selenium.NewChromeDriverService("chromedriver", port, opts...)
	if err != nil {
		return "", err
	}
	defer service.Stop()

	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	args := []string{"--incognito"}
	if headless {
		args = append(args, "--headless")
	}
	if proxyString != "" {
		proxyArgStr := "--proxy-server=" + proxyString
		args = append(args, proxyArgStr)
	}
	ua := "--user-agent=" + getUserAgent()
	args = append(args, ua)
	caps.AddChrome(chrome.Capabilities{
		Args: args,
	})
	wd, err := selenium.NewRemote(caps, "http://127.0.0.1:"+strconv.Itoa(port)+"/wd/hub")
	if err != nil {
		return "", err
	}
	_ = wd.SetPageLoadTimeout(time.Duration(timeout) * time.Second)
	wd.Refresh()
	stop := make(chan struct{})
	dirpath := "captured_traffic/" + getHostname(link) + "/"
	os.MkdirAll(dirpath, os.ModePerm)
	// Open output pcap file and write header
	savepath := dirpath + strconv.Itoa(randomInt(100000, 999999)) + ".pcap"
	go captureTraffic(savepath, eth, stop)
	wd.Get(link)
	_, err = wd.FindElement(selenium.ByXPATH, `/html/body`)
	if err != nil {
		close(stop)
		os.Remove(savepath)
		return "", err
	}
	close(stop)
	_ = wd.Quit()

	return savepath, nil
}

func checkInterface(i string) error {
	// Find all devices
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return err
	}
	localip, err := externalIP()
	if err != nil {
		return err
	}
	var ok bool
	var devicesNames []string
	var suggestion string
	// Print device information
	for _, device := range devices {
		devicesNames = append(devicesNames, device.Name)
		for _, address := range device.Addresses {
			if localip == address.IP.String() {
				suggestion = device.Name
			}
		}
		if i == device.Name {
			ok = true
		}
	}
	if !ok {
		err := fmt.Errorf("no such device %v; list of devices: %v; better use this one: %v", i, devicesNames, suggestion)
		return err
	}
	return nil
}

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

func captureTraffic(savepath, eth string, stop chan struct{}) {
	var (
		deviceName  string = eth
		snapshotLen int32  = 1024
		promiscuous bool   = false
		// err         error
		timeout time.Duration = -1 * time.Second
		handle  *pcap.Handle
	)

	f, _ := os.Create(savepath)
	defer f.Close()
	w := pcapgo.NewWriter(f)
	w.WriteFileHeader(uint32(snapshotLen), layers.LinkTypeEthernet)

	// Open the device for capturing
	handle, _ = pcap.OpenLive(deviceName, snapshotLen, promiscuous, timeout)
	defer handle.Close()

	var filter string = "tcp and port 443"
	handle.SetBPFFilter(filter)
	// Start processing packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		if isClosed(stop) {
			break
		}
	}
}

func randomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

func pickUnusedPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	port := l.Addr().(*net.TCPAddr).Port
	if err := l.Close(); err != nil {
		return 0, err
	}
	return port, nil
}

func getHostname(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		log.Fatal(err)
	}
	return u.Hostname()
}

func isClosed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

func getUserAgent() string {
	useragents := []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3835.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3831.6 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3818.0 Safari/537.36 Edg/77.0.189.3",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3790.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3782.0 Safari/537.36 Edg/76.0.152.0"}
	rand.Seed(time.Now().Unix())
	return useragents[rand.Intn(len(useragents))]
}
