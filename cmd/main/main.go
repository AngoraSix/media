package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"

	media "angorasix.com/media"
	config "angorasix.com/media/config"
	images "angorasix.com/media/gen/images"
	strategies "angorasix.com/media/strategies"
)

// Build id
var Build string

func main() {
	// Define command line flags, add any other flag required to configure the
	// service.
	var (
		hostF              = flag.String("host", config.DefaultServerConfig.Hostname, "Server host (valid values: 0.0.0.0)")
		domainF            = flag.String("domain", "", "Host domain name (overrides host domain specified in service design)")
		httpPortF          = flag.String("http-port", config.DefaultServerConfig.Port, "HTTP port (overrides host HTTP port specified in service design)")
		secureF            = flag.Bool("secure", false, "Use secure scheme (https or grpcs)")
		dbgF               = flag.Bool("debug", false, "Log request and response bodies")
		strategy           = flag.String("strategy", config.DefaultServiceConfig.Strategy, "The image persistence strategy")
		staticDir          = flag.String("staticDir", config.DefaultServiceConfig.StaticDir, "The static directory that contains the images")
		uploadDir          = flag.String("uploadDir", config.DefaultServiceConfig.UploadDir, "The directory to upload the images to")
		bucketName         = flag.String("bucketName", config.DefaultServiceConfig.BucketName, "The bucket name")
		projectID          = flag.String("projectID", config.DefaultServiceConfig.ProjectID, "The project Id")
		storageAPIHost     = flag.String("storageAPIHost", config.DefaultServiceConfig.StorageAPIHost, "The Storage API Host")
		thumbnailMaxHeight = flag.Uint("thumbnailMaxHeight", config.DefaultServiceConfig.ThumbnailMaxHeight, "The max height to use when creating an image thumbnail")
		thumbnailMaxWidth  = flag.Uint("thumbnailMaxWidth", config.DefaultServiceConfig.ThumbnailMaxWidth, "The max width to use when creating an image thumbnail")
		thumbnailMaxSize   = flag.Uint("thumbnailMaxSize", config.DefaultServiceConfig.ThumbnailMaxSize, "The max size (in bytes) to use when creating an image thumbnail")
		thumbnailQuality   = flag.Uint("thumbnailQuality", config.DefaultServiceConfig.ThumbnailQuality, "The quality of the thumbnail to create (0-100)")
	)

	flag.Parse()

	// Setup logger. Replace logger with your own log package of choice.
	var (
		logger *log.Logger
	)
	{
		logger = log.New(os.Stderr, "[media] ", log.Ltime)
	}

	serviceConfig := config.ServiceConfig{
		Strategy:           *strategy,
		StaticDir:          *staticDir,
		UploadDir:          *uploadDir,
		BucketName:         *bucketName,
		ProjectID:          *projectID,
		StorageAPIHost:     *storageAPIHost,
		ThumbnailMaxHeight: *thumbnailMaxHeight,
		ThumbnailMaxWidth:  *thumbnailMaxWidth,
		ThumbnailMaxSize:   *thumbnailMaxSize,
		ThumbnailQuality:   *thumbnailQuality,
	}
	selectedStrategy, err := strategies.CreateStrategyFromConfig(&serviceConfig)

	if err != nil {
		panic(err)
	}

	logger.Printf("DEPLOYING BUILD[%s]", Build)

	// Initialize the services.
	var (
		imagesSvc images.Service
	)
	// GERGERGER
	{
		imagesSvc = media.NewImages(logger, selectedStrategy, &serviceConfig)
	}

	// Wrap the services in endpoints that can be invoked from other services
	// potentially running in different processes.
	var (
		imagesEndpoints *images.Endpoints
	)
	{
		imagesEndpoints = images.NewEndpoints(imagesSvc)
	}

	// Create channel used by both the signal handler and server goroutines
	// to notify the main goroutine when to stop the server.
	errc := make(chan error)

	// Setup interrupt handler. This optional step configures the process so
	// that SIGINT and SIGTERM signals cause the services to stop gracefully.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		errc <- fmt.Errorf("%s", <-c)
	}()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	// Start the servers and send errors (if any) to the error channel.
	switch *hostF {
	case "0.0.0.0":
		{
			addr := "http://0.0.0.0:80"
			u, err := url.Parse(addr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid URL %#v: %s\n", addr, err)
				os.Exit(1)
			}
			if *secureF {
				u.Scheme = "https"
			}
			if *domainF != "" {
				u.Host = *domainF
			}
			if *httpPortF != "" {
				h := strings.Split(u.Host, ":")[0]
				u.Host = h + ":" + *httpPortF
			} else if u.Port() == "" {
				u.Host += ":80"
			}
			handleHTTPServer(ctx, u, imagesEndpoints, &wg, errc, logger, *dbgF)
		}

	default:
		fmt.Fprintf(os.Stderr, "invalid host argument: %q (valid hosts: 0.0.0.0)\n", *hostF)
	}

	// Wait for signal.
	logger.Printf("exiting (%v)", <-errc)

	// Send cancellation signal to the goroutines.
	cancel()

	wg.Wait()
	logger.Println("exited")
}
