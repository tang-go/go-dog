package snowflake

//  雪花算法

// 算法模式如下：
// |  1   |    41     |    10    |    12    |
// | sign | timestamp | workerId | sequence |
import (
	"reflect"
	"sync"
	"time"
)

var (
	machineID     int64 // 机器 id 占10位, 十进制范围是 [ 0, 1023 ]
	sn            int64 // 序列号占 12 位,十进制范围是 [ 0, 4095 ]
	lastTimeStamp int64 // 上次的时间戳(毫秒级), 1秒=1000毫秒, 1毫秒=1000微秒,1微秒=1000纳秒
)

func init() {
	lastTimeStamp = time.Now().UnixNano() / 1000000
}

func _SetMachineID(mid int64) {
	// 把机器 id 左移 12 位,让出 12 位空间给序列号使用
	machineID = mid << 12
}

func _GetSnowflakeID() int64 {
	curTimeStamp := time.Now().UnixNano() / 1000000
	// 同一毫秒
	if curTimeStamp == lastTimeStamp {
		sn++
		// 序列号占 12 位,十进制范围是 [ 0, 4095 ]
		if sn > 4095 {
			time.Sleep(time.Millisecond)
			curTimeStamp = time.Now().UnixNano() / 1000000
			lastTimeStamp = curTimeStamp
			sn = 0
		}

		// 取 64 位的二进制数 0000000000 0000000000 0000000000 0001111111111 1111111111 1111111111  1 ( 这里共 41 个 1 )和时间戳进行并操作
		// 并结果( 右数 )第 42 位必然是 0,  低 41 位也就是时间戳的低 41 位
		rightBinValue := curTimeStamp & 0x1FFFFFFFFFF
		// 机器 id 占用10位空间,序列号占用12位空间,所以左移 22 位; 经过上面的并操作,左移后的第 1 位,必然是 0
		rightBinValue <<= 22
		id := rightBinValue | machineID | sn
		return id
	}
	if curTimeStamp > lastTimeStamp {
		sn = 0
		lastTimeStamp = curTimeStamp
		// 取 64 位的二进制数 0000000000 0000000000 0000000000 0001111111111 1111111111 1111111111  1 ( 这里共 41 个 1 )和时间戳进行并操作
		// 并结果( 右数 )第 42 位必然是 0,  低 41 位也就是时间戳的低 41 位
		rightBinValue := curTimeStamp & 0x1FFFFFFFFFF
		// 机器 id 占用10位空间,序列号占用12位空间,所以左移 22 位; 经过上面的并操作,左移后的第 1 位,必然是 0
		rightBinValue <<= 22
		id := rightBinValue | machineID | sn
		return id
	}
	if curTimeStamp < lastTimeStamp {
		return 0
	}
	return 0
}

//_Duplicate 去重
func _Duplicate(a interface{}) (ret []int64) {
	va := reflect.ValueOf(a)
	for i := 0; i < va.Len(); i++ {
		if i > 0 && reflect.DeepEqual(va.Index(i-1).Interface(), va.Index(i).Interface()) {
			continue
		}
		ret = append(ret, va.Index(i).Interface().(int64))
	}
	return ret
}

const (
	maxcount = 10000
)

//SnowFlake 雪花算法对象
type SnowFlake struct {
	max   int
	ids   []int64
	count int
	pos   int
	mutex sync.Mutex
}

//NewSnowFlake 新建一个对象
func NewSnowFlake(wID int64) *SnowFlake {
	_SetMachineID(wID)
	return &SnowFlake{
		max:   maxcount,
		count: 0,
		pos:   0,
	}
}

//GetID 获取ID
func (s *SnowFlake) GetID() int64 {
	s.mutex.Lock()
	if s.count <= 0 {
		for n := 0; n < s.max; n++ {
			s.ids = append(s.ids, _GetSnowflakeID())
		}
		s.ids = _Duplicate(s.ids)
		s.count = len(s.ids)
		s.pos = 0
	}
	id := s.ids[s.pos]
	s.count--
	s.pos++
	s.mutex.Unlock()
	return id
}
