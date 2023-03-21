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

	"github.com/robfig/cron/v3"
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

func webHook(public_ip, domain, text string, ip_strings []string) {
	data := map[string]interface{}{
		"ipv6": map[string]string{
			// "text":    "Detect the current server's DDNS status: ",
			"text":    text,
			"result":  "Success",
			"addr":    public_ip,
			"domains": domain,
			"resolve": strings.Join(ip_strings, ","),
		},
	}
	DDNS_bot_webhook := sendUrl("<DDNS_BOT_WEBHOOK>", "POST", data)
	fmt.Println(DDNS_bot_webhook)
}

func main() {

	// 新建一个定时任务对象
	// 根据cron表达式进行时间调度，cron可以精确到秒，大部分表达式格式也是从秒开始。
	//crontab := cron.New()  默认从分开始进行时间调度
	crontab := cron.New(cron.WithSeconds()) //精确到秒
	//定义定时器调用的任务函数
	task := func() {
		fmt.Println("hello world", time.Now())
		fmt.Println("It's time to execute the scheduled task!")
		// 在此处编写需要执行的任务代码
		domain := "<DOMAIN>"
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
		if public_ip != ip_strings[0] {
			text := "The public network ip is inconsistent with the ip of domain name resolution: "
			webHook(public_ip, domain, text, ip_strings)
			update_domain_url := "http://dynv6.com/api/update?hostname=" + domain + "&token=<dynv6_TOKEN>=" + public_ip + "&ipv6prefix="
			result := sendUrl(update_domain_url, "GET", nil)
			fmt.Println(result)
		} else {
			text := "Detect the current server's DDNS status:"
			webHook(public_ip, domain, text, ip_strings)
		}
	}
	//定时任务
	spec := "0 0 8,14,21,0 * * ?"
	// 添加定时任务,
	crontab.AddFunc(spec, task)
	// 启动定时器
	crontab.Start()
	// 定时任务是另起协程执行的,这里使用 select 简答阻塞.实际开发中需要
	// 根据实际情况进行控制
	select {} //阻塞主线程停止

}
