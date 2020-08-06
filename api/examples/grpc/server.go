// +build examples

package grpcserver

import (
	"flag"
	stdlog "log"
	"net"
	"os"

	conf "github.com/cryptogarageinc/server-common-go/pkg/configuration"
	"github.com/cryptogarageinc/server-common-go/pkg/database/orm"
	"github.com/cryptogarageinc/server-common-go/pkg/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	configPath = flag.String("config", "", "Path to the configuration file to use.")
	appName    = flag.String("appname", "", "The name of the application. Will be use as a prefix for environment variables.")
	envname    = flag.String("e", "", "environment (ex., \"development\"). Should match with the name of the configuration file.")
	migrate    = flag.Bool("migrate", false, "If set performs a db migration before starting.")
)

// Config contains the configuration parameters for the server.
type Config struct {
	Address  string `configkey:"server.address" validate:"required"`
	TLS      bool   `configkey:"server.tls"`
	CertFile string `configkey:"server.certfile" validate:"required_with=TLS"`
	KeyFile  string `configkey:"server.keyfile" validate:"required_with=TLS"`
}

func newInitializedLog(config *conf.Configuration) *log.Log {
	logConfig := &log.Config{}
	config.InitializeComponentConfig(logConfig)
	logger := log.NewLog(logConfig)
	logger.Initialize()
	return logger
}

func newInitializedOrm(config *conf.Configuration, log *log.Log) *orm.ORM {
	ormConfig := &orm.Config{}
	config.InitializeComponentConfig(ormConfig)
	ormInstance := orm.NewORM(ormConfig, log)
	err := ormInstance.Initialize()

	if err != nil {
		panic("Could not initialize database.")
	}

	return ormInstance
}

func main() {
	flag.Parse()

	if *configPath == "" {
		stdlog.Fatal("No configuration path specified")
	}

	if *appName == "" {
		stdlog.Fatal("No configuration name specified")
	}

	if *envname != "" {
		os.Setenv("GRPC_SERVER_ENV", *envname)
	}

	config := conf.NewConfiguration(*appName, *envname, []string{*configPath})
	err := config.Initialize()

	if err != nil {
		stdlog.Fatalf("Could not read configuration %v.", err)
	}

	serverConfig := &Config{}

	config.InitializeComponentConfig(serverConfig)

	lis, err := net.Listen("tcp", serverConfig.Address)
	if err != nil {
		stdlog.Fatalf("failed to listen: %v", err)
	}

	opts := make([]grpc.ServerOption, 0)

	// should be use to ensure default authentication
	if serverConfig.TLS {
		certFile := serverConfig.CertFile
		keyFile := serverConfig.KeyFile
		if certFile == "" {
			stdlog.Fatal("Need to provide the path to the certificate file")
		}
		if keyFile == "" {
			stdlog.Fatal("Need to provide the path to the key file")
		}
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			stdlog.Fatalf("Failed to generate credentials %v", err)
		}
		opts = append(opts, grpc.Creds(creds))
	}

	logInstance := newInitializedLog(config)
	ormInstance := newInitializedOrm(config, logInstance)

	if *migrate {
		err := doMigration(logInstance, ormInstance)

		if err != nil {
			stdlog.Fatalf("Failed to apply migration %v", err)
		}
	}

	grpcServer := grpc.NewServer(opts...)

	// register service here ...
	// myservice.RegisterXServer(grpcServer, someControllerImplementation)

	stdlog.Printf("Ready to listen on %v", serverConfig.Address)
	grpcServer.Serve(lis)
}

func doMigration(l *log.Log, o *orm.ORM) error {
	migrator := orm.NewMigrator(
		o,
	)

	return migrator.Initialize()
}
