package redis

import "time"

type Lock struct {
	ops Ops
}

const (
	LockPrefix       = "dlock:"
	LockValue        = "DLOCK"
	DefaultLeaseTime = 30 * time.Second
)

func NewLock(ops Ops) *Lock {
	return &Lock{
		ops: ops,
	}
}

/*
*
获取一个锁，如果锁的key已经存在根据传入的第二个参数的时间进行等待，直到超时或者拿到锁为止。
切记业务完成后一定要及时释放锁。
*/
func (lock *Lock) GetDLock(key string, waitSecond int) (error, bool) {
	successGetLock, err := lock.ops.SetNX(LockPrefix+key, LockValue)
	if err != nil {
		return err, successGetLock
	}
	if !successGetLock {
		beginTime := time.Now()
		for {
			successGetLock, err = lock.ops.SetNXTtl(LockPrefix+key, LockValue, DefaultLeaseTime)
			if successGetLock {
				break
			}
			if err != nil {
				break
			}
			currentTime := time.Now()
			diffTime := currentTime.Sub(beginTime)
			if diffTime.Seconds() >= float64(waitSecond) {
				break
			}
		}
	}
	return err, successGetLock
}

/*
*
获取一个锁，如果锁的key已经存在根据传入的第二个参数的时间进行等待，直到超时或者拿到锁为止。
切记业务完成后一定要及时释放锁,设置默认的锁过期时间防止忘记释放锁。
*/
func (lock *Lock) GetDLockWithTtl(key string, waitSecond int, duration time.Duration) (error, bool) {
	successGetLock, err := lock.ops.SetNXTtl(LockPrefix+key, LockValue, duration)
	if err != nil {
		return err, successGetLock
	}
	if !successGetLock {
		beginTime := time.Now()
		for {
			successGetLock, err = lock.ops.SetNXTtl(LockPrefix+key, LockValue, duration)
			if successGetLock {
				break
			}
			if err != nil {
				break
			}
			currentTime := time.Now()
			diffTime := currentTime.Sub(beginTime)
			if diffTime.Seconds() >= float64(waitSecond) {
				break
			}
		}
	}
	return err, successGetLock
}

/*
*
释放锁
*/
func (lock *Lock) DisposeLock(key string) (error, bool) {
	res, err := lock.ops.DeleteKey(LockPrefix + key)
	return err, res > 0
}

// SetIfAbsent sets the value if the key does not exist (alias for SetNX)
func (lock *Lock) SetIfAbsent(key string, value interface{}) (bool, error) {
	return lock.ops.SetIfAbsent(key, value)
}

// SetIfAbsentWithTTL sets the value with TTL if the key does not exist
func (lock *Lock) SetIfAbsentWithTTL(key string, value interface{}, duration time.Duration) (bool, error) {
	return lock.ops.SetIfAbsentWithTTL(key, value, duration)
}
