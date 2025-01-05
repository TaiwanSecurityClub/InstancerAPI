package router
import (
    "fmt"
    "bytes"
    "time"
    "regexp"
    "strconv"
    "text/template"

    "github.com/gin-gonic/gin"

    "github.com/TaiwanSecurityClub/InstancerAPI/middlewares/token"
    "github.com/TaiwanSecurityClub/InstancerAPI/utils/config"
    "github.com/TaiwanSecurityClub/InstancerAPI/utils/errutil"
    "github.com/TaiwanSecurityClub/InstancerAPI/models/instance"
)

type statusdata struct {
    AccessPoint string `json:"accesspoint"`
    ExpiredAt time.Time `json:"expiredat"`
}

type userdata struct {
    ID int `json:"userid"`
}

var router *gin.RouterGroup

func Init(r *gin.RouterGroup) {
    router = r
    router.GET("/", token.CheckAuth, status)
    router.GET("/flag", token.CheckAuth, flag)
    router.POST("/create", token.CheckAuth, create)
    router.POST("/destroy", token.CheckAuth, destroy)
}

func status(c *gin.Context) {
    name := c.Query("userid")
    if name == "" {
        errutil.AbortAndStatus(c, 404)
        return
    }
    ins := instance.GetInstance(name)
    if ins == nil {
        errutil.AbortAndStatus(c, 404)
        return
    }
    data := statusdata{
        ExpiredAt: ins.ExpiredAt,
    }
    data.AccessPoint = ""
    re := regexp.MustCompile(`:[0-9]+$`)
    for i, port := range ins.Ports {
        switch config.GetMode(i) {
        case config.Forward:
            tmp := fmt.Sprintf("%s://%s:%d", config.BaseScheme, config.BaseHost, port)
            data.AccessPoint += fmt.Sprintf("<a href=\"%s\">%s</a><br/>", tmp, tmp)
        case config.Proxy:
            tmp := fmt.Sprintf("%s://%s%d.%s%s", config.BaseScheme, ins.ID, i, config.BaseHost, re.FindString(c.Request.Host))
            data.AccessPoint += fmt.Sprintf("<a href=\"%s\">%s</a><br/>", tmp, tmp)
        case config.Command:
            tmp, err := template.New("command").Parse(config.GetCommand(i))
            if err != nil {
                errutil.AbortAndStatus(c, 500)
                return
            }
            var buf bytes.Buffer
            err = tmp.Execute(&buf, struct {
                BaseHost string
                Port uint16
            } {
                BaseHost: config.BaseHost,
                Port: port,
            })
            if err != nil {
                errutil.AbortAndStatus(c, 500)
                return
            }
            data.AccessPoint += fmt.Sprintf("<code>%s</code><br/>", buf.String())
        }
    }
    c.JSON(200, data)
}

func flag(c *gin.Context) {
    name := c.Query("userid")
    if name == "" {
        errutil.AbortAndStatus(c, 404)
        return
    }
    ins := instance.GetInstance(name)
    if ins == nil {
        errutil.AbortAndStatus(c, 404)
        return
    }
    c.String(200, ins.GetFlag())
}

func create(c *gin.Context) {
    var user userdata
    if err := c.ShouldBindJSON(&user); err != nil {
        errutil.AbortAndStatus(c, 400)
    }
    name := strconv.Itoa(user.ID)
    _, err := instance.Up(name)
    if err != nil {
        panic(err)
    }
    c.JSON(200, true)
}

func destroy(c *gin.Context) {
    var user userdata
    if err := c.ShouldBindJSON(&user); err != nil {
        errutil.AbortAndStatus(c, 400)
    }
    name := strconv.Itoa(user.ID)
    err := instance.Down(name)
    if err != nil {
        panic(err)
    }
    c.JSON(200, true)
}
