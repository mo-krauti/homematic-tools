# homematic-tools

## develop

Set values of template.env in your env.  
Run standalone tool:
```bash
./scripts/run.sh homematic-powerfox
```

## deployment

The deployment assumes you have set up `homematic` as an ssh alias in your ssh config.

### cross compile
```bash
./scripts/build-arm.sh homematic-powerfox
./scripts/build-arm.sh homematic-sma-web
./scripts/build-arm.sh homematic-stromgedacht
```

### environment
```bash
cp template.env deployment/supervisor/credentials.env
```
Now set the variables to your credentials and device hostnames.

### supervisord
```bash
./scripts/deploy-supervisord.sh
```

### tools
```bash
./scripts/deploy.sh homematic-powerfox
./scripts/deploy.sh homematic-sma-web
./scripts/deploy.sh homematic-stromgedacht
```

### manage
Start supervisord
```bash
ssh homematic sh /usr/local/addons/homematic-tools/supervisor/start.sh
```
Use supervisord ctl
```bash
ssh homematic /usr/local/addons/homematic-tools/supervisor/supervisord ctl --help
```
Access the web interface at `<device-ip>:9001`

### auto-start
Add to `rc.local` on homematic device:
```bash
/bin/sh /usr/local/addons/homematic-tools/supervisor/start.sh
```
