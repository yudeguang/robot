package robot

import (
	"errors"
	"log"
	"strings"
	"syscall"
	"unsafe"
)

var (
	modpsapi                     = syscall.NewLazyDLL("psapi.dll")
	procOpenProcess              = mod_kernel32.NewProc("OpenProcess")
	procTerminateProcess         = mod_kernel32.NewProc("TerminateProcess")
	procCloseHandle              = mod_kernel32.NewProc("CloseHandle")
	procCreateToolhelp32Snapshot = mod_kernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32First           = mod_kernel32.NewProc("Process32FirstW")
	procProcess32Next            = mod_kernel32.NewProc("Process32NextW")
	procGetModuleFileNameEx      = modpsapi.NewProc("GetModuleFileNameExW")
)

type ProcInfo struct {
	PId   uint32 //进程Id
	PPId  uint32 //父进程Id
	PName string //进程名
	PPath string //进程全路径
}
type tagPROCESSENTRY32 struct {
	dwSize              uint32
	cntUsage            uint32
	th32ProcessID       uint32
	th32DefaultHeapID   uintptr
	th32ModuleID        uint32
	cntThreads          uint32
	th32ParentProcessID uint32
	pcPriClassBase      uint32
	dwFlags             uint32
	szExeFile           [260]uint16
}

/*
根据进程Id结束进程,如果权限不足需要先调用SetDebugPrivilege提升权限
*/
func KillProcessById(uProcId uint32) error {
	pHandle, _, err := procOpenProcess.Call(
		uintptr(0x0400|0x0001),
		uintptr(0),
		uintptr(uProcId))
	if int(pHandle) == 0 {
		return errors.New("OpenProcess:" + err.Error())
	}
	defer func() { //close handle on exit
		procCloseHandle.Call(uintptr(pHandle))
	}()
	ret, _, err := procTerminateProcess.Call(
		uintptr(pHandle),
		uintptr(1))
	if int(ret) == 0 {
		return errors.New("TerminateProcess:" + err.Error())
	}
	return nil
}

/*
根据进程名结束进程,如果权限不足需要先调用SetDebugPrivilege提升权限
*/
func KillEXEByName(name string) error {
	pidarr, err := GetProcessIdByName(name)
	if err != nil {
		return err
	}
	for _, id := range pidarr {
		err = KillProcessById(id)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
获取进程名为指定进程名的进程Id
*/
func GetProcessIdByName(name string) ([]uint32, error) {
	psarr, err := EnumProcess() //先枚举进程
	if err != nil {
		return nil, err
	}
	pidArr := make([]uint32, 0)
	name = strings.ToLower(name)
	for _, p := range psarr {
		v := strings.ToLower(p.PName)
		if v == name {
			pidArr = append(pidArr, p.PId)
		}
	}
	return pidArr, nil
}

//判断某程序是否已经运行，比如 qq.exe
func FindEXE(name string) bool {
	processlist, err := EnumProcess()
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range processlist {
		if strings.ToLower(v.PName) == strings.ToLower(name) {
			return true
		}
	}
	return false
}

/*
枚举运行的进程
*/
func EnumProcess() ([]ProcInfo, error) {
	//打开进程快照
	pHandle, _, err := procCreateToolhelp32Snapshot.Call(uintptr(0x2), uintptr(0x0))
	if int(pHandle) == -1 {
		return nil, errors.New("CreateToolhelp32Snapshot:" + err.Error())
	}
	defer func() { //close handle on exit
		procCloseHandle.Call(uintptr(pHandle))
	}()

	//枚举进程
	psArray := make([]ProcInfo, 0)
	proc := tagPROCESSENTRY32{}
	proc.dwSize = uint32(unsafe.Sizeof(proc))
	rt, _, err := procProcess32First.Call(uintptr(pHandle), uintptr(unsafe.Pointer(&proc)))
	if int(rt) != 1 {
		return nil, errors.New("Process32First:" + err.Error())
	}
	for {
		if int(rt) != 1 { //发生错误，有可能是最后一个进程了
			if err != syscall.ERROR_NO_MORE_FILES {
				return nil, errors.New("Process32Next:" + err.Error())
			}
			break //进程循环完了
		}
		path, _ := GetProcessFullPath(proc.th32ProcessID)
		psinfo := ProcInfo{
			PName: syscall.UTF16ToString(proc.szExeFile[:]),
			PId:   uint32(proc.th32ProcessID),
			PPId:  uint32(proc.th32ParentProcessID),
			PPath: path,
		}
		psArray = append(psArray, psinfo)
		rt, _, err = procProcess32Next.Call(uintptr(pHandle), uintptr(unsafe.Pointer(&proc)))
	}
	return psArray, nil
}

/*
获得进程全路径
*/
func GetProcessFullPath(pid uint32) (string, error) {
	if procGetModuleFileNameEx.Find() != nil { //函数不存在认为成功
		return "", nil
	}
	var szFilePath [1024]uint16
	pHandle, _, err := procOpenProcess.Call(
		uintptr(0x0400|0x0010),
		uintptr(0),
		uintptr(pid))
	if int(pHandle) == 0 {
		return "", errors.New("OpenProcess:" + err.Error())
	}
	defer func() { //close handle on exit
		procCloseHandle.Call(uintptr(pHandle))
	}()
	ret, _, err := procGetModuleFileNameEx.Call(uintptr(pHandle),
		uintptr(0),
		uintptr(unsafe.Pointer(&szFilePath)),
		uintptr(1024))
	if int(ret) <= 0 {
		return "", errors.New("GetModuleFileNameEx:" + err.Error())
	}
	return syscall.UTF16ToString(szFilePath[:int(ret)]), nil
}

/*
创建进程,使用一个新的console
*/
func CreateConsoleProcess(name string, args []string, ishidden bool, iswaitexit bool) (uint32, error) {
	cmdline := "\"" + name + "\""
	if args != nil {
		for _, arg := range args {
			cmdline += " \"" + arg + "\""
		}
	}
	si := new(syscall.StartupInfo)
	si.Cb = uint32(unsafe.Sizeof(*si))
	si.Flags = syscall.STARTF_USESHOWWINDOW
	if ishidden {
		si.ShowWindow = syscall.SW_HIDE
	} else {
		si.ShowWindow = syscall.SW_SHOW
	}
	pi := new(syscall.ProcessInformation)
	flags := syscall.CREATE_UNICODE_ENVIRONMENT | 0x00000010 //CREATE_NEW_CONSOLE
	err := syscall.CreateProcess(nil, syscall.StringToUTF16Ptr(cmdline), nil, nil, false, uint32(flags), nil, nil, si, pi)
	if err != nil {
		return 0, err
	}
	if iswaitexit {
		syscall.WaitForSingleObject(pi.Process, syscall.INFINITE)
	}
	syscall.CloseHandle(syscall.Handle(pi.Thread))
	syscall.CloseHandle(syscall.Handle(pi.Process))
	return pi.ProcessId, nil
}
