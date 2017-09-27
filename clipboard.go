package robot

/*
#include <windows.h>
BOOL SetClipboard(const char* ptext)
{
	if(!OpenClipboard(NULL))
	{
		return FALSE;
	}
	EmptyClipboard();
	if(ptext != NULL)
	{
		HGLOBAL h = GlobalAlloc(GHND|GMEM_SHARE, strlen(ptext)+1);
		if(h)
		{
			strcpy((LPSTR)GlobalLock(h),ptext);
			GlobalUnlock(h);
			SetClipboardData(CF_TEXT, h);
		}
	}
	CloseClipboard();
	return TRUE;
}
char* GetClipboard()
{
	char* result = NULL;
	if(!IsClipboardFormatAvailable(CF_TEXT))
	{
		return result;
	}
	if(!OpenClipboard(NULL))
	{
		return result;
	}
	HGLOBAL hMem = GetClipboardData(CF_TEXT);
	if (hMem != NULL)
	{
		LPCTSTR lpStr = (LPCTSTR)GlobalLock(hMem);
		if (lpStr != NULL)
		{
			result = malloc(strlen(lpStr)+1);
			strcpy(result,lpStr);
			GlobalUnlock(hMem);
		}
	}
	CloseClipboard();
	return result;
}
*/
import "C"
import (
	"fmt"
	"github.com/axgle/mahonia"
	"unsafe"
)

//剪切板操作对象
type Clipboard struct {
}

//设置文本到剪切板,设置为nil则表示清空剪切板
func (this *Clipboard) NewClipboard(val interface{}) error {
	if val == nil {
		C.SetClipboard(nil)
	} else {
		enc := mahonia.NewEncoder("gbk")
		stext := enc.ConvertString(fmt.Sprint(val))
		psztext := C.CString(stext)
		defer C.free(unsafe.Pointer(psztext))
		C.SetClipboard(psztext)
	}
	return nil
}

//从剪切板中获取文本数据
func (this *Clipboard) GetClipboard() string {
	res := C.GetClipboard()
	if res == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(res))
	stext := C.GoString(res)
	enc := mahonia.NewDecoder("gbk")
	stext = enc.ConvertString(stext)
	return stext
}
