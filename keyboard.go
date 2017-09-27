package robot

import (
	"fmt"
	"syscall"
	"time"
)

var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")

	procDeleteFile = modkernel32.NewProc("DeleteFileW")
)

//定义按键编号
var (
	VK_SHIFT     = byte(0x10)
	VK_CTRL      = byte(0x11)
	VK_END       = byte(0x23)
	VK_HOME      = byte(0x24)
	VK_RETURN    = byte(0x0D)
	VK_CAPS_LOCK = byte(0x14) //CAPS LOCK 键 即大写键
	VK_MOUSE_R   = byte(0x02) //鼠标右键

	VK_DELETE              = byte(0x2E) //DELETE 键
	VK_MULTIPLICATION_SIGN = byte(0x6A) //(*) 键
	VK_PLUS_SIGN           = byte(0x6B) //(+) 键
	VK_ENTER               = byte(0x6C) //ENTER 键
	VK_MINUS_SIGN          = byte(0x6D) //(-) 键
	VK_DECIMAL_POINT       = byte(0x6E) //(.) 键
	VK_DIVISION_SIGN       = byte(0x6F) //(/) 键

	VK_a = byte(0x41)
	VK_b = byte(0x42)
	VK_c = byte(0x43)
	VK_d = byte(0x44)
	VK_e = byte(0x45)
	VK_f = byte(0x46)
	VK_g = byte(0x47)
	VK_h = byte(0x48)
	VK_i = byte(0x49)
	VK_j = byte(0x4A)
	VK_k = byte(0x4B)
	VK_l = byte(0x4C)
	VK_m = byte(0x4D)
	VK_n = byte(0x4E)
	VK_o = byte(0x4F)
	VK_p = byte(0x50)
	VK_q = byte(0x51)
	VK_r = byte(0x52)
	VK_s = byte(0x53)
	VK_t = byte(0x54)
	VK_u = byte(0x55)
	VK_v = byte(0x56)
	VK_w = byte(0x57)
	VK_x = byte(0x58)
	VK_y = byte(0x59)

	// VK_z = byte(0x59) Z不知道用什么表示

	VK_0 = byte(0x60)
	VK_1 = byte(0x61)
	VK_2 = byte(0x62)
	VK_3 = byte(0x63)
	VK_4 = byte(0x64)
	VK_5 = byte(0x65)
	VK_6 = byte(0x66)
	VK_7 = byte(0x67)
	VK_8 = byte(0x68)
	VK_9 = byte(0x69)
)
var (
	moduser32         = syscall.NewLazyDLL("User32.dll")
	procMouseEvent    = moduser32.NewProc("mouse_event")
	procKeyboardEvent = moduser32.NewProc("keybd_event")
)

//此函数仅处理0-9 a-z(含大写)以及-号（因为此程序字符只包含上述部分）
func KeyboardWrite(word string) (err error) {

	for _, v := range word {
		switch string(v) {
		case "-":
			Key(VK_MINUS_SIGN)
		case "a":
			Key(VK_a)
		case "b":
			Key(VK_b)
		case "c":
			Key(VK_c)
		case "d":
			Key(VK_d)
		case "e":
			Key(VK_e)
		case "f":
			Key(VK_f)
		case "g":
			Key(VK_g)
		case "h":
			Key(VK_h)
		case "i":
			Key(VK_i)
		case "j":
			Key(VK_j)
		case "k":
			Key(VK_k)
		case "l":
			Key(VK_l)
		case "m":
			Key(VK_m)
		case "n":
			Key(VK_n)
		case "o":
			Key(VK_o)
		case "p":
			Key(VK_p)
		case "q":
			Key(VK_q)
		case "r":
			Key(VK_r)
		case "s":
			Key(VK_s)
		case "t":
			Key(VK_t)
		case "u":
			Key(VK_u)
		case "v":
			Key(VK_v)
		case "w":
			Key(VK_w)
		case "x":
			Key(VK_x)
		case "y":
			Key(VK_y)
		case "z":
			KeyDown('Z')
			KeyUp('Z')

		case "A":

			Key(VK_a)

		case "B":

			Key(VK_b)

		case "C":

			Key(VK_c)

		case "D":

			Key(VK_d)

		case "E":

			Key(VK_e)

		case "F":

			Key(VK_f)

		case "G":

			Key(VK_g)

		case "H":

			Key(VK_h)

		case "I":

			Key(VK_i)

		case "J":

			Key(VK_j)

		case "K":

			Key(VK_k)

		case "L":

			Key(VK_l)

		case "M":

			Key(VK_m)

		case "N":

			Key(VK_n)

		case "O":

			Key(VK_o)

		case "P":

			Key(VK_p)

		case "Q":

			Key(VK_q)

		case "R":

			Key(VK_r)

		case "S":

			Key(VK_s)

		case "T":

			Key(VK_t)

		case "U":

			Key(VK_u)

		case "V":

			Key(VK_v)

		case "W":

			Key(VK_w)

		case "X":

			Key(VK_x)

		case "Y":

			Key(VK_y)

		case "Z":

			KeyDown('Z')
			KeyUp('Z')

		case "0":
			Key(VK_0)
		case "1":
			Key(VK_1)
		case "2":
			Key(VK_2)
		case "3":
			Key(VK_3)
		case "4":
			Key(VK_4)
		case "5":
			Key(VK_5)
		case "6":
			Key(VK_6)
		case "7":
			Key(VK_7)
		case "8":
			Key(VK_8)
		case "9":
			Key(VK_9)
		default:
			return fmt.Errorf("转换失败，遇到未定义字符")

		}
	}
	return nil
}

//同时按下键以及释放键两个动作
func Key(key byte) {
	KeyDown(key)
	KeyUp(key)
}

func CtrlA() { //按下 ctrl+a 全选
	KeyDown(VK_CTRL)
	KeyDown('A')
	KeyUp('A')
	KeyUp(VK_CTRL)
	time.Sleep(time.Millisecond * 100)
}

func CtrlC() { //按下 ctrl+C 复制
	KeyDown(VK_CTRL)
	KeyDown('C')
	KeyUp('C')
	KeyUp(VK_CTRL)
	time.Sleep(time.Millisecond * 100)
}

func CtrlS() { //按下 ctrl+S 保存
	KeyDown(VK_CTRL)
	KeyDown('S')
	KeyUp('S')
	KeyUp(VK_CTRL)
	time.Sleep(time.Millisecond * 100)
}

func CtrlV() { //按下 ctrl+V 粘贴
	KeyDown(VK_CTRL)
	KeyDown('V')
	KeyUp('V')
	KeyUp(VK_CTRL)
	time.Sleep(time.Millisecond * 100)
}
func CtrlX() { //按下 ctrl+X 剪切
	KeyDown(VK_CTRL)
	KeyDown('X')
	KeyUp('X')
	KeyUp(VK_CTRL)
	time.Sleep(time.Millisecond * 100)
}

//按下键
func KeyDown(key byte) {
	procKeyboardEvent.Call(uintptr(key), uintptr(0), uintptr(0), 0)
}

//释放键
func KeyUp(key byte) {
	procKeyboardEvent.Call(uintptr(key), uintptr(0), uintptr(2), 0)
}

/*
移动鼠标移动到x,y坐标,hasclick表明是否需要点击一下
x,y按屏幕分辨率计算(分辨率，直接查询系统分辨率)
*/

func MoveMouse(x, y int, hasclick bool) {
	nRealX := int(x * 65535 / 1440)
	nRealY := int(y * 65535 / 900)
	procMouseEvent.Call(uintptr(0x0001|0x8000), uintptr(nRealX), uintptr(nRealY), 0, 0)
	if hasclick {
		procMouseEvent.Call(uintptr(0x0002), uintptr(nRealX), uintptr(nRealY), 0, 0)
		procMouseEvent.Call(uintptr(0x0004), uintptr(nRealX), uintptr(nRealY), 0, 0)
	}
	//都延时1秒
	time.Sleep(10 * time.Millisecond)
}
func MoveMouseDoubleClick(x, y int, hasclick bool) {
	nRealX := int(x * 65535 / 1440)
	nRealY := int(y * 65535 / 900)
	procMouseEvent.Call(uintptr(0x0001|0x8000), uintptr(nRealX), uintptr(nRealY), 0, 0)
	if hasclick {
		procMouseEvent.Call(uintptr(0x0002), uintptr(nRealX), uintptr(nRealY), 0, 0)
		procMouseEvent.Call(uintptr(0x0004), uintptr(nRealX), uintptr(nRealY), 0, 0)
		time.Sleep(100 * time.Millisecond)
		procMouseEvent.Call(uintptr(0x0002), uintptr(nRealX), uintptr(nRealY), 0, 0)
		procMouseEvent.Call(uintptr(0x0004), uintptr(nRealX), uintptr(nRealY), 0, 0)
	}
	//都延时1秒
	time.Sleep(1 * time.Second)
}

//鼠标移动到某个位置，然后点击右键
func MoveMouseRitht(x, y int) {
	nRealX := int(x * 65535 / 1440)
	nRealY := int(y * 65535 / 900)
	procMouseEvent.Call(uintptr(0x0001|0x8000), uintptr(nRealX), uintptr(nRealY), 0, 0)

	procMouseEvent.Call(uintptr(0x0008), uintptr(nRealX), uintptr(nRealY), 0, 0)
	procMouseEvent.Call(uintptr(0x0010), uintptr(nRealX), uintptr(nRealY), 0, 0)

	//都延时1秒
	time.Sleep(10 * time.Millisecond)
}

//同时按住CTRL LeftMouse
func CtrlLeftMouse(x, y int, hasclick bool) {
	KeyDown(VK_CTRL)
	MoveMouse(x, y, hasclick)
	KeyUp(VK_CTRL)

}
