package config

import (
    "os"
    "strconv"
    "time"

    netaddr "github.com/dspinhirne/netaddr-go"
)

var Debug bool
var Port string
var Token string
var MaxPort uint16
var MinPort uint16
var Validity time.Duration
var BaseScheme string
var BaseHost string
var FlagPrefix string
var FlagMsg string
var ChalDir string
var SubNetPool *netaddr.IPv4Net
var Prefix uint8
var DBservice string
var DBuser string
var DBpasswd string
var DBhost string
var DBport string
var DBname string
var DBdebug bool
var ProxyMode bool
var NCMode bool

func init() {
    loadenv()
    var err error
    debugstr, exists := os.LookupEnv("DEBUG")
    if !exists {
        Debug = false
    } else {
        Debug, err = strconv.ParseBool(debugstr)
        if err != nil {
            Debug = false
        }
    }
    dbdebugstr, exists := os.LookupEnv("DBDEBUG")
    if !exists {
        DBdebug = true
    } else {
        DBdebug, err = strconv.ParseBool(dbdebugstr)
        if err != nil {
            DBdebug = false
        }
    }
    proxymodestr, exists := os.LookupEnv("PROXYMODE")
    if !exists {
        ProxyMode = false
    } else {
        ProxyMode, err = strconv.ParseBool(proxymodestr)
        if err != nil {
            ProxyMode = false
        }
    }
    ncmodestr, exists := os.LookupEnv("NCMODE")
    if !exists || ProxyMode {
        NCMode = false
    } else {
        NCMode, err = strconv.ParseBool(ncmodestr)
        if err != nil {
            NCMode = false
        }
    }
    Port = os.Getenv("PORT")
    ChalDir = os.Getenv("CHALDIR")
    Token = os.Getenv("TOKEN")
    FlagPrefix = os.Getenv("FLAGPREFIX")
    FlagMsg = os.Getenv("FLAGMSG")
    BaseScheme = os.Getenv("BASESCHEME")
    BaseHost = os.Getenv("BASEHOST")
    DBservice, exists = os.LookupEnv("DBSERVICE")
    if !exists {
        DBservice = "sqlite"
    }
    DBuser = os.Getenv("DBUSER")
    DBpasswd = os.Getenv("DBPASSWD")
    DBhost = os.Getenv("DBHOST")
    DBport = os.Getenv("DBPORT")
    DBname = os.Getenv("DBNAME")
    subnetpoolstr, exists := os.LookupEnv("SUBNETPOOL")
    if !exists {
        SubNetPool, _ = netaddr.ParseIPv4Net("172.16.0.0/16")
    } else {
        SubNetPool, err = netaddr.ParseIPv4Net(subnetpoolstr)
        if err != nil {
            SubNetPool, _ = netaddr.ParseIPv4Net("172.16.0.0/16")
        }
    }
    prefixstr, exists := os.LookupEnv("PREFIX")
    if !exists {
        Prefix = 24
    } else {
        tmp, err := strconv.ParseUint(prefixstr, 10, 8)
        Prefix = uint8(tmp)
        if err != nil {
            Prefix = 24
        }
    }
    maxportstr, exists := os.LookupEnv("MAXPORT")
    if !exists {
        MaxPort = 30000
    } else {
        tmp, err := strconv.ParseUint(maxportstr, 10, 16)
        MaxPort = uint16(tmp)
        if err != nil {
            MaxPort = 30000
        }
    }
    minportstr, exists := os.LookupEnv("MINPORT")
    if !exists {
        MinPort = 30000
    } else {
        tmp, err := strconv.ParseUint(minportstr, 10, 16)
        MinPort = uint16(tmp)
        if err != nil {
            MinPort = 30000
        }
    }
    validitystr, exists := os.LookupEnv("VALIDITY")
    if !exists {
        Validity = 3 * time.Minute
    } else {
        Validity, err = time.ParseDuration(validitystr)
        if err != nil {
            Validity = 3 * time.Minute
        }
    }
}
