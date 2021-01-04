package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"pino/data"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	syscall "golang.org/x/sys/unix"
)

var appName string = os.Getenv("APP_NAME")

var diskPath string
var port int
var scrapeInterval int
var isDebug bool

// Handy conversions
const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

var (
	diskSizeGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: appName,
			Name:      "disk_total_size",
			Help:      "Total size of the disk in MB",
		})
	diskFreeGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: appName,
			Name:      "disk_free_size",
			Help:      "Total free size of the disk in MB ",
		})
	diskUsedGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: appName,
			Name:      "disk_used_size",
			Help:      "Total usage of the disk in MB",
		})
)

func getDiskUsage(path string) (disk data.DiskStatus) {

	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		log.Fatal("Unable to get metadata")
		return
	}

	all := fs.Blocks * uint64(fs.Bsize)
	free := fs.Bavail * uint64(fs.Bsize)
	used := all - free

	return data.DiskStatus{
		All:  all,
		Free: free,
		Used: used,
	}

}

// Init parses cmd flags
func Init() {

	//log.SetFormatter(&log.JSONFormatter{})

	log.Debug("Parsing flags...")

	flag.IntVar(&port, "port", 8080, "The port where the webserver should bind")
	flag.StringVar(&diskPath, "path", "/", "The target path where interested to monitor")
	flag.IntVar(&scrapeInterval, "interval", 30000, "Interval between scrapes")
	flag.BoolVar(&isDebug, "isDebug", false, "Set to true to run the app in debug mode")

	flag.Parse()

	if isDebug {
		log.SetLevel(log.DebugLevel)
	}

	log.WithFields(log.Fields{
		"port":     port,
		"path":     diskPath,
		"interval": scrapeInterval,
		"isDebug":  isDebug,
	}).Info("Starting webserver")

}

func main() {

	Init()

	http.Handle("/metrics", promhttp.Handler())

	prometheus.MustRegister(diskSizeGauge)
	prometheus.MustRegister(diskFreeGauge)
	prometheus.MustRegister(diskUsedGauge)

	go func() {
		for {
			result := getDiskUsage(diskPath)

			log.WithFields(log.Fields{
				"disk usage": fmt.Sprintf("%+v", result),
			}).Debug("Got result")
			diskSizeGauge.Set(float64(result.All) / float64(MB))
			diskFreeGauge.Set(float64(result.Free) / float64(MB))
			diskUsedGauge.Set(float64(result.Used) / float64(MB))

			time.Sleep(time.Millisecond * time.Duration(scrapeInterval))
		}
	}()

	http.ListenAndServe(":"+strconv.Itoa(port), nil)

}
