package main

import (
	"deficonnect/defipayapi/app"
	"fmt"
	"io"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/caarlos0/env"
	"github.com/decred/slog"
	flags "github.com/jessevdk/go-flags"
)

const (
	defaultConfigFilename = "coinzion.conf"
	sampleConfigFileName  = "./sample-coinzion.conf"
	defaultLogFilename    = "coinzion.log"
	defaultDataDirname    = "data"
	defaultLogLevel       = "info"
	defaultLogDirname     = "logs"
	defaultDbHost         = "0.0.0.0"
	defaultDbPort         = "5432"
	defaultDbUser         = "postgres"
	defaultDbPass         = "postgres"
	defaultDbName         = "coinzion"
)

var (
	defaultHomeDir           = "./"
	defaultConfigFile        = filepath.Join(defaultHomeDir, defaultConfigFilename)
	defaultLogDir            = filepath.Join(defaultHomeDir, defaultLogDirname)
	defaultDataDir           = filepath.Join(defaultHomeDir, defaultDataDirname)
	dcrdHomeDir              = "./"
	defaultDaemonRPCCertFile = filepath.Join(dcrdHomeDir, "rpc.cert")
	defaultMaxLogZips        = 16

	defaultHost               = "0.0.0.0"
	defaultHTTPProfPath       = "/p"
	defaultAPIProto           = "http"
	defaultPort               = "7070"
	defaultCacheControlMaxAge = 86400
	defaultServerHeader       = "coinzion"
)

type config struct {
	// General application behavior
	HomeDir      string `short:"A" long:"appdata" description:"Path to application home directory" env:"PDANALYTICS_APPDATA_DIR"`
	ConfigFile   string `short:"C" long:"configfile" description:"Path to configuration file" env:"PDANALYTICS_CONFIG_FILE"`
	DataDir      string `short:"b" long:"datadir" description:"Directory to store data" env:"PDANALYTICS_DATA_DIR"`
	LogDir       string `long:"logdir" description:"Directory to log output." env:"PDANALYTICS_LOG_DIR"`
	MaxLogZips   int    `long:"max-log-zips" description:"The number of zipped log files created by the log rotator to be retained. Setting to 0 will keep all."`
	OutFolder    string `short:"f" long:"outfolder" description:"Folder for file outputs" env:"PDANALYTICS_OUT_FOLDER"`
	ShowVersion  bool   `short:"V" long:"version" description:"Display version information and exit"`
	TestNet      bool   `long:"testnet" description:"Use the test network (default mainnet)" env:"PDANALYTICS_USE_TESTNET"`
	SimNet       bool   `long:"simnet" description:"Use the simulation test network (default mainnet)" env:"PDANALYTICS_USE_SIMNET"`
	DebugLevel   string `short:"d" long:"debuglevel" description:"Logging level {trace, debug, info, warn, error, critical}" env:"PDANALYTICS_LOG_LEVEL"`
	Quiet        bool   `short:"q" long:"quiet" description:"Easy way to set debuglevel to error" env:"PDANALYTICS_QUIET"`
	HTTPProfile  bool   `long:"httpprof" short:"p" description:"Start HTTP profiler." env:"PDANALYTICS_ENABLE_HTTP_PROFILER"`
	HTTPProfPath string `long:"httpprofprefix" description:"URL path prefix for the HTTP profiler." env:"PDANALYTICS_HTTP_PROFILER_PREFIX"`
	CPUProfile   string `long:"cpuprofile" description:"File for CPU profiling." env:"PDANALYTICS_CPU_PROFILER_FILE"`
	UseGops      bool   `short:"g" long:"gops" description:"Run with gops diagnostics agent listening. See github.com/google/gops for more information." env:"PDANALYTICS_USE_GOPS"`
	ReloadHTML   bool   `long:"reload-html" description:"Reload HTML templates on every request" env:"DCRDATA_RELOAD_HTML"`
	NoHttp       bool   `long:"nohttp" description:"Disables http server from running"`

	// Postgresql Configuration
	DBHost string `long:"dbhost" description:"Database host" env:"DBHOST"`
	DBPort string `long:"dbport" description:"Database port" env:"DBPORT"`
	DBUser string `long:"dbuser" description:"Database username" env:"DBUSER"`
	DBPass string `long:"dbpass" description:"Database password" env:"DBPASS"`
	DBName string `long:"dbname" description:"Database name" env:"DBNAME"`

	// EMAIL
	MailgunDomain string `long:"mailgudomain" env:"MAILGUNDOMAIN"`
	MailgunAPIKey string `long:"mailgunapikey" env:"MAILGUNAPIKEY"`

	NoAutoPayout string `long:"noautopayout" env:"NOAUTOPAYOUT"`

	app.BlockchainConfig

	// API/server
	APIProto           string `long:"apiproto" description:"Protocol for API (http or https)" env:"PDANALYTICS_ENABLE_HTTPS"`
	APIListen          string `long:"apilisten" description:"Listen address for API. default localhost:7777, :17778 testnet, :17779 simnet" env:"PDANALYTICS_LISTEN_URL"`
	ServerHeader       string `long:"server-http-header" description:"Set the HTTP response header Server key value. Valid values are \"off\", \"version\", or a custom string."`
	CacheControlMaxAge int    `long:"cachecontrol-maxage" description:"Set CacheControl in the HTTP response header to a value in seconds for clients to cache the response. This applies only to FileServer routes." env:"DCRDATA_MAX_CACHE_AGE"`
}

func defaultConfig() config {
	cfg := config{
		HomeDir:            defaultHomeDir,
		DataDir:            defaultDataDir,
		LogDir:             defaultLogDir,
		DBHost:             defaultDbHost,
		DBPort:             defaultDbPort,
		DBUser:             defaultDbUser,
		DBPass:             defaultDbPass,
		DBName:             defaultDbName,
		MaxLogZips:         defaultMaxLogZips,
		ConfigFile:         defaultConfigFile,
		DebugLevel:         defaultLogLevel,
		HTTPProfPath:       defaultHTTPProfPath,
		APIProto:           defaultAPIProto,
		CacheControlMaxAge: defaultCacheControlMaxAge,
		ServerHeader:       defaultServerHeader,
	}

	return cfg
}

// cleanAndExpandPath expands environment variables and leading ~ in the passed
// path, cleans the result, and returns it.
func cleanAndExpandPath(path string) string {
	// NOTE: The os.ExpandEnv doesn't work with Windows cmd.exe-style
	// %VARIABLE%, but the variables can still be expanded via POSIX-style
	// $VARIABLE.
	path = os.ExpandEnv(path)

	if !strings.HasPrefix(path, "~") {
		return filepath.Clean(path)
	}

	// Expand initial ~ to the current user's home directory, or ~otheruser to
	// otheruser's home directory.  On Windows, both forward and backward
	// slashes can be used.
	path = path[1:]

	var pathSeparators string
	if runtime.GOOS == "windows" {
		pathSeparators = string(os.PathSeparator) + "/"
	} else {
		pathSeparators = string(os.PathSeparator)
	}

	userName := ""
	if i := strings.IndexAny(path, pathSeparators); i != -1 {
		userName = path[:i]
		path = path[i:]
	}

	homeDir := ""
	var u *user.User
	var err error
	if userName == "" {
		u, err = user.Current()
	} else {
		u, err = user.Lookup(userName)
	}
	if err == nil {
		homeDir = u.HomeDir
	}
	// Fallback to CWD if user lookup fails or user has no home directory.
	if homeDir == "" {
		homeDir = "."
	}

	return filepath.Join(homeDir, path)
}

// normalizeNetworkAddress checks for a valid local network address format and
// adds default host and port if not present. Invalidates addresses that include
// a protocol identifier.
func normalizeNetworkAddress(a, defaultHost, defaultPort string) (string, error) {
	if strings.Contains(a, "://") {
		return a, fmt.Errorf("Address %s contains a protocol identifier, which is not allowed", a)
	}
	if a == "" {
		return defaultHost + ":" + defaultPort, nil
	}
	host, port, err := net.SplitHostPort(a)
	if err != nil {
		if strings.Contains(err.Error(), "missing port in address") {
			normalized := a + ":" + defaultPort
			host, port, err = net.SplitHostPort(normalized)
			if err != nil {
				return a, fmt.Errorf("Unable to address %s after port resolution: %v", normalized, err)
			}
		} else {
			return a, fmt.Errorf("Unable to normalize address %s: %v", a, err)
		}
	}
	if host == "" {
		host = defaultHost
	}
	if port == "" {
		port = defaultPort
	}
	return host + ":" + port, nil
}

// validLogLevel returns whether or not logLevel is a valid debug log level.
func validLogLevel(logLevel string) bool {
	_, ok := slog.LevelFromString(logLevel)
	return ok
}

// parseAndSetDebugLevels attempts to parse the specified debug level and set
// the levels accordingly.  An appropriate error is returned if anything is
// invalid.
func parseAndSetDebugLevels(debugLevel string) error {
	// When the specified string doesn't have any delimiters, treat it as
	// the log level for all subsystems.
	if !strings.Contains(debugLevel, ",") && !strings.Contains(debugLevel, "=") {
		// Validate debug log level.
		if !validLogLevel(debugLevel) {
			str := "The specified debug level [%v] is invalid"
			return fmt.Errorf(str, debugLevel)
		}

		// Change the logging level for all subsystems.
		setLogLevels(debugLevel)

		return nil
	}

	// Split the specified string into subsystem/level pairs while detecting
	// issues and update the log levels accordingly.
	for _, logLevelPair := range strings.Split(debugLevel, ",") {
		if !strings.Contains(logLevelPair, "=") {
			str := "The specified debug level contains an invalid " +
				"subsystem/level pair [%v]"
			return fmt.Errorf(str, logLevelPair)
		}

		// Extract the specified subsystem and log level.
		fields := strings.Split(logLevelPair, "=")
		subsysID, logLevel := fields[0], fields[1]

		// Validate log level.
		if !validLogLevel(logLevel) {
			str := "The specified debug level [%v] is invalid"
			return fmt.Errorf(str, logLevel)
		}

		setLogLevel(subsysID, logLevel)
	}

	return nil
}

func copyFile(sourec, destination string) error {
	from, err := os.Open(sourec)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(destination, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}

	return nil
}

// loadConfig initializes and parses the config using a config file and command
// line options.
func loadConfig() (*config, error) {
	loadConfigError := func(err error) (*config, error) {
		return nil, err
	}

	// Default config
	cfg := defaultConfig()
	defaultConfigNow := defaultConfig()

	// Load settings from environment variables.
	err := env.Parse(&cfg)
	if err != nil {
		return loadConfigError(err)
	}

	// If appdata was specified but not the config file, change the config file
	// path, and record this as the new default config file location.
	if defaultHomeDir != cfg.HomeDir && defaultConfigNow.ConfigFile == cfg.ConfigFile {
		cfg.ConfigFile = filepath.Join(cfg.HomeDir, defaultConfigFilename)
		// Update the defaultConfig to avoid an error if the config file in this
		// "new default" location does not exist.
		defaultConfigNow.ConfigFile = cfg.ConfigFile
	}

	// Pre-parse the command line options to see if an alternative config file
	// or the version flag was specified. Override any environment variables
	// with parsed command line flags.
	preCfg := cfg
	preParser := flags.NewParser(&preCfg, flags.HelpFlag|flags.PassDoubleDash)
	_, flagerr := preParser.Parse()

	if flagerr != nil {
		e, ok := flagerr.(*flags.Error)
		if !ok || e.Type != flags.ErrHelp {
			preParser.WriteHelp(os.Stderr)
		}
		if ok && e.Type == flags.ErrHelp {
			preParser.WriteHelp(os.Stdout)
			os.Exit(0)
		}
		return loadConfigError(flagerr)
	}

	// Show the version and exit if the version flag was specified.
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	if preCfg.ShowVersion {
		fmt.Printf("%s version 1.0 (Go version %s)\n", appName, runtime.Version())
		os.Exit(0)
	}

	// If a non-default appdata folder is specified on the command line, it may
	// be necessary adjust the config file location. If the the config file
	// location was not specified on the command line, the default location
	// should be under the non-default appdata directory. However, if the config
	// file was specified on the command line, it should be used regardless of
	// the appdata directory.
	if defaultHomeDir != preCfg.HomeDir && defaultConfigNow.ConfigFile == preCfg.ConfigFile {
		preCfg.ConfigFile = filepath.Join(preCfg.HomeDir, defaultConfigFilename)
		// Update the defaultConfig to avoid an error if the config file in this
		// "new default" location does not exist.
		defaultConfigNow.ConfigFile = preCfg.ConfigFile
	}

	// Config file name for logging.
	configFile := "NONE (defaults)"
	parser := flags.NewParser(&cfg, flags.Default)

	// Parse command line options again to ensure they take precedence.
	_, err = parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); !ok || e.Type != flags.ErrHelp {
			parser.WriteHelp(os.Stderr)
		}
		return loadConfigError(err)
	}

	// logRotator = nil

	// Initialize log rotation. After log rotation has been initialized, the
	// logger variables may be used. This creates the LogDir if needed.
	if cfg.MaxLogZips < 0 {
		cfg.MaxLogZips = 0
	}
	initLogRotator(filepath.Join(cfg.LogDir, defaultLogFilename), cfg.MaxLogZips)

	log.Infof("Log folder:  %s", cfg.LogDir)
	log.Infof("Config file: %s", configFile)

	// Output folder
	cfg.OutFolder = cleanAndExpandPath(cfg.OutFolder)

	// Ensure HTTP profiler is mounted with a valid path prefix.
	if cfg.HTTPProfile && (cfg.HTTPProfPath == "/" || len(defaultHTTPProfPath) == 0) {
		return loadConfigError(fmt.Errorf("httpprofprefix must not be \"\" or \"/\""))
	}

	// Parse, validate, and set debug log level(s).
	if cfg.Quiet {
		cfg.DebugLevel = "error"
	}

	// Parse, validate, and set debug log level(s).
	if err := parseAndSetDebugLevels(cfg.DebugLevel); err != nil {
		fmt.Fprintln(os.Stderr, err)
		parser.WriteHelp(os.Stderr)
		return loadConfigError(err)
	}

	port := defaultPort
	log.Info("Env $PORT :", os.Getenv("PORT"))
	if os.Getenv("PORT") != "" {
		_, err = strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			log.Critical(err)
			log.Critical("$PORT must be set")
		}
		port = os.Getenv("PORT")
	}
	// Check the supplied APIListen address
	if cfg.APIListen == "" {
		cfg.APIListen = defaultHost + ":" + port
	} else {
		cfg.APIListen, err = normalizeNetworkAddress(cfg.APIListen, defaultHost, port)
		if err != nil {
			return loadConfigError(err)
		}
	}

	switch cfg.ServerHeader {
	case "off":
		cfg.ServerHeader = ""
	case "version":
		cfg.ServerHeader = "coinzion - 1.0"
	}

	return &cfg, nil
}
