package easy_user

import (
	"context"
	"errors"
	"github.com/injoyai/base/crypt/md5"
	"github.com/injoyai/conv"
	"strings"
	"time"
)

var (
	ErrAuthFail = errors.New("验证失败")
	ErrInput    = errors.New("用户名或密码错误")
	ErrNotFind  = errors.New("用户不存在")
)

type User struct {
	ID          int64  `json:"id"`                    //主键
	Username    string `json:"username" xorm:"index"` //账号
	Password    string `json:"password"`              //密码
	LoginDevice string `json:"loginDevice"`           //登录设备
	LoginIP     string `json:"loginIP"`               //登录ip
	LoginNum    int    `json:"loginNum"`              //登录次数
	LoginTime   int64  `json:"loginTime"`             //登录时间
	LoginValid  int64  `json:"loginValid"`            //登录有效期,秒,0永久有效
	Photo       string `json:"photo"`                 //头像
}

// EditInfo 修改信息
func (this *User) EditInfo(request *EditInfoRequest) string {
	this.LoginValid = request.LoginValid
	this.Photo = request.Photo
	return "LoginValid,Photo"
}

// EditPassword 修改密码
func (this *User) EditPassword(req *EditPasswordRequest, onPassword func(password string) error) (string, error) {
	if len(req.NewPassword) == 0 {
		return "", errors.New("密码不能为空")
	}
	if this.Password != md5.Encrypt(req.OldPassword) {
		return "", errors.New("密码错误")
	}
	if onPassword != nil {
		if err := onPassword(req.NewPassword); err != nil {
			return "", err
		}
	}
	this.Password = md5.Encrypt(req.NewPassword)
	return "Password", nil
}

// Login 登录
func (this *User) Login(request *LoginRequest, device, addr string, onToken func(user *User) string) (*LoginInfo, string, error) {
	if this.Username != request.Username {
		return nil, "", ErrInput
	}
	//校验密码
	if this.Password != md5.Encrypt(request.Password) {
		return nil, "", ErrInput
	}
	this.LoginNum++
	this.LoginTime = time.Now().Unix()
	this.LoginDevice = device
	this.LoginIP = strings.Split(addr, ":")[0]
	info := this.NewLoginInfo(onToken)
	return info, "LoginNum,LoginDate,LoginIP", nil
}

// LoginRequest 登录结构
type LoginRequest struct {
	Username string `json:"username"` //账号
	Password string `json:"password"` //密码
}

func (this *LoginRequest) Check() error {
	if len(this.Username) == 0 && len(this.Password) == 0 {
		return errors.New("请输入用户名和密码")
	}
	if len(this.Username) == 0 {
		return errors.New("请输入用户名和密码")
	}
	if len(this.Password) == 0 {
		return errors.New("请输入用户名和密码")
	}
	return nil
}

// EditPasswordRequest 修改密码结构
type EditPasswordRequest struct {
	OldPassword string `json:"oldPassword"` //旧密码
	NewPassword string `json:"newPassword"` //新密码
}

// EditInfoRequest 修改信息结构
type EditInfoRequest struct {
	LoginValid int64  `json:"loginValid"` //登录有效期,秒,0永久有效
	Photo      string `json:"photo"`      //头像
}

// UserCreateRequest 新建用户结构
type UserCreateRequest struct {
	Username   string `json:"username"`   //账号
	Password   string `json:"password"`   //密码
	LoginValid int64  `json:"loginValid"` //登录有效期,秒,0永久有效
	Photo      string `json:"photo"`      //头像
}

func (this *UserCreateRequest) New(onPassword func(password string) error) (*User, error) {
	if len(this.Username) == 0 {
		return nil, errors.New("用户名不能为空")
	}
	if onPassword != nil {
		if err := onPassword(this.Password); err != nil {
			return nil, err
		}
	}
	return &User{
		Username:   this.Username,
		Password:   md5.Encrypt(this.Password),
		LoginValid: this.LoginValid,
		Photo:      this.Photo,
	}, nil
}

func (this *User) NewLoginInfo(onToken func(user *User) string) *LoginInfo {
	return &LoginInfo{
		Username: this.Username,
		Device:   this.LoginDevice,
		IP:       this.LoginIP,
		Time:     this.LoginTime,
		Token: func() string {
			if onToken != nil {
				return this.Username + "_" + onToken(this)
			}
			return this.Username + "_" + md5.Encrypt(conv.String(this.LoginTime))
		}(),
		Valid: conv.Select(this.LoginValid <= 0, 0, this.LoginTime+this.LoginValid),
	}
}

type LoginInfo struct {
	ID       int64              `json:"id"`
	Username string             `json:"username"`
	Device   string             `json:"device"`
	IP       string             `json:"ip"`
	Time     int64              `json:"time"`
	Token    string             `json:"token"`
	Valid    int64              `json:"valid"`
	cancel   context.CancelFunc //登出
}

func (this *LoginInfo) Logout() {
	if this.cancel != nil {
		this.cancel()
	}
}

// OnInvalid 开启协程监听过期
func (this *LoginInfo) OnInvalid(onInvalid ...func(user *LoginInfo)) {
	if this.cancel != nil {
		this.cancel()
	}
	if this.Valid <= 0 || len(onInvalid) == 0 {
		//永久有效
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	this.cancel = cancel
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
		case <-time.After(time.Second * time.Duration(this.Valid-(time.Now().Unix()-this.Time))):
		}
		for _, f := range onInvalid {
			f(this)
		}
	}(ctx)
}
