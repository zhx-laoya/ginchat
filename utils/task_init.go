package utils

import "time"

type TimerFunc func(interface{}) bool

/*
delay 首次延迟: 从调用函数开始到首次执行 fun 方法的时间间隔。
tick 间隔 : 每次执行完 fun 方法后，等待的时间间隔
fun 定时执行的方法
param 方法的参数
*/
// 实现了一个定时执行 fun 方法的功能。在每次执行完 fun 方法后
// ，会等待设定的时间间隔，然后再次触发执行，实现周期性的定时操作。
func Timer(delay, tick time.Duration, fun TimerFunc, param interface{}) {
	go func() {
		if fun == nil {
			return
		}
		// 创建了一个 time.Timer 对象 t，并设置其首次触发的时间为 delay。
		// 然后进入一个无限循环，使用 select 语句监听 t.C 通道的触发事件。
		t := time.NewTimer(delay)
		for {
			select {
			case <-t.C:
				if !fun(param) {
					return
				}
				t.Reset(tick)
			}
		}
	}()
}
