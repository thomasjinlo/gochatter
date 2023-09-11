package utils

import (
    "fmt"
    "net"
    "os"
    "strings"
)

func GetPrivateIp() string {
    interfaces, err := net.Interfaces()
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }

    // Iterate through the network interfaces
    for _, iface := range interfaces {
        // Filter out loopback and non-up interfaces
        if strings.HasPrefix(iface.Name, "lo") || (iface.Flags&net.FlagUp == 0) {
            continue
        }

        // Get the addresses associated with the interface
        addrs, err := iface.Addrs()
        if err != nil {
            fmt.Println("Error:", err)
            continue
        }

        for _, addr := range addrs {
            if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
                if ipnet.IP.To4() != nil {
                    return ipnet.IP.String()
                }
            }
        }
    }
    return ""
}
