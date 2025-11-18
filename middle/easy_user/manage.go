package easy_user

import (
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/injoyai/base/maps"
	"xorm.io/xorm"
)

type Config struct {
	AdminUsername string        //管理员账号
	AdminPassword string        //管理员密码
	Valid         time.Duration //登录有效期
	SuperToken    string        //超级token,可选
	DBType        string        //数据库类型
	DSN           string        //数据源
	Multipoint    bool          //是否开启多点登录

	OnInvalid  func(info *LoginInfo)       //登录过期回调,可选
	OnPassword func(password string) error //校验密码,可选
	OnToken    func(user *User) string     //生成token,可选
}

func NewManage(cfg Config) (*Manage, error) {
	if len(cfg.AdminUsername) == 0 {
		cfg.AdminUsername = "admin"
	}
	if len(cfg.AdminPassword) == 0 {
		cfg.AdminPassword = "admin"
	}
	if len(cfg.DBType) == 0 && len(cfg.DSN) == 0 {
		cfg.DBType = "sqlite3"
		cfg.DSN = "./data/database/db.db"
	}
	if cfg.OnPassword == nil {
		cfg.OnPassword = func(password string) error {
			if len(password) < 6 {
				return errors.New("密码长度不能小于6")
			}
			if strings.Contains(password, " ") {
				return errors.New("密码不能使用空格")
			}
			return nil
		}
	}
	db, err := xorm.NewEngine(cfg.DBType, cfg.DSN)
	if err != nil {
		return nil, err
	}
	if err := db.Sync2(new(User), new(LoginInfo)); err != nil {
		return nil, err
	}

	//初始化超级管理员
	has, err := db.Where("Username=?", cfg.AdminUsername).Get(new(User))
	if err != nil {
		return nil, err
	} else if !has {
		if _, err := db.Insert(&User{
			Username: cfg.AdminUsername,
			Password: cfg.AdminPassword,
		}); err != nil {
			return nil, err
		}
	}

	//加载用户登录信息

	return &Manage{
		DB:   db,
		Cfg:  cfg,
		User: maps.NewGeneric[string, *User](),
		Info: maps.NewGeneric[string, *maps.Generic[string, *LoginInfo]](),
	}, nil
}

type Manage struct {
	DB   *xorm.Engine
	Cfg  Config
	User *maps.Generic[string, *User]
	Info *maps.Generic[string, *maps.Generic[string, *LoginInfo]]
}

func (this *Manage) Auth(token string) (*User, error) {

	// 超级token
	if len(this.Cfg.SuperToken) > 0 && token == this.Cfg.SuperToken {
		return this.GetUser(this.Cfg.AdminUsername, true)
	}

	// 校验token长度
	if len(token) == 0 {
		return nil, ErrAuthFail
	}

	// 解密token
	bytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, ErrAuthFail
	}

	// 提取token信息
	list := strings.Split(string(bytes), "_")
	if len(list) != 2 {
		return nil, ErrAuthFail
	}
	username := list[0]

	// 获取用户信息
	u, err := this.GetUser(username, true)
	if err != nil {
		return u, err
	}

	// 校验token
	ls, _ := this.Info.Get(username)
	if ls != nil {
		info, _ := ls.Get(token)
		if info == nil {
			return u, ErrAuthFail
		}
		if info.Valid < time.Now().Unix() {
			ls.Del(token)
			return u, ErrAuthFail
		}
	}

	return u, nil
}

func (this *Manage) GetList(index, size int) ([]*User, error) {
	data := []*User{}
	err := this.DB.Limit(size, index*size).Find(&data)
	return data, err
}

// GetUser 获取用户信息,优先从缓存读取,不存在则从数据读取,并更新到缓存
func (this *Manage) GetUser(username string, tryCache bool) (*User, error) {

	//尝试从缓存读取用户
	if tryCache {
		if u, _ := this.User.Get(username); u != nil {
			return u, nil
		}
	}

	//从数据库读取用户
	user := new(User)
	has, err := this.DB.Where("Username=?", username).Get(user)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrNotFind
	}

	//添加到缓存
	this.User.Set(user.Username, user)

	return user, nil
}

func (this *Manage) CreateUser(req *UserCreateRequest) error {
	u, err := req.New(this.Cfg.OnPassword)
	if err != nil {
		return err
	}
	_, err = this.DB.Insert(u)
	if err != nil {
		return err
	}
	this.User.Set(u.Username, u)
	return nil
}

func (this *Manage) DelUser(username string) error {
	if username == this.Cfg.AdminUsername {
		return errors.New("无法删除超级管理员")
	}
	_, err := this.DB.Where("Username=?", username).Delete(&User{})
	if err != nil {
		return err
	}
	this.User.Del(username)
	return nil
}

// EditSelfInfo 修改个人信息
func (this *Manage) EditSelfInfo(u *User, req *EditInfoRequest) error {
	cols := u.EditInfo(req)
	_, err := this.DB.Where("Username=?", u.Username).Cols(cols).Update(u)
	return err
}

func (this *Manage) EditPassword(u *User, req *EditPasswordRequest) error {
	cols, err := u.EditPassword(req, this.Cfg.OnPassword)
	if err != nil {
		return err
	}
	if _, err = this.DB.Where("Username=?", u.Username).Cols(cols).Update(u); err != nil {
		return err
	}
	//清除token
	ls, _ := this.Info.Get(u.Username)
	if ls != nil {
		ls.Range(func(key string, value *LoginInfo) bool {
			value.Logout()
			return true
		})
	}
	return nil
}

// Login 登录,token有效,则修改token有效时长,无效则重新生成
func (this *Manage) Login(request *LoginRequest, device, addr string) (*User, error) {
	if err := request.Check(); err != nil {
		return nil, err
	}
	user, err := this.GetUser(request.Username, true)
	if err != nil && err != ErrNotFind {
		return nil, err
	} else if err != nil {
		return nil, ErrInput
	}

	//登录操作
	info, cols, err := user.Login(request, device, addr, this.Cfg.OnToken)
	if err != nil {
		return nil, err
	}

	//是否是多点登录
	if !this.Cfg.Multipoint {
		//全部登出操作
		ls, _ := this.Info.Get(user.Username)
		if ls != nil {
			ls.Range(func(key string, value *LoginInfo) bool {
				value.Logout()
				return true
			})
		}
	}

	//更新到数据库
	if _, err := this.DB.Where("Username=?", user.Username).Cols(cols).Update(user); err != nil {
		return nil, err
	}

	//添加登录信息
	if _, err := this.DB.Insert(info); err != nil {
		return nil, err
	}

	//添加到缓存
	ls, _ := this.Info.GetOrSetByHandler(user.Username, func() (*maps.Generic[string, *LoginInfo], error) {
		return maps.NewGeneric[string, *LoginInfo](), nil
	})
	ls.Set(info.Token, info)

	//登录过期事件
	info.OnInvalid(this.Cfg.OnInvalid, func(user *LoginInfo) {
		//登录过期删除登录信息
		if _, err := this.DB.Where("ID=?", user.ID).Delete(&LoginInfo{}); err == nil {
			if ls, _ := this.Info.Get(user.Username); ls != nil {
				ls.Del(user.Token)
			}
		}
	})

	return user, nil
}
