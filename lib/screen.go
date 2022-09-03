/**
 * @Author Kokutas
 * @Description //TODO
 * @Date 2021/2/10 23:37
 **/
package lib

import "github.com/kbinani/screenshot"
// 屏幕信息相关
// github.com/kbinani/screenshot 这个包可以实现屏幕录制

func ScreenSize()(width,height int){
	for i := 0; i < screenshot.NumActiveDisplays(); i++ {
		bounds :=screenshot.GetDisplayBounds(i)
		height = bounds.Max.X
		width = bounds.Max.Y
	}
	return width,height
}
