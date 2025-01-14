package template

import (
	"log/slog"
	"os"
	"text/template"
)

type OVPNConfig struct {
	Hostname string `json:"hostname"`
	Protocol string `json:"protocol"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
}

func WriteOVPNFile(hostname string, ip string, protocol string, port int) error {
	ovpnConfig := OVPNConfig{
		Hostname: hostname,
		Protocol: protocol,
		IP:       ip,
		Port:     port,
	}

	templateName := "template.ovpn.tmpl"
	dirName := "/etc/nordvpn"
	dir := os.DirFS(dirName)
	tmpl, err := template.New(templateName).ParseFS(dir, templateName)
	if err != nil {
		slog.Warn("Template not found", "dir", dirName, "error", err)
		dirName = "templates"
		dir = os.DirFS("templates")
		tmpl, err = template.New(templateName).ParseFS(dir, templateName)
	}
	if err != nil {
		slog.Warn("Template not found", "dir", dirName, "error", err)
		return err
	} else {
		slog.Info("Template found", "dir", dirName)
	}
	outputName := "/tmp/nordvpn.ovpn"
	out, err := os.Create(outputName)
	if err != nil {
		return err
	}
	err = tmpl.Execute(out, ovpnConfig)
	if err != nil {
		return err
	}
	return nil
}
