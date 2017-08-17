package metrics

import "github.com/prometheus/client_golang/prometheus"

const namespace = "gpldr"

// UploadedFilesTotal is a prometheus counter to monitor the number of uploaded
// files
var UploadedFilesTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "uploaded_files_total",
		Help:      "Count of files uploaded",
	},
	[]string{"ip"},
)

// UploadedFilesSizeTotal is a prometheus counter to monitor the size of
// uploaded files
var UploadedFilesSizeTotal = prometheus.NewCounter(
	prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "uploaded_files_size_total",
		Help:      "Size of uploaded files",
	},
)

func init() {
	prometheus.MustRegister(UploadedFilesSizeTotal, UploadedFilesTotal)
}
