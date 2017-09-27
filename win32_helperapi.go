/*
Package 创建WINDOWS服务的相关功能
*/
package robot

import (
	"C"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"unsafe"
)

//WIN API结构体定义
type ST_SERVICE_TABLE_ENTRY struct {
	lpServiceName *uint16
	lpEntryFun    uintptr
}

//服务控制结构体
type ST_SERVICE_STATUS struct {
	dwServiceType             uint32
	dwCurrentState            uint32
	dwControlsAccepted        uint32
	dwWin32ExitCode           uint32
	dwServiceSpecificExitCode uint32
	dwCheckPoint              uint32
	dwWaitHint                uint32
}

//错误信息定义
var (
	ErrParams    = errors.New("params error")
	ErrNotInfo   = errors.New("must call API SetServiceInfo set service info first")
	ErrNameError = errors.New("the service name error")
)

//服务信息结构体
type ST_ServiceInfo struct {
	ServiceName string //服务名
	ServiceDesc string //描述
	IsAutoStart bool   //是否自动启动
	IsDesktop   bool   //是否允许桌面交互
	PFN_OnStart func() //启动时调用的函数
	PFN_OnStop  func() //停止服务调用的函数
}

var gConfig *ST_ServiceInfo = nil

/*
判断本进程是否已运行
*/
func ThisIsRunning() bool {
	mLckFile := os.Args[0] + ".lck"
	//用独占打开方式并不关闭来判断
	ret, _, _ := procCreateFile.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(mLckFile))),
		uintptr(0x80000000|0x40000000),
		uintptr(0),
		uintptr(0),
		uintptr(4),
		uintptr(0),
		uintptr(0))
	if int(ret) == -1 {
		return true
	}
	return false
}

/*
获得当前EXE程序的路径
*/
func GetExeFileName() string {
	var nsize = uint32(1024)
	var buffer [1024]uint16
	r1, _, _ := procGetModuleFileName.Call(uintptr(0), uintptr(unsafe.Pointer(&buffer[0])), uintptr(nsize))
	if int(r1) <= 0 {
		return ""
	}
	path := string(syscall.UTF16ToString(buffer[:int(r1)]))
	return path
}

//以下定义常用的服务支持
//设置服务信息,安装、卸载、运行前需要先调用
func SetServiceInfo(p *ST_ServiceInfo) error {
	if p == nil {
		return ErrNotInfo
	}
	if p.ServiceName == "" || p.ServiceDesc == "" {
		return ErrNotInfo
	}
	//判断服务名是否违规
	reg := regexp.MustCompile("^([a-zA-Z0-9_]{1,100})$")
	if !reg.MatchString(p.ServiceName) {
		return ErrNameError
	}
	if p.PFN_OnStart == nil || p.PFN_OnStop == nil {
		return ErrNotInfo
	}
	gConfig = p
	return nil
}

/*
执行安装或卸载操作,自己处理命令行,
返回是否需要继续往下执行,是否错误
*/
func RunService() error {
	if gConfig == nil {
		return ErrNotInfo
	}
	//启动服务
	t := []ST_SERVICE_TABLE_ENTRY{
		{syscall.StringToUTF16Ptr(gConfig.ServiceName), syscall.NewCallback(serviceMain)},
		{nil, 0},
	}
	log.Println("开始运行服务......")
	ret, _, err := procStartServiceCtrlDispatcher.Call(uintptr(unsafe.Pointer(&t[0])))
	if ret == 0 {
		log.Println("运行服务失败:", err)
	} else {
		log.Println("停止服务")
	}
	return nil
}

//安装服务
func InstallService() error {
	if gConfig == nil {
		return ErrNotInfo
	}
	service := &ClsService{}
	err := service.OpenSCManager()
	if err != nil {
		return err
	}
	defer service.CloseSCManager()
	//获取路径
	exePath := os.Args[0]
	if !filepath.IsAbs(exePath) {
		var err error
		exePath, err = filepath.Abs(exePath)
		if err != nil {
			return err
		}
	}
	cmdLine := fmt.Sprintf("\"%s\" -svr:%s", exePath, gConfig.ServiceName)
	return service.CreateService(cmdLine, gConfig.ServiceName, gConfig.ServiceDesc, gConfig.IsAutoStart, gConfig.IsDesktop)
}

//卸载服务
func UninstallService() error {
	if gConfig == nil {
		return ErrNotInfo
	}
	service := &ClsService{}
	err := service.OpenSCManager()
	if err != nil {
		return err
	}
	defer service.CloseSCManager()
	return service.DeleteService(gConfig.ServiceName)
}

//判断是否运行在服务模式下
func IsRunAsService() bool {
	if gConfig == nil {
		panic(ErrNotInfo.Error())
	}
	svrArg := "-svr:" + gConfig.ServiceName
	if len(os.Args) > 1 && strings.ToLower(os.Args[1]) == strings.ToLower(svrArg) {
		return true
	}
	return false
}

//主工作服务
var gServiceStatus *ST_SERVICE_STATUS = nil
var chan_exit chan int
var gServiceHandle uintptr = 0

func serviceMain(argc uint32, argv **uint16) int {
	chan_exit = make(chan int)
	gServiceStatus = &ST_SERVICE_STATUS{dwServiceType: 0x30, dwCurrentState: 0x2, dwControlsAccepted: 0x1}
	res, _, err := procRegisterServiceCtrlHandler.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(gConfig.ServiceName))),
		syscall.NewCallback(eventControlHandle))
	if res == 0 {
		log.Println("注册服务回调失败:", err)
		return 0
	}
	gServiceHandle = res
	gServiceStatus.dwCurrentState = 0x4 //设置为启动,不管启动与否都继续往下走
	res, _, err = procSetServiceStatus.Call(gServiceHandle, uintptr(unsafe.Pointer(gServiceStatus)))
	if res == 0 {
		log.Println("设置服务为启动状态失败:", err)
	}
	//开始真正的任务处理
	gConfig.PFN_OnStart()
	<-chan_exit
	return 0
}

//控制事件回调
func eventControlHandle(ctllcode uint32) uintptr {
	switch ctllcode {
	case 0x1:
		gConfig.PFN_OnStop()
		gServiceStatus.dwWin32ExitCode = 0
		gServiceStatus.dwCurrentState = 1 //SERVICE_STOPPED
		gServiceStatus.dwCheckPoint = 0
		gServiceStatus.dwWaitHint = 0
		res, _, err := procSetServiceStatus.Call(gServiceHandle, uintptr(unsafe.Pointer(gServiceStatus)))
		if res == 0 {
			log.Println("设置服务停止状态失败:", err)
		}
		chan_exit <- 1
	case 0x4:
		break
	}
	return 0
}
