package network

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"simple-docker/common"
	"simple-docker/container"
	"strings"
	"text/tabwriter"
)

var (
	drivers  = map[string]Driver{}
	networks = map[string]*Network{}
)

// Network 网络
type Network struct {
	Name    string
	IpRange *net.IPNet
	Driver  string
}

// Endpoint 网络端点
type Endpoint struct {
	ID          string           `json:"id"`
	Device      netlink.Veth     `json:"dev"`
	IPAddress   net.IP           `json:"ip"`
	MacAddress  net.HardwareAddr `json:"mac"`
	Network     *Network
	PortMapping []string
}

// Driver 网络驱动接口
type Driver interface {
	// Name 驱动名
	Name() string
	// Create 创建网络
	Create(subnet string, name string) (*Network, error)
	// Delete 删除网络
	Delete(network Network) error
	// Connect 连接容器网络端点到网络
	Connect(network *Network, endpoint *Endpoint) error
	// Disconnect 从网络上移除容器网络端点
	Disconnect(network Network, endpoint *Endpoint) error
}

func (nw *Network) dump(dumpPath string) error {
	if _, err := os.Stat(dumpPath); err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(dumpPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	nwPath := path.Join(dumpPath, nw.Name)
	nwFile, err := os.OpenFile(nwPath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logrus.Errorf("error：%v", err)
		return err
	}
	defer func(nwFile *os.File) {
		err := nwFile.Close()
		if err != nil {

		}
	}(nwFile)

	nwJson, _ := json.Marshal(nw)
	_, err = nwFile.Write(nwJson)
	if err != nil {
		logrus.Errorf("write network file, error: %v", err)
		return err
	}
	return nil
}

func (nw *Network) remove(dumpPath string) error {
	if _, err := os.Stat(path.Join(dumpPath, nw.Name)); err != nil && os.IsNotExist(err) {
		return nil
	}
	return os.Remove(path.Join(dumpPath, nw.Name))
}

func (nw *Network) load(dumpPath string) error {
	nwConfigFile, err := os.Open(dumpPath)
	if err != nil {
		return err
	}
	defer func(nwConfigFile *os.File) {
		err := nwConfigFile.Close()
		if err != nil {

		}
	}(nwConfigFile)

	nwJson := make([]byte, 2000)
	n, err := nwConfigFile.Read(nwJson)
	if err != nil {
		return err
	}

	err = json.Unmarshal(nwJson[:n], nw)
	if err != nil {
		logrus.Errorf("json unmarshal nw info, err: %v", err)
		return err
	}
	return nil
}

// Init 初始化网络驱动
func Init() error {
	var bridgeDriver = BridgeNetworkDriver{}
	drivers[bridgeDriver.Name()] = &bridgeDriver

	if _, err := os.Stat(common.DefaultNetworkPath); err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(common.DefaultNetworkPath, os.ModePerm); err != nil {
			return err
		}
	}
	// 递归遍历目录
	err := filepath.Walk(common.DefaultNetworkPath, func(nwPath string, info os.FileInfo, err error) error {
		if strings.HasSuffix(nwPath, "/") {
			return nil
		}
		_, nwName := path.Split(nwPath)
		nw := &Network{
			Name: nwName,
		}

		if err := nw.load(nwPath); err != nil {
			logrus.Errorf("error load network: %s", err)
		}

		networks[nwName] = nw
		return nil
	})

	if err != nil {
		logrus.Errorf("file path walk, err: %v", err)
		return err
	}
	logrus.Infof("networks: %v", networks)

	return nil
}

// CreateNetwork 创建网络
func CreateNetwork(driver, subnet, name string) error {
	_, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		logrus.Errorf("parse cidr, err: %v", err)
		return err
	}
	// 分配一个IP地址
	ip, err := ipAllocator.Allocate(ipNet)
	if err != nil {
		logrus.Errorf("allocate ip, err: %v", err)
	}
	ipNet.IP = ip

	// 创建网络
	nw, err := drivers[driver].Create(ipNet.String(), name)
	if err != nil {
		return err
	}

	// 将对象保存到文件中
	err = nw.dump(common.DefaultNetworkPath)
	if err != nil {
		logrus.Errorf("dump network, err: %v", err)
		return err
	}
	return nil
}

// Connect 连接网络
func Connect(networkName string, containerInfo *container.Info) error {
	network, ok := networks[networkName]
	if !ok {
		return fmt.Errorf("no Such network: %s", networkName)
	}

	// 分配容器IP地址
	ip, err := ipAllocator.Allocate(network.IpRange)
	if err != nil {
		return err
	}

	// 创建网络端点
	ep := &Endpoint{
		ID:          fmt.Sprintf("%s-%s", containerInfo.Id, networkName),
		IPAddress:   ip,
		Network:     network,
		PortMapping: containerInfo.PortMapping,
	}
	// 调用网络驱动挂载和配置网络端点
	if err = drivers[network.Driver].Connect(network, ep); err != nil {
		return err
	}
	// 给容器的namespace配置容器网络设备IP地址
	if err = configEndpointIpAddressAndRoute(ep, containerInfo); err != nil {
		return err
	}

	// 配置端口映射
	err = configPortMapping(ep, containerInfo)
	if err != nil {
		logrus.Errorf("config port mapping, err: %v", err)
		return err
	}
	return nil
}

// 给容器的namespace配置容器网络设备IP地址
func configEndpointIpAddressAndRoute(ep *Endpoint, cinfo *container.Info) error {
	peerLink, err := netlink.LinkByName(ep.Device.PeerName)
	if err != nil {
		logrus.Errorf("fail config endpoint: %v", err)
		return err
	}
	defer enterContainerNetns(&peerLink, cinfo)()

	interfaceIP := *ep.Network.IpRange
	interfaceIP.IP = ep.IPAddress

	if err = setInterfaceIP(ep.Device.PeerName, interfaceIP.String()); err != nil {
		return fmt.Errorf("%v,%s", ep.Network, err)
	}

	if err = setInterfaceUP(ep.Device.PeerName); err != nil {
		return err
	}

	if err = setInterfaceUP("lo"); err != nil {
		return err
	}

	_, cidr, _ := net.ParseCIDR("0.0.0.0/0")

	defaultRoute := &netlink.Route{
		LinkIndex: peerLink.Attrs().Index,
		Gw:        ep.Network.IpRange.IP,
		Dst:       cidr,
	}

	if err = netlink.RouteAdd(defaultRoute); err != nil {
		return err
	}

	return nil
}

func enterContainerNetns(enLink *netlink.Link, cinfo *container.Info) func() {
	f, err := os.OpenFile(fmt.Sprintf("/proc/%s/ns/net", cinfo.Pid), os.O_RDONLY, 0)
	if err != nil {
		logrus.Errorf("error get container net namespace, %v", err)
	}

	nsFD := f.Fd()
	runtime.LockOSThread()

	// 修改 veth peer 另外一端移到容器的 namespace 中
	if err = netlink.LinkSetNsFd(*enLink, int(nsFD)); err != nil {
		logrus.Errorf("set link netns, err: %v", err)
	}

	// 获取当前的网络 namespace
	origns, err := netns.Get()
	if err != nil {
		logrus.Errorf("get current netns, err: %v", err)
	}

	// 设置当前进程到新的网络 namespace，并在函数执行完成之后再恢复到之前的 namespace
	if err = netns.Set(netns.NsHandle(nsFD)); err != nil {
		logrus.Errorf("error set netns, %v", err)
	}
	return func() {
		err := netns.Set(origns)
		if err != nil {
			return
		}
		err = origns.Close()
		if err != nil {
			return
		}
		runtime.UnlockOSThread()
		err = f.Close()
		if err != nil {
			return
		}
	}
}

// 配置端口映射关系
func configPortMapping(ep *Endpoint, cinfo *container.Info) error {
	for _, pm := range ep.PortMapping {
		portMapping := strings.Split(pm, ":")
		if len(portMapping) != 2 {
			logrus.Errorf("port mapping format error, %v", pm)
			continue
		}
		iptablesCmd := fmt.Sprintf("-t nat -A PREROUTING -p tcp -m tcp --dport %s -j DNAT --to-destination %s:%s",
			portMapping[0], ep.IPAddress.String(), portMapping[1])
		cmd := exec.Command("iptables", strings.Split(iptablesCmd, " ")...)
		// err := cmd.Run()
		output, err := cmd.Output()
		if err != nil {
			logrus.Errorf("iptables Output, %v", output)
			continue
		}
	}
	return nil
}

// ListNetwork 遍历网络
func ListNetwork() {
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	_, _ = fmt.Fprint(w, "NAME\tIpRange\tDriver\n")
	for _, nw := range networks {
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n",
			nw.Name,
			nw.IpRange.String(),
			nw.Driver,
		)
	}
	if err := w.Flush(); err != nil {
		logrus.Errorf("Flush error %v", err)
		return
	}
}

// DeleteNetwork 删除网络
func DeleteNetwork(networkName string) error {
	nw, ok := networks[networkName]
	if !ok {
		return fmt.Errorf("no Such Network: %s", networkName)
	}

	if err := ipAllocator.Release(nw.IpRange, &nw.IpRange.IP); err != nil {
		return fmt.Errorf("remove network gateway ip, err: %v", err)
	}

	if err := drivers[nw.Driver].Delete(*nw); err != nil {
		return fmt.Errorf("remove network driver, err: %v", err)
	}

	return nw.remove(common.DefaultNetworkPath)
}
