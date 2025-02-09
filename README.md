# wifimelon
wifimelon is a lightweight, open source, Golang based Wi-Fi hacking framework focused on the use of several concurrent Wi-Fi adapters capable of packet injection for all kinds of shenanigans. This repository is currently under active development. Use only on test networks that you own.


# V1.1 usage: 
Set up config.conf with your Wi-Fi adapters (leave any that you can't fill). Identify the network BSSID('s) that you want to deauth as well as their channels (```sudo airodump-ng wlanX```). Set your Wi-Fi adapters, in order, to said channels with ```sudo iwconfig wlanX channel XX```, run ```sudo airmon-ng check kill``` if you're also using your integrated Wi-Fi adapter, and then hammer away: ```sudo go run main.go mdeauth BS:SI:D1 BS:SI:D2 BS:SI:D3```
