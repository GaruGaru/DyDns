# DyDns
## Dynamic DNS Client for Namecheap

### Installation

```bash
go get github.com/GaruGaru/DyDns

```

### Usage 

```bash
dydns --domain=<domain> --host=<entry> --password=<dynamic_dns_psw> --delay=<update_delay_min>  
```
 
### Docker 

```bash
docker run -e "DOMAIN=<domain>" -e "HOST=<entry>" -e "PASSWORD=<dynamic_dns_psw>" garugaru/dydns
```




