package param

type AddReq struct {
	X int64 `json:"x" description:"加数X" type:"int64"`
	Y int64 `json:"y" description:"加数Y" type:"int64"`
}

type AddRsp struct {
	Z int64 `json:"z" description:"结果Z" type:"int64"`
}
