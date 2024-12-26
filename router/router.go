package router
import (
    "fmt"
    "time"
    "regexp"
    "strconv"

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
    if config.ProxyMode {
        re := regexp.MustCompile(`:[0-9]+$`)
        data.AccessPoint = fmt.Sprintf("%s://%s.%s%s", config.BaseScheme, ins.ID, config.BaseHost, re.FindString(c.Request.Host))
        data.AccessPoint = fmt.Sprintf("<a href=\"%s\">%s</a>", data.AccessPoint, data.AccessPoint)
    } else if config.NCMode {
        data.AccessPoint = fmt.Sprintf("%s %s %d", config.BaseScheme, config.BaseHost, ins.Port)
        data.AccessPoint = fmt.Sprintf("<code>%s</code>", data.AccessPoint)
    } else {
        data.AccessPoint = fmt.Sprintf("%s://%s:%d", config.BaseScheme, config.BaseHost, ins.Port)
        data.AccessPoint = fmt.Sprintf("<a href=\"%s\">%s</a>", data.AccessPoint, data.AccessPoint)
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
