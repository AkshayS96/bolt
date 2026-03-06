package cmd

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"bolt/internal/utils"

	"github.com/spf13/cobra"
)

func init() {
	dnsCmd := &cobra.Command{
		Use:   "dns",
		Short: "DNS lookup tool",
	}

	dnsLookupCmd := &cobra.Command{
		Use:   "lookup <domain>",
		Short: "Resolve a domain to IP addresses",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			domain := args[0]
			utils.Key.Printf("DNS lookup: %s\n\n", domain)

			ips, err := net.LookupIP(domain)
			if err != nil {
				utils.PrintError("Lookup failed: " + err.Error())
				return
			}
			for _, ip := range ips {
				if ip.To4() != nil {
					utils.PrintKeyValue("A:", ip.String())
				} else {
					utils.PrintKeyValue("AAAA:", ip.String())
				}
			}

			cname, err := net.LookupCNAME(domain)
			if err == nil && cname != "" {
				utils.PrintKeyValue("CNAME:", cname)
			}

			mxRecords, err := net.LookupMX(domain)
			if err == nil {
				for _, mx := range mxRecords {
					utils.PrintKeyValue("MX:", fmt.Sprintf("%s (priority: %d)", mx.Host, mx.Pref))
				}
			}

			nsRecords, err := net.LookupNS(domain)
			if err == nil {
				for _, ns := range nsRecords {
					utils.PrintKeyValue("NS:", ns.Host)
				}
			}

			txtRecords, err := net.LookupTXT(domain)
			if err == nil {
				for _, txt := range txtRecords {
					utils.PrintKeyValue("TXT:", txt)
				}
			}
		},
	}

	dnsCmd.AddCommand(dnsLookupCmd)

	// --- Ping ---
	pingCmd := &cobra.Command{
		Use:   "ping <host>",
		Short: "Check if a host is reachable (TCP ping)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			host := args[0]
			if !strings.Contains(host, ":") {
				host = host + ":80"
			}
			start := time.Now()
			conn, err := net.DialTimeout("tcp", host, 5*time.Second)
			elapsed := time.Since(start)
			if err != nil {
				utils.PrintError(fmt.Sprintf("Host unreachable: %s (%s)", host, err.Error()))
				return
			}
			conn.Close()
			utils.PrintSuccess(fmt.Sprintf("Host reachable: %s (%s)", host, elapsed.Round(time.Millisecond)))
		},
	}

	// --- Port ---
	portCmd := &cobra.Command{
		Use:   "port",
		Short: "Port utilities (check, kill)",
	}

	portCheckCmd := &cobra.Command{
		Use:   "check <port>",
		Short: "Check if a port is in use",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			addr := ":" + args[0]
			ln, err := net.Listen("tcp", addr)
			if err != nil {
				utils.PrintWarning(fmt.Sprintf("Port %s is IN USE", args[0]))
			} else {
				ln.Close()
				utils.PrintSuccess(fmt.Sprintf("Port %s is AVAILABLE", args[0]))
			}
		},
	}

	portKillCmd := &cobra.Command{
		Use:   "kill <port>",
		Short: "Kill process running on a port",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			port := args[0]
			var killCmd *exec.Cmd
			if runtime.GOOS == "windows" {
				killCmd = exec.Command("cmd", "/c", fmt.Sprintf("for /f \"tokens=5\" %%a in ('netstat -aon ^| findstr :%s') do taskkill /F /PID %%a", port))
			} else {
				killCmd = exec.Command("sh", "-c", fmt.Sprintf("lsof -ti:%s | xargs kill -9 2>/dev/null", port))
			}
			output, err := killCmd.CombinedOutput()
			if err != nil {
				utils.PrintWarning(fmt.Sprintf("No process found on port %s (or already killed)", port))
			} else {
				utils.PrintSuccess(fmt.Sprintf("Killed process on port %s", port))
				if len(strings.TrimSpace(string(output))) > 0 {
					fmt.Println(string(output))
				}
			}
		},
	}

	portCmd.AddCommand(portCheckCmd, portKillCmd)

	// --- IP ---
	ipCmd := &cobra.Command{
		Use:   "ip",
		Short: "Show local IP addresses",
		Run: func(cmd *cobra.Command, args []string) {
			addrs, err := net.InterfaceAddrs()
			if err != nil {
				utils.PrintError("Failed to get interfaces: " + err.Error())
				return
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						utils.PrintKeyValue("IPv4:", ipnet.IP.String())
					}
				}
			}
		},
	}

	rootCmd.AddCommand(dnsCmd, pingCmd, portCmd, ipCmd)
}
