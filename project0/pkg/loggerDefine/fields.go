package loggerDefine

type Field struct {
	Key string
	Val any
}

// 包装通用方法，先穿关键值进来，然后操作组装一个关键key给回去
func Error(err error) Field {
	return Field{Key: "error", Val: err}
}


func Int64(key string,val int64) Field{
	return Field{
		Key: key,
		Val: val,
	}
}

func Int(key string,val int) Field{
	return Field{
		Key: key,
		Val: val,
	}
}


func String(key string,val string) Field{
	return Field{
		Key: key,
		Val: val,
	}
}