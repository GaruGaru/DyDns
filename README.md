# DyDns
## Dynamic DNS Client for Namecheap

### Installation

```bash
go get github.com/GaruGaru/DyDns
```

### Usage 

```bash
dydns --domain=<domain> --entries=<entry>,<entry2> --password=<dynamic_dns_psw> --delay=<update_delay_min>
```
 
### Docker 

```bash
docker run -e "DOMAIN=<domain>" -e "ENTRIES=<entry>,<entry2>" -e "PASSWORD=<dynamic_dns_psw>" garugaru/dydns
```




