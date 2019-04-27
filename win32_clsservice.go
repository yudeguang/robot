package robot

import (
	//"C"
	"fmt"
	"syscall"
	"unsafe"
)

//API函数定义
var (
	mod_advapi32                   = syscall.NewLazyDLL("advapi32.dll")
	procOpenSCManager              = mod_advapi32.NewProc("OpenSCManagerW")
	procCloseServiceHandle         = mod_advapi32.NewProc("CloseServiceHandle")
	procOpenService                = mod_advapi32.NewProc("OpenServiceW")
	procCreateService              = mod_advapi32.NewProc("CreateServiceW")
	procChangeServiceConfig2       = mod_advapi32.NewProc("ChangeServiceConfig2W")
	procControlService             = mod_advapi32.NewProc("ControlService")
	procDeleteService              = mod_advapi32.NewProc("DeleteService")
	procRegisterServiceCtrlHandler = mod_advapi32.NewProc("RegisterServiceCtrlHandlerW")
	procSetServiceStatus           = mod_advapi32.NewProc("SetServiceStatus")
	procStartServiceCtrlDispatcher = mod_advapi32.NewProc("StartServiceCtrlDispatcherW")
	procQueryServiceConfig         = mod_advapi32.NewProc("QueryServiceConfigW")
	mod_kernel32                   = syscall.NewLazyDLL("kernel32.dll")
	procCreateFile                 = mod_kernel32.NewProc("CreateFileW")
	procGetModuleFileName          = mod_kernel32.NewProc("GetModuleFileNameW")
)

//服务配置信息
type ST_ServiceConfig struct {
	dwServiceType      uint32
	dwStartType        uint32
	dwErrorControl     uint32
	lpBinaryPathName   *uint16
	lpLoadOrderGroup   *uint16
	dwTagId            uint32
	lpDependencies     *uint16
	lpServiceStartName *uint16
	lpDisplayName      *uint16
	data               [4096]byte
}

//服务操作类
type ClsService struct {
	hSCManager uintptr
}

//打开服务控制台
func (this *ClsService) OpenSCManager() error {
	var err error
	this.hSCManager, _, err = procOpenSCManager.Call(uintptr(0),
		uintptr(0),
		uintptr(0x000f003f)) //SC_MANAGER_ALL_ACCESS
	if this.hSCManager == 0 {
		return err
	}
	return nil
}

//关闭服务控制台
func (this *ClsService) CloseSCManager() {
	if this.hSCManager != 0 {
		procCloseServiceHandle.Call(this.hSCManager)
		this.hSCManager = 0
	}
}

//判断服务是否已经安装
func (this *ClsService) IsExists(servicename string) (bool, error) {
	h, _, _ := procOpenService.Call(this.hSCManager,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(servicename))),
		uintptr(0x0001)) //SERVICE_QUERY_CONFIG
	if h == 0 {
		return false, nil
	}
	procCloseServiceHandle.Call(h)
	return true, nil
}

//创建服务
func (this *ClsService) CreateService(exepath, servicename, desc string, autostart bool, allowdesktop bool) error {
	serviceType, startType := 0x10, 3 //SERVICE_WIN32_OWN_PROCESS,SERVICE_DEMAND_START
	if allowdesktop {
		serviceType = serviceType | 0x100 //允许与桌面交互,SERVICE_INTERACTIVE_PROCESS
	}
	if autostart {
		startType = 2 //自动启动 SERVICE_AUTO_START
	}
	//安装服务
	h, _, err := procCreateService.Call(this.hSCManager, //hSCManager,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(servicename))), //lpServiceName
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(servicename))), //lpDisplayNam
		uintptr(0x000f01ff),                                            //dwDesiredAccess SERVICE_ALL_ACCESS
		uintptr(serviceType),                                           //dwServiceType SERVICE_WIN32_OWN_PROCESS
		uintptr(startType),                                             //dwStartType SERVICE_DEMAND_START or SERVICE_AUTO_START
		uintptr(0x1),                                                   //dwErrorControl SERVICE_ERROR_NORMAL
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(exepath))), //lpBinaryPathName
		uintptr(0x0), //lpLoadOrderGroup
		uintptr(0x0), //lpdwTagId
		uintptr(0x0), //lpDependencies
		uintptr(0x0), //lpServiceStartName
		uintptr(0x0)) //lpPassword
	if h == 0 {
		return err
	}
	defer procCloseServiceHandle.Call(h)
	//设置描述
	type ST_SERVICE_DESCRIPTIONW struct {
		lpDescription *uint16
	}
	pdesc := &ST_SERVICE_DESCRIPTIONW{}
	pdesc.lpDescription = syscall.StringToUTF16Ptr(desc)
	procChangeServiceConfig2.Call(h,
		uintptr(1), //SERVICE_CONFIG_DESCRIPTION
		uintptr(unsafe.Pointer(pdesc)))
	return nil
}

//删除服务
func (this *ClsService) DeleteService(servicename string) error {
	h, _, err := procOpenService.Call(this.hSCManager,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(servicename))),
		uintptr(0x00010020)) //SERVICE_STOP | DELETE
	if h == 0 {
		return err
	}
	defer procCloseServiceHandle.Call(h)
	//停止服务
	procControlService.Call(h,
		uintptr(0x1), //SERVICE_CONTROL_STOP
		uintptr(0x0)) //don't recv last status
	ret, _, err := procDeleteService.Call(h)
	if ret == 0 {
		return err
	}
	return nil
}
func utf16PtrToString(cstr *uint16) string {
	us := make([]uint16, 0)
	if cstr != nil {
		for p := uintptr(unsafe.Pointer(cstr)); ; p += 2 {
			u := *(*uint16)(unsafe.Pointer(p))
			if u == 0 {
				break
			}
			us = append(us, u)
		}
	}
	return syscall.UTF16ToString(us)
}

//查询服务信息
func QueryServiceBinaryPath(sername string) (string, error) {
	hServiceMgr, _, _ := procOpenSCManager.Call(uintptr(0), uintptr(0), uintptr(0xF003F))
	if uint32(hServiceMgr) == 0 {
		return "", fmt.Errorf("call OpenSCManager return false")
	}
	defer procCloseServiceHandle.Call(uintptr(hServiceMgr))
	hService, _, _ := procOpenService.Call(uintptr(hServiceMgr),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(sername))),
		uintptr(0xF003F))
	if uint32(hService) == 0 {
		return "", fmt.Errorf("call OpenService return false")
	}
	defer procCloseServiceHandle.Call(uintptr(hService))

	var nNeedSize uint32 = 0
	sconfig := ST_ServiceConfig{}
	ret, _, _ := procQueryServiceConfig.Call(uintptr(hService),
		uintptr(unsafe.Pointer(&sconfig)),
		uintptr(unsafe.Sizeof(sconfig)),
		uintptr(unsafe.Pointer(&nNeedSize)))
	if uint32(ret) == 0 {
		return "", fmt.Errorf("call QueryServiceConfig return false")
	}
	path := utf16PtrToString(sconfig.lpBinaryPathName)
	return path, nil
}
