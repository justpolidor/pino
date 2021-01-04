# PINO

pino exports disk stats about a specific volume (path) to a prometheus compatible endpoint ``` /metrics ```.
It exports: available disk space, used disk space, total disk space in MB.

### Usage
```
docker run jpolidor/pino:1.0.0 -h
docker run jpolidor/pino:1.0.0 -isDebug=true -interval=1000 -path="/mnt/1" -port=8080 

```