package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

func resolveIPs(domain string) ([]net.IP, error) {
	// 解析IP地址
	addrs, err := net.LookupIP(domain)
	if err != nil {
		return nil, err
	}

	ips := make([]net.IP, 0, len(addrs))
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			ips = append(ips, ipv4)
		} else {
			ipv6 := addr.To16()
			if ipv6 != nil {
				ips = append(ips, ipv6)
			}
		}
	}

	return ips, nil
}

func sendUrl(url, method string, data map[string]interface{}) string {

	// 创建 HTTP 客户端
	client := &http.Client{}

	if method == "GET" {
		// 创建 GET 请求
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			panic(err)
		}
		// 发送请求
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		// 读取响应
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		// 打印响应内容
		// fmt.Println(string(body))
		return string(body)
	} else if method == "POST" {
		jsonData, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		// 创建 POST 请求
		req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", "application/json")
		// 发送请求
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		// 读取响应
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		// 打印响应内容
		return string(body)
	}
	return ""
}

func main() {
	fmt.Println("let's start the program: ")
	// 在每天的早上八点，下午两点和晚上九点检测DDNS是否可用
	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		now := time.Now()
		if now.Hour() == 8 || now.Hour() == 14 || now.Hour() == 21 {
			fmt.Println("It's time to execute the scheduled task!")
			// 在此处编写需要执行的任务代码
			domain := "DOMAIN"
			// 解析IP地址
			ips, err := resolveIPs(domain)
			if err != nil {
				panic(err)
			}
			var ip_strings []string
			for _, ip := range ips {
				ip_strings = append(ip_strings, ip.String())
			}
			public_ip := sendUrl("http://6.ipw.cn", "GET", nil)
			data := map[string]interface{}{
				"ipv6": map[string]string{
					"text":    "Detect the current server's DDNS status: ",
					"result":  "Success",
					"addr":    public_ip,
					"domains": domain,
					"resolve": strings.Join(ip_strings, ","),
				},
			}
			DDNS_bot_webhook := sendUrl("DDNS_BOT_WEBHOOK", "POST", data)
			fmt.Println(DDNS_bot_webhook)
		}
	}

}
