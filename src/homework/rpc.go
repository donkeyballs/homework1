package homework

const (
	PREPARE  = "PREPARE"
	COMMIT   = "COMMIT"
	ROLLBACK = "ROLLBACK"
	ACK      = "ACK"
	Checked  = "YES"
	DONE     = "DONE"
)

type NUM struct {
	index  int //更新的下标
	Number int //需要更新的数字
}

//命令
type OrderInfo struct {
	Stage string //PREPARE , COMMIT, ROLLBACK
	Args  [3]int //传递的参数
}

//回复
type Reply struct {
	Stage  string // ACK  PREPARED
	IndexI int
	Data   NUM
}
