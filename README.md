# http-forward
A http forward tool

### Usage
```bash
Usage of main.exe:
  -l string
        需要绑定本地的http服务HOST (default "0.0.0.0:9001")
  -r string
        需要转发的http服务HOST (default "www.baidu.com")
```

### example

```bash
main.exe -l 192.168.111.132:8080 -r www.baidu.com
```
then
I can use baidu by visiting 
`http://192.168.111.132:8080/`
or
`https://192.168.111.132:8080/`

> Please ignore my .idea file :)