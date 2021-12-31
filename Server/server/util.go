package server

import (
	"fmt"
	"os"
)

/**
 * 错误显示并退出
 * @param err   错误类型数据
 * @param info  错误信息提示内容
 **/
func errorExit(err error, info string) {

	if err != nil {

		fmt.Println(info + "  " + err.Error())

		os.Exit(1)

	}

}

/**
 * 错误检查
 * @param err   错误类型数据
 * @param info  错误信息提示内容
 **/
func checkError(err error, info string) bool {

	if err != nil {

		fmt.Println(info + "  " + err.Error())

		return false
	}

	return true
}
