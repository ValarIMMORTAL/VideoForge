package processor

// 调用MoneyPrinterTurbo的/api/v1/videos（post） 及api/v1/tasks接口（get）
func GenerateVideo(videos VideoParams) {
	//todo 请求/api/v1/videos接口 将返回的taskid 传到api/v1/tasks中

	//启动协程，将轮询api/v1/tasks接口，当生成完毕时返回
}
