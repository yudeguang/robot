package robot

import (
	"bytes"
	"fmt"
	"runtime"
	"syscall"
	"unsafe"
)

var (
	PAGE_SIZE                    = uint32(4096)
	START_ADDR                   = uint32(0x00010000)
	END_ADDR                     = uint32(0x7FFF0000)
	MEM_COMMIT                   = uint32(0x1000)
	PAGE_READWRITE               = uint32(0x04)
	PAGE_EXECUTE_READWRITE       = uint32(0x40)
	PAGE_WRITECOPY               = uint32(0x08)
	PAGE_EXECUTE_WRITECOPY       = uint32(0x80)
	modUser32                    = syscall.NewLazyDLL("user32.dll")
	modKernel32                  = syscall.NewLazyDLL("Kernel32.dll")
	procFindWindow               = modUser32.NewProc("FindWindowW")
	procGetWindowThreadProcessId = modUser32.NewProc("GetWindowThreadProcessId")
	procVirtualQueryEx           = modKernel32.NewProc("VirtualQueryEx")
	procReadProcessMemory        = modKernel32.NewProc("ReadProcessMemory")
)

func NewMemReader() *MemReader {
	p := &MemReader{}
	return p
}

//从内存中获取相关数据
func ReadMemoryFromExe(exeName string) ([]byte, error) {
	reader := NewMemReader()
	arrIds, err := GetProcessIdByName(exeName)
	if err != nil || len(arrIds) == 0 {
		return nil, err
	}
	err = reader.OpenById(uint32(arrIds[0]))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	buffer := bytes.NewBuffer(nil)
	err = reader.readMemory(buffer)
	if err != nil {
		return nil, err
	}
	data := buffer.Bytes()
	return data, nil
}

type tagMEMORY_BASIC_INFORMATION struct {
	BaseAddress       uint32
	AllocationBase    uint32
	AllocationProtect uint32
	RegionSize        uint32
	State             uint32
	Protect           uint32
	Type              uint32
}

var MEMORY_INFO_SIZE int = 28

type MemReader struct {
	hProcHandler syscall.Handle
}

//根据进程ID打开内存地址空间
func (m *MemReader) OpenById(procId uint32) error {
	m.Close()
	da := uint32(0x0400 | 0x0010 | 0x0020 | 0x0008)
	h, err := syscall.OpenProcess(da, false, procId)
	if err != nil {
		return err
	}
	m.hProcHandler = h
	return nil
}

//关闭进程
func (m *MemReader) Close() {
	if int(m.hProcHandler) != 0 {
		syscall.CloseHandle(m.hProcHandler)
		m.hProcHandler = 0
	}
}

//读取内存到一个bytes.Buffer中，
func (m *MemReader) readMemory(buffer *bytes.Buffer) error {
	buffer.Reset()
	runtime.GC()
	var tmpBuffer []byte = nil
	pMemInfo := &tagMEMORY_BASIC_INFORMATION{}
	beginAddr := START_ADDR
	for beginAddr < END_ADDR {
		err := m.queryMemoryInfomation(beginAddr, pMemInfo)
		if err != nil {
			return err
		}
		if m.isAccessible(pMemInfo) {
			if (pMemInfo.RegionSize % PAGE_SIZE) != 0 {
				return fmt.Errorf("内存的对齐方式错误")
			}
			if uint32(len(tmpBuffer)) < pMemInfo.RegionSize {
				tmpBuffer = make([]byte, pMemInfo.RegionSize)
			}
			//log.Printf("%.8X 可读,BASE:%.8X 长度:%d\r\n", beginAddr, pMemInfo.BaseAddress, pMemInfo.RegionSize)
			var readedSize uint32 = 0
			r, _, _ := procReadProcessMemory.Call(uintptr(m.hProcHandler),
				uintptr(pMemInfo.BaseAddress), uintptr(unsafe.Pointer(&tmpBuffer[0])),
				uintptr(pMemInfo.RegionSize), uintptr(unsafe.Pointer(&readedSize)))
			if int(r) == 1 {
				//全0的数字不要
				used := false
				for i := uint32(0); i < readedSize; i++ {
					if tmpBuffer[i] != 0x0 {
						used = true
						break
					}
				}
				if used {
					buffer.Write(tmpBuffer[:readedSize])
				}
			}
		}
		beginAddr = beginAddr + pMemInfo.RegionSize
	}
	return nil
}

//其它函数
func (m *MemReader) isAccessible(pMemInfo *tagMEMORY_BASIC_INFORMATION) bool {
	if pMemInfo.State == MEM_COMMIT {
		if ((pMemInfo.Protect & PAGE_READWRITE) != 0) ||
			((pMemInfo.Protect & PAGE_EXECUTE_READWRITE) != 0) {
			return true
		}
	}
	return false
}
func (m *MemReader) queryMemoryInfomation(baseAddr uint32, pMemInfo *tagMEMORY_BASIC_INFORMATION) error {
	r, _, err := procVirtualQueryEx.Call(uintptr(m.hProcHandler), uintptr(baseAddr),
		uintptr(unsafe.Pointer(pMemInfo)), uintptr(MEMORY_INFO_SIZE))
	if int(r) != MEMORY_INFO_SIZE {
		return err
	}
	return nil
}
