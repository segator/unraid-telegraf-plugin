# unraid-telegraf-plugin
Export unraid metrics as a influxDB line Protocol.
influxDB line protocol can be handled by telegraf using inputs.exec with dataType influx

##Standalone
As a standalone binary this app will print on the stdout the influxDB protocol line with the results.

### Build
```bash
go get -d ./
go build -o unraid-influxdb-line
```

### Run
```bash
chmod +x unraid-influxdb-line
#Get Help
./unraid-influxdb-line
Usage of ./unraid-influxdb-line:
  -modules string
        modules to extract data separated with coma, default(disk,ups)  available disk,ups (default "disk,ups")
  -unraid-disks-path string
        Path where disks.ini is located, by default /var/local/emhttp/disks.ini (default "/var/local/emhttp/disks.ini")
```

## Docker with telegraf
Telegraf execute the standalone binary extract the data and can be processed as you want.
### Build
```bash
docker build -t segator/unraid-telegraf .
```

### Run
You need to edit your telegraf.conf adding on the inputs.exec the execution of /app/unraid-influxdb-line
#### telegraf.conf Example
```
[[inputs.exec]]
   commands = [
     "/app/unraid-influxdb-line --unraid-disks-path /rootfs/var/local/emhttp/disks.ini"
   ]
```

After this you can run this docker(Remember change the path of the telegraf.conf
```
docker run -d --name=telegraf --restart=always --privileged --net=host -v /boot/config/telegraf/telegraf.conf:/etc/telegraf/telegraf.conf:ro \
           -v /var/run/utmp:/var/run/utmp \
           -v /var/run/docker.sock:/var/run/docker.sock \
           -v /:/rootfs \
           -v /sys:/rootfs/sys \
           -v /etc:/rootfs/etc \
           -v /proc:/rootfs/proc \
           -e HOST_PROC=/rootfs/proc \
           -e HOST_SYS=/rootfs/sys \
           -e HOST_ETC=/rootfs/etc \
           -e HOST_MOUNT_PREFIX=/rootfs \
           segator/unraid-telegraf
```
