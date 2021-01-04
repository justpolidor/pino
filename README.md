# PINO

pino exports disk stats about a specific volume (path) to a prometheus compatible endpoint ``` /metrics ```.
It exports: available disk space, used disk space, total disk space in MB.

### Usage

First export an environment variable named ```APP_NAME``` that will prefix the prometheus metric

```
docker run jpolidor/pino:1.0.0 -h
docker run -e APP_NAME=pino -p 8080:8080 jpolidor/pino:1.0.0 -isDebug=true -interval=1000 -path="/mnt/1" -port=8080 

```
### Example output 

```
curl <host>:<port>/metrics
[...]
# HELP pino_disk_free_size Total free size of the disk in MB
# TYPE pino_disk_free_size gauge
pino_disk_free_size 44854.88671875
# HELP pino_disk_total_size Total size of the disk in MB
# TYPE pino_disk_total_size gauge
pino_disk_total_size 59819.81640625
# HELP pino_disk_used_size Total usage of the disk in MB
# TYPE pino_disk_used_size gauge
pino_disk_used_size 14964.9296875
```
