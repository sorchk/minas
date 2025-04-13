// package tool 提供各种工具函数，本文件实现了验证码相关功能
package tool

import (
	"encoding/json"
	"errors"
	"os"
	"server/utils/cache"
	"server/utils/config"
	"strconv"
	"strings"

	"github.com/wenlng/go-captcha/captcha" // 使用第三方验证码库
)

// GetCaptcha 生成并返回一个新的验证码
// 返回包含验证码图片、缩略图和密钥的Captcha结构体
func GetCaptcha() Captcha {
	capt := captcha.GetCaptcha()
	path, _ := os.Getwd()
	// ====================================================
	// Method: SetBackground(color []string);
	// Desc: 设置验证码背景图
	// ====================================================
	capt.SetBackground([]string{
		path + "/public/captcha/images/1.jpg",
		path + "/public/captcha/images/2.jpg",
		path + "/public/captcha/images/3.jpg",
		path + "/public/captcha/images/4.jpg",
		path + "/public/captcha/images/5.jpg",
	})
	// ====================================================
	// Method: SetFont(fonts []string);
	// Desc: 设置验证码字体
	// ====================================================
	capt.SetFont([]string{
		path + "/public/captcha/fonts/fzshengsksjw_cu.ttf",
		path + "/public/captcha/fonts/fzssksxl.ttf",
		path + "/public/captcha/fonts/hyrunyuan.ttf",
	})
	// ====================================================
	// Method: SetImageSize(size *Size);
	// Desc: 设置验证码主图的尺寸
	// ====================================================
	capt.SetImageSize(&captcha.Size{Width: 300, Height: 300})
	// ====================================================
	// Method: SetThumbSize(size *Size);
	// Desc: 设置验证码缩略图的尺寸
	// ====================================================
	capt.SetThumbSize(&captcha.Size{Width: 150, Height: 40})
	// ====================================================
	// Method: SetFontDPI(val int);
	// Desc: 设置验证码字体的随机DPI，最好是72
	// ====================================================
	capt.SetFontDPI(72)
	// ====================================================
	// Method: SetTextRangLen(val *captcha.RangeVal);
	// Desc: 设置验证码文本总和的随机长度范围
	// ====================================================
	capt.SetTextRangLen(&captcha.RangeVal{Min: 6, Max: 7})
	// ====================================================
	// Method: SetRangFontSize(val *captcha.RangeVal);
	// Desc: 设置验证码文本的随机大小
	// ====================================================
	capt.SetRangFontSize(&captcha.RangeVal{Min: 32, Max: 42})
	// ====================================================
	// Method: SetRangCheckTextLen(val *captcha.RangeVal);
	// Desc:设置验证码校验文本的随机长度的范围
	// ====================================================
	capt.SetRangCheckTextLen(&captcha.RangeVal{Min: 2, Max: 4})
	// ====================================================
	// Method: SetRangCheckFontSize(val *captcha.RangeVal);
	// Desc:设置验证码文本的随机大小
	// ====================================================
	capt.SetRangCheckFontSize(&captcha.RangeVal{Min: 24, Max: 30})
	// ====================================================
	// Method: SetTextRangFontColors(colors []string);
	// Desc: 设置验证码文本的随机十六进制颜色
	// ====================================================
	capt.SetTextRangFontColors([]string{
		"#1d3f84",
		"#3a6a1e",
		"#006600",
		"#005db9",
		"#aa002a",
		"#875400",
		"#6e3700",
		"#333333",
		"#660033",
	})
	// ====================================================
	// Method: SetImageFontAlpha(val float64);
	// Desc:设置验证码字体的透明度
	// ====================================================
	capt.SetImageFontAlpha(0.5)
	// ====================================================
	// Method: SetTextRangAnglePos(pos []*RangeVal);
	// Desc:设置验证码文本的角度
	// ====================================================
	capt.SetTextRangAnglePos([]*captcha.RangeVal{
		{Min: 1, Max: 15},
		{Min: 15, Max: 30},
		{Min: 30, Max: 45},
		{Min: 315, Max: 330},
		{Min: 330, Max: 345},
		{Min: 345, Max: 359},
	})
	// ====================================================
	// Method: SetImageFontDistort(val int);
	// Desc:设置验证码字体扭曲程度
	// ====================================================
	capt.SetImageFontDistort(captcha.ThumbBackgroundDistortLevel2)
	// ====================================================
	// Method: SetThumbBgColors(colors []string);
	// Desc: 设置缩略验证码背景的随机十六进制颜色
	// ====================================================
	capt.SetThumbBgColors([]string{
		"#1d3f84",
		"#3a6a1e",
	})
	// ====================================================
	// Method: SetThumbBackground(colors []string);
	// Desc:设置缩略验证码随机图像背景
	// ====================================================
	capt.SetThumbBackground([]string{
		path + "/public/captcha/images/thumb/r1.jpg",
		path + "/public/captcha/images/thumb/r2.jpg",
		path + "/public/captcha/images/thumb/r3.jpg",
		path + "/public/captcha/images/thumb/r4.jpg",
		path + "/public/captcha/images/thumb/r5.jpg",
	})
	// ====================================================
	// Method: SetThumbBgDistort(val int);
	// Desc:设置缩略验证码的扭曲程度
	// ====================================================
	capt.SetThumbBgDistort(captcha.ThumbBackgroundDistortLevel2)
	// ====================================================
	// Method: SetThumbBgCirclesNum(val int);
	// Desc:设置验证码背景的圈点数
	// ====================================================
	capt.SetThumbBgCirclesNum(20)
	// ====================================================
	// Method: SetThumbBgSlimLineNum(val int);
	// Desc:设置验证码背景的线条数
	// ====================================================
	capt.SetThumbBgSlimLineNum(3)

	// ====================================================
	// 生成验证码并处理结果
	dots, b64, tb64, key, err := capt.Generate()
	if err != nil {
		return Captcha{} // 生成失败时返回空对象
	}

	// 将点位置信息转为JSON并存入缓存系统
	dot, err := json.Marshal(dots)
	errSet := cache.GetCacheSystem().SetExpire(key, string(dot), config.CONF.Jwt.Expire)
	if err == nil && errSet == nil {
		// 成功生成验证码并缓存，返回完整信息
		return Captcha{
			ImageBase64: b64,  // 主验证码图片(Base64编码)
			ThumbBase64: tb64, // 缩略图(Base64编码)
			CaptchaKey:  key,  // 验证码唯一标识符
		}
	} else {
		// 缓存失败时返回空对象
		return Captcha{}
	}
}

// Check 验证用户提交的验证码是否正确
// 参数:
//   - key: 验证码唯一标识符
//   - dots: 用户点击位置的坐标，格式为"x1,y1,x2,y2..."
//
// 返回:
//   - string: 验证结果描述("check success"或"check fail")
//   - error: 验证失败时的错误信息
func Check(key string, dots string) (string, error) {
	// 检查参数是否为空
	if dots == "" && key == "" {
		return "check fail", errors.New("dots or key param is empty")
	}

	// 从缓存系统获取验证码信息
	cacheDots, err1 := cache.GetCacheSystem().Get(key)
	if err1 == nil && len(cacheDots) > 0 {
		// 解析用户提交的坐标点
		src := strings.Split(dots, ",")
		var dct map[int]captcha.CharDot
		if err2 := json.Unmarshal([]byte(cacheDots), &dct); err2 != nil {
			return "check fail", errors.New("illegal key")
		}

		// 验证用户点击位置是否正确
		chkRet := false
		if len(src) >= len(dct)*2 {
			chkRet = true
			for i, dot := range dct {
				j := i * 2   // X坐标在偶数位置
				k := i*2 + 1 // Y坐标在奇数位置
				a, _ := strconv.Atoi(src[j])
				b, _ := strconv.Atoi(src[k])
				// 判断点击位置是否在字符范围内
				chkRet = checkDist(a, b, dot.Dx, dot.Dy, dot.Width, dot.Height)
				if !chkRet {
					break
				}
			}
		}

		// 验证成功返回结果
		if chkRet && (len(dct)*2) == len(src) {
			return "check success", nil
		} else {
			return "check fail", errors.New("check fail")
		}
	} else {
		// 缓存中未找到验证码信息
		return "check fail", err1
	}
}

/**
 * @Description: 计算两点之间的距离，判断用户点击位置是否在指定的字符区域内
 * @param sx 用户点击的X坐标
 * @param sy 用户点击的Y坐标
 * @param dx 验证码文字的X坐标
 * @param dy 验证码文字的Y坐标
 * @param width 验证码文字的宽度
 * @param height 验证码文字的高度
 * @return bool 返回是否在有效区域内
 */
func checkDist(sx, sy, dx, dy, width int, height int) bool {
	// 判断点击坐标是否在字符的有效范围内
	return sx >= dx &&
		sx <= dx+width &&
		sy <= dy &&
		sy >= dy-height
}

// Captcha 验证码结构体，包含验证码的所有信息
type Captcha struct {
	ImageBase64 string `json:"image_base64"` // 主验证码图片(Base64编码)
	ThumbBase64 string `json:"thumb_base64"` // 缩略图(Base64编码)
	CaptchaKey  string `json:"captcha_key"`  // 验证码唯一标识符
}
