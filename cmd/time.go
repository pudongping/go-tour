package cmd

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/pudongping/go-tour/internal/timer"
)

var calculateTime string
var duration string

var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "时间格式处理",
	Long:  "时间格式处理",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var nowTimeCmd = &cobra.Command{
	Use:   "now",
	Short: "获取当前时间",
	Long:  "获取当前时间",
	Run: func(cmd *cobra.Command, args []string) {
		nowTime := timer.GetNowTime()
		// 2006-01-02 15:04:05 时间格式是约定的进行时间格式化的标准格式
		// 等同于 time.Now().Format(time.RFC3339)
		log.Printf("输出结果: %s, %d", nowTime.Format("2006-01-02 15:04:05"), nowTime.Unix())
	},
}

var calculateTimeCmd = &cobra.Command{
	Use:   "calc",
	Short: "计算所需时间",
	Long:  "计算所需时间",
	Run: func(cmd *cobra.Command, args []string) {
		var currentTimer time.Time
		var layout = "2006-01-02 15:04:05"
		if calculateTime == "" {
			currentTimer = timer.GetNowTime()
		} else {
			var err error
			space := strings.Count(calculateTime, " ")
			space1 := strings.Count(calculateTime, "-")
			space2 := strings.Count(calculateTime, ":")
			if space == 0 {
				layout = "2006-01-02"
			}
			if space == 1 || (space1 == 0 && space2 == 0) {
				layout = "2006-01-02 15:04:05"
			}
			calculateTime = strings.Trim(calculateTime, " ")
			location, _ := time.LoadLocation("Asia/Shanghai")
			// Parse 方法会尝试在入参的参数中中分析并读取时区信息，但是如果入参的参数没有指定时区信息的话，那么就会默认使用 UTC 时间
			// 因此我们最好使用 ParseInLocation 方法并指定时区
			currentTimer, err = time.ParseInLocation(layout, calculateTime, location)
			if err != nil {
				t, _ := strconv.Atoi(calculateTime)
				currentTimer = time.Unix(int64(t), 0)
			}
		}
		t, err := timer.GetCalculateTime(currentTimer, duration)
		if err != nil {
			log.Fatalf("timer.GetCalculateTime err: %v", err)
		}

		log.Printf("输出结果: %s, %d", t.Format(layout), t.Unix())
	},
}

func init() {
	timeCmd.AddCommand(nowTimeCmd)
	timeCmd.AddCommand(calculateTimeCmd)

	calculateTimeCmd.Flags().StringVarP(&calculateTime, "calculate", "c", "", ` 需要计算的时间，有效单位为时间戳或已格式化后的时间 `)
	calculateTimeCmd.Flags().StringVarP(&duration, "duration", "d", "", ` 持续时间，有效时间单位为"ns", "us" (or "µ s"), "ms", "s", "m", "h"`)
}
