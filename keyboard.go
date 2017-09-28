package robot

import (
	"fmt"
	"strings"
	"syscall"
	"time"
)

var (
	modkernel32    = syscall.NewLazyDLL("kernel32.dll")
	procDeleteFile = modkernel32.NewProc("DeleteFileW")
)

//时间间隔
var sleepTime = time.Millisecond * 10

//把默认的分辨率写成1440X900
var screen_whith = 1440
var screen_height = 900

//设置目标电脑屏幕分辨率
func NewScreenSize(whith, height int) {
	if whith > 0 && height > 0 {
		screen_whith = whith
		screen_height = height
	}
}

//定义按键编号,仅列出常用的一些符号
var (

	// 左徽标键： VK_LWIN (91)
	// 右徽标键： VK_LWIN (92)
	// 鼠标右键快捷键：VK_APPS (93)
	// vK_MOUSE_R   = byte(0x02) //鼠标右键

	vK_BACK      = byte(0x08) //退格 backspace
	vK_TAB       = byte(0x09) //TAB
	vK_RETURN    = byte(0x0D) //回车
	vK_SHIFT     = byte(0x10) //shift
	vK_CTRL      = byte(0x11) //ctrl
	vK_ALT       = byte(0x12) //ALT
	vK_CAPS_LOCK = byte(0x14) //CAPS LOCK 键 即大写键
	vK_ESC       = byte(0x1B) //ESC
	vK_SPACE     = byte(0x20) //空格

	vK_MULTIPLICATION_SIGN = byte(0x6A) //(*) 键
	vK_PLUS_SIGN           = byte(0x6B) //(+) 键
	vK_ENTER               = byte(0x6C) //ENTER 键 应该是小键盘区域那个与vK_RETURN差别不大
	vK_MINUS_SIGN          = byte(0x6D) //(-) 键
	vK_DECIMAL_POINT       = byte(0x6E) //(.) 键
	vK_DIVISION_SIGN       = byte(0x6F) //(/) 键

	vK_HOME     = byte(0x24)
	vK_PageUp   = byte(0x21)
	vK_PageDown = byte(0x22)
	vK_END      = byte(0x23)
	vK_INSERT   = byte(0x2D)
	vK_DELETE   = byte(0x2E) //DELETE 键
	//方向键
	vK_LEFT  = byte(0x25) //方向键(←)：
	vK_UP    = byte(0x26) //方向键(↑)：
	vK_RIGHT = byte(0x27) //方向键(→)
	vK_DOWN  = byte(0x28) //方向键(↓)

	// F1到F12
	vK_F1  = byte(0x70)
	vK_F2  = byte(0x71)
	vK_F3  = byte(0x72)
	vK_F4  = byte(0x73)
	vK_F5  = byte(0x74)
	vK_F6  = byte(0x75)
	vK_F7  = byte(0x76)
	vK_F8  = byte(0x77)
	vK_F9  = byte(0x78)
	vK_F10 = byte(0x79)
	vK_F11 = byte(0x7A)
	vK_F12 = byte(0x7B)
	// a-z
	vK_a = byte(0x41)
	vK_b = byte(0x42)
	vK_c = byte(0x43)
	vK_d = byte(0x44)
	vK_e = byte(0x45)
	vK_f = byte(0x46)
	vK_g = byte(0x47)
	vK_h = byte(0x48)
	vK_i = byte(0x49)
	vK_j = byte(0x4A)
	vK_k = byte(0x4B)
	vK_l = byte(0x4C)
	vK_m = byte(0x4D)
	vK_n = byte(0x4E)
	vK_o = byte(0x4F)
	vK_p = byte(0x50)
	vK_q = byte(0x51)
	vK_r = byte(0x52)
	vK_s = byte(0x53)
	vK_t = byte(0x54)
	vK_u = byte(0x55)
	vK_v = byte(0x56)
	vK_w = byte(0x57)
	vK_x = byte(0x58)
	vK_y = byte(0x59)
	vK_z = byte(0x5A)
	// 0-9
	vK_0 = byte(0x60)
	vK_1 = byte(0x61)
	vK_2 = byte(0x62)
	vK_3 = byte(0x63)
	vK_4 = byte(0x64)
	vK_5 = byte(0x65)
	vK_6 = byte(0x66)
	vK_7 = byte(0x67)
	vK_8 = byte(0x68)
	vK_9 = byte(0x69)
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
			Key(vK_MINUS_SIGN)
		case "a":
			Key(vK_a)
		case "b":
			Key(vK_b)
		case "c":
			Key(vK_c)
		case "d":
			Key(vK_d)
		case "e":
			Key(vK_e)
		case "f":
			Key(vK_f)
		case "g":
			Key(vK_g)
		case "h":
			Key(vK_h)
		case "i":
			Key(vK_i)
		case "j":
			Key(vK_j)
		case "k":
			Key(vK_k)
		case "l":
			Key(vK_l)
		case "m":
			Key(vK_m)
		case "n":
			Key(vK_n)
		case "o":
			Key(vK_o)
		case "p":
			Key(vK_p)
		case "q":
			Key(vK_q)
		case "r":
			Key(vK_r)
		case "s":
			Key(vK_s)
		case "t":
			Key(vK_t)
		case "u":
			Key(vK_u)
		case "v":
			Key(vK_v)
		case "w":
			Key(vK_w)
		case "x":
			Key(vK_x)
		case "y":
			Key(vK_y)
		case "z":
			KeyDown('Z')
			KeyUp('Z')

		case "A":

			Key(vK_a)

		case "B":

			Key(vK_b)

		case "C":

			Key(vK_c)

		case "D":

			Key(vK_d)

		case "E":

			Key(vK_e)

		case "F":

			Key(vK_f)

		case "G":

			Key(vK_g)

		case "H":

			Key(vK_h)

		case "I":

			Key(vK_i)

		case "J":

			Key(vK_j)

		case "K":

			Key(vK_k)

		case "L":

			Key(vK_l)

		case "M":

			Key(vK_m)

		case "N":

			Key(vK_n)

		case "O":

			Key(vK_o)

		case "P":

			Key(vK_p)

		case "Q":

			Key(vK_q)

		case "R":

			Key(vK_r)

		case "S":

			Key(vK_s)

		case "T":

			Key(vK_t)

		case "U":

			Key(vK_u)

		case "V":

			Key(vK_v)

		case "W":

			Key(vK_w)

		case "X":

			Key(vK_x)

		case "Y":

			Key(vK_y)

		case "Z":

			KeyDown('Z')
			KeyUp('Z')

		case "0":
			Key(vK_0)
		case "1":
			Key(vK_1)
		case "2":
			Key(vK_2)
		case "3":
			Key(vK_3)
		case "4":
			Key(vK_4)
		case "5":
			Key(vK_5)
		case "6":
			Key(vK_6)
		case "7":
			Key(vK_7)
		case "8":
			Key(vK_8)
		case "9":
			Key(vK_9)
		default:
			return fmt.Errorf("转换失败，遇到未定义字符")

		}
	}
	return nil
}

//按下键
func KeyDown(key byte) {
	procKeyboardEvent.Call(uintptr(key), uintptr(0), uintptr(0), 0)
}

//释放键
func KeyUp(key byte) {
	procKeyboardEvent.Call(uintptr(key), uintptr(0), uintptr(2), 0)
}

//用文本的方式输入
func Press(key string) {
	//清洗容易出错的部分键值
	key = strings.ToLower(key)
	if key == "backspace" || key == "back space" {
		key = "back"
	} else if key == "table" {
		key = "tab"
	} else if key == "caps" || key == "caps lock" || key == "capslock" {
		key = "caps_lock"
	} else if key == "page_up" || key == "page up" {
		key = "pageup"
	} else if key == "page_down" || key == "page down" {
		key = "pagedown"
	}
	switch key {
	case "back":
		Key(vK_BACK)
	case "tab":
		Key(vK_TAB)
	case "return":
		Key(vK_RETURN)
	case "shift":
		Key(vK_SHIFT)
	case "ctrl":
		Key(vK_CTRL)
	case "alt":
		Key(vK_ALT)
	case "caps_lock":
		Key(vK_CAPS_LOCK)
	case "esc":
		Key(vK_ESC)
	case "space":
		Key(vK_SPACE)
	case "*":
		Key(vK_MULTIPLICATION_SIGN)
	case "+":
		Key(vK_PLUS_SIGN)
	case "enter":
		Key(vK_ENTER)
	case "-":
		Key(vK_MINUS_SIGN)
	case ".":
		Key(vK_DECIMAL_POINT)
	case `/`:
		Key(vK_DIVISION_SIGN)
	case `home`:
		Key(vK_HOME)
	case `pageup`:
		Key(vK_PageUp)
	case `pagedown`:
		Key(vK_PageDown)
	case `end`:
		Key(vK_END)
	case `insert`:
		Key(vK_INSERT)
	case `delete`:
		Key(vK_DELETE)
	case `left`:
		Key(vK_LEFT)
	case `up`:
		Key(vK_UP)
	case `right`:
		Key(vK_RIGHT)
	case `down`:
		Key(vK_DOWN)
	case `f1`:
		Key(vK_F1)
	case `f2`:
		Key(vK_F2)
	case `f3`:
		Key(vK_F3)
	case `f4`:
		Key(vK_F4)
	case `f5`:
		Key(vK_F5)
	case `f6`:
		Key(vK_F6)
	case `f7`:
		Key(vK_F7)
	case `f8`:
		Key(vK_F8)
	case `f9`:
		Key(vK_F9)
	case `f10`:
		Key(vK_F10)
	case `f11`:
		Key(vK_F11)
	case `f12`:
		Key(vK_F12)
	default:
		panic("未知输入，请改用Key()函数直接输入")
	}
}

//同时按下键以及释放键两个动作
func Key(key byte) {
	KeyDown(key)
	KeyUp(key)
}

func CtrlA() { //按下 ctrl+a 全选
	KeyDown(vK_CTRL)
	KeyDown('A')
	KeyUp('A')
	KeyUp(vK_CTRL)
	time.Sleep(sleepTime)
}

func CtrlC() { //按下 ctrl+C 复制
	KeyDown(vK_CTRL)
	KeyDown('C')
	KeyUp('C')
	KeyUp(vK_CTRL)
	time.Sleep(sleepTime)
}

func CtrlS() { //按下 ctrl+S 保存
	KeyDown(vK_CTRL)
	KeyDown('S')
	KeyUp('S')
	KeyUp(vK_CTRL)
	time.Sleep(sleepTime)
}

func CtrlV() { //按下 ctrl+V 粘贴
	KeyDown(vK_CTRL)
	KeyDown('V')
	KeyUp('V')
	KeyUp(vK_CTRL)
	time.Sleep(sleepTime)
}
func CtrlX() { //按下 ctrl+X 剪切
	KeyDown(vK_CTRL)
	KeyDown('X')
	KeyUp('X')
	KeyUp(vK_CTRL)
	time.Sleep(sleepTime)
}

//同时按住CTRL LeftMouse
func CtrlLeftMouseClick(x, y int, hasclick bool) {
	KeyDown(vK_CTRL)
	MoveMouseClick(x, y)
	KeyUp(vK_CTRL)
}

/*
移动鼠标移动到x,y坐标,hasclick表明是否需要点击一下
x,y按屏幕分辨率计算(分辨率，直接查询系统分辨率)
*/
func MoveMouse(X, Y int) {
	X = int(X * 65535 / screen_whith)
	Y = int(Y * 65535 / screen_height)
	procMouseEvent.Call(uintptr(0x0001|0x8000), uintptr(X), uintptr(Y), 0, 0)
	time.Sleep(sleepTime)
}

//鼠标移动到某个位置，默认是左侧鼠标点击
func MoveMouseClick(X, Y int) {
	X = int(X * 65535 / screen_whith)
	Y = int(Y * 65535 / screen_height)
	procMouseEvent.Call(uintptr(0x0001|0x8000), uintptr(X), uintptr(Y), 0, 0)
	procMouseEvent.Call(uintptr(0x0002), uintptr(X), uintptr(Y), 0, 0)
	procMouseEvent.Call(uintptr(0x0004), uintptr(X), uintptr(Y), 0, 0)
	time.Sleep(sleepTime)
}

//鼠标移动到某个位置，然后点击右键
func MoveMouseRithtClick(X, Y int) {
	X = int(X * 65535 / screen_whith)
	Y = int(Y * 65535 / screen_height)
	procMouseEvent.Call(uintptr(0x0001|0x8000), uintptr(X), uintptr(Y), 0, 0)
	procMouseEvent.Call(uintptr(0x0008), uintptr(X), uintptr(Y), 0, 0)
	procMouseEvent.Call(uintptr(0x0010), uintptr(X), uintptr(Y), 0, 0)
	time.Sleep(sleepTime)
}

func MoveMouseDoubleClick(X, Y int) {
	X = int(X * 65535 / screen_whith)
	Y = int(Y * 65535 / screen_height)

	procMouseEvent.Call(uintptr(0x0001|0x8000), uintptr(X), uintptr(Y), 0, 0)

	procMouseEvent.Call(uintptr(0x0002), uintptr(X), uintptr(Y), 0, 0)
	procMouseEvent.Call(uintptr(0x0004), uintptr(X), uintptr(Y), 0, 0)
	time.Sleep(sleepTime)

	procMouseEvent.Call(uintptr(0x0002), uintptr(X), uintptr(Y), 0, 0)
	procMouseEvent.Call(uintptr(0x0004), uintptr(X), uintptr(Y), 0, 0)
	time.Sleep(sleepTime)
}
