# homematic-stromgedacht

## Quickstart

Create a system variable with the name `"stromgedacht_wert"` in your homematic system.

```bash
go build homematic-stromgedacht.go
chmod +x homematic-stromgedacht
./homematic-stromgedacht http://HOMEMATIC_IP HOMEMATIC_USER HOMEMATIC_PW POST_CODE
```
