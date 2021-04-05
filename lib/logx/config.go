package logx

type LogConf struct {
	ServiceName         string `json:",optional"`                                    // 服务名称，可选
	Mode                string `json:",default=console,options=console|file|volume"` // 日志模式，console-命令行，file-文件，volume-输出至docker挂载的文件内
	TimeFormat          string `json:",optional"`                                    // 自定义日志时间格式
	Path                string `json:",default=logs"`                                // 日志存储路径，默认放在 logs 目录
	Level               string `json:",default=info,options=info|error|fatal"`       // 日志级别，默认信息级，可选信息级|错误级|严重级
	Compress            bool   `json:",optional"`                                    // 是否开启gzip压缩
	KeepDays            int    `json:",optional"`                                    // 日志保留天数
	StackCooldownMillis int    `json:",default=100"`                                 // 写日志时间间隔，默认100毫秒
}
