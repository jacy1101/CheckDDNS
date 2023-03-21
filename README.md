## 检测DDNS当前是否可用

- 支持Telegram
- 支持IPv6
- 增加域名解析
- 目前仅支持[dynv6](https://dynv6.com/)一家DDNS平台
- 每日早上8点、下午2点和晚上9点自动检测DDNS状态，配合[DDNS-GO](https://github.com/jeessy2/ddns-go)使用

>如果OpenWRT遇到时间不对的问题，参考[https://github.com/jeessy2/ddns-go/issues/497](https://github.com/jeessy2/ddns-go/issues/497)，安装`zoneinfo-asia`并重启OpenWRT。

需要添加[DDNSBOT](https://t.me/KDDNS_bot)

<img src="images/Snipaste_2023-03-20_17-46-50.png" width="auto"  height="auto"/>