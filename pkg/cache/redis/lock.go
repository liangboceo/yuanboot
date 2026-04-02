package redis

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Lock struct {
	ops           Ops
	uuid          string
	watchdogMutex sync.Mutex
	watchdogDone  chan struct{}
}

const (
	LockPrefix       = "dlock:"
	LockValue        = "DLOCK"
	DefaultWaitTime  = 10 * time.Second
	DefaultLeaseTime = 30 * time.Second
	WatchdogInterval = 10 * time.Second
)

// Lua scripts for atomic operations
var (
	// Try to acquire lock with TTL
	scriptTryLock = `
		if redis.call('setnx', KEYS[1], ARGV[1]) == 1 then
			redis.call('pexpire', KEYS[1], ARGV[2])
			return 1
		else
			if redis.call('get', KEYS[1]) == ARGV[1] then
				redis.call('pexpire', KEYS[1], ARGV[2])
				return 1
			end
			return 0
		end
	`

	// Unlock script - only delete if value matches
	scriptUnlock = `
		if redis.call('get', KEYS[1]) == ARGV[1] then
			return redis.call('del', KEYS[1])
		else
			return 0
		end
	`

	// Extend lock TTL script
	scriptExtend = `
		if redis.call('get', KEYS[1]) == ARGV[1] then
			return redis.call('pexpire', KEYS[1], ARGV[2])
		else
			return 0
		end
	`
)

func NewLock(ops Ops) *Lock {
	return &Lock{
		ops:          ops,
		uuid:         generateUUID(),
		watchdogDone: make(chan struct{}),
	}
}

func generateUUID() string {
	return fmt.Sprintf("%s-%d-%d",
		uuid()[:8],
		time.Now().UnixNano(),
		rand.Intn(100000))
}

func uuid() string {
	return fmt.Sprintf("%x%x%x",
		time.Now().UnixNano(),
		rand.Uint64(),
		rand.Uint64())
}

/*
TryLock attempts to acquire a lock with the specified lease time.
Returns true if lock is acquired, false otherwise.

Parameters:
  - key: lock key
  - leaseTime: lock expiration time (if <= 0, uses default)

Example:

	err, ok := lock.TryLock("mykey", 30*time.Second)
*/
func (lock *Lock) TryLock(key string, leaseTime time.Duration) (error, bool) {
	if leaseTime <= 0 {
		leaseTime = DefaultLeaseTime
	}

	fullKey := LockPrefix + key
	leaseMillis := leaseTime.Milliseconds()

	result, err := lock.ops.Eval(scriptTryLock, []string{fullKey}, lock.uuid, leaseMillis)
	if err != nil {
		return err, false
	}

	acquired := result.(int64) == 1
	if acquired {
		// Start watchdog for automatic extension
		go lock.startWatchdog(key, leaseTime)
	}

	return nil, acquired
}

/*
TryLockWithWait attempts to acquire a lock with wait timeout.
This is similar to Redisson's tryLock() method.

Parameters:
  - key: lock key
  - waitTime: maximum time to wait for lock
  - leaseTime: lock expiration time

Example:

	err, ok := lock.TryLockWithWait("mykey", 10*time.Second, 30*time.Second)
*/
func (lock *Lock) TryLockWithWait(key string, waitTime time.Duration, leaseTime time.Duration) (error, bool) {
	if leaseTime <= 0 {
		leaseTime = DefaultLeaseTime
	}

	fullKey := LockPrefix + key
	leaseMillis := leaseTime.Milliseconds()

	// Try to acquire lock immediately
	result, err := lock.ops.Eval(scriptTryLock, []string{fullKey}, lock.uuid, leaseMillis)
	if err != nil {
		return err, false
	}

	if result.(int64) == 1 {
		// Start watchdog for automatic extension
		go lock.startWatchdog(key, leaseTime)
		return nil, true
	}

	// Wait for lock with polling
	pollInterval := 50 * time.Millisecond

	beginTime := time.Now()
	for {
		result, err := lock.ops.Eval(scriptTryLock, []string{fullKey}, lock.uuid, leaseMillis)
		if err != nil {
			return err, false
		}

		if result.(int64) == 1 {
			go lock.startWatchdog(key, leaseTime)
			return nil, true
		}

		// Check timeout
		elapsed := time.Since(beginTime)
		if elapsed >= waitTime {
			return nil, false
		}

		// Sleep to avoid CPU spinning
		time.Sleep(pollInterval)
	}
}

/*
Unlock releases the lock.
This is safe - only the lock holder can release the lock.

Example:

	err := lock.Unlock("mykey")
*/
func (lock *Lock) Unlock(key string) error {
	// Stop watchdog first
	lock.stopWatchdog()

	fullKey := LockPrefix + key
	_, err := lock.ops.Eval(scriptUnlock, []string{fullKey}, lock.uuid)
	return err
}

/*
ForceUnlock forcibly releases the lock regardless of holder.
Use with caution!

Example:

	err := lock.ForceUnlock("mykey")
*/
func (lock *Lock) ForceUnlock(key string) error {
	lock.stopWatchdog()

	fullKey := LockPrefix + key
	_, err := lock.ops.DeleteKey(fullKey)
	return err
}

/*
Extend extends the lock expiration time.
Only the lock holder can extend the lock.

Example:

	err := lock.Extend("mykey", 30*time.Second)
*/
func (lock *Lock) Extend(key string, leaseTime time.Duration) error {
	fullKey := LockPrefix + key
	leaseMillis := leaseTime.Milliseconds()

	result, err := lock.ops.Eval(scriptExtend, []string{fullKey}, lock.uuid, leaseMillis)
	if err != nil {
		return err
	}

	if result.(int64) == 0 {
		return fmt.Errorf("failed to extend lock: not the owner")
	}

	return nil
}

/*
IsLocked checks if the lock is currently held.

Example:

	ok := lock.IsLocked("mykey")
*/
func (lock *Lock) IsLocked(key string) (bool, error) {
	fullKey := LockPrefix + key
	return lock.ops.Exists(fullKey)
}

// startWatchdog starts the watchdog goroutine for automatic lock extension
func (lock *Lock) startWatchdog(key string, leaseTime time.Duration) {
	lock.watchdogMutex.Lock()
	defer lock.watchdogMutex.Unlock()

	// Check if already running
	select {
	case <-lock.watchdogDone:
		lock.watchdogDone = make(chan struct{})
	default:
	}

	go func() {
		// Extend interval should be 1/3 of lease time, minimum 10 seconds
		extendInterval := leaseTime / 3
		if extendInterval < WatchdogInterval {
			extendInterval = WatchdogInterval
		}

		ticker := time.NewTicker(extendInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Extend lock TTL
				fullKey := LockPrefix + key
				result, err := lock.ops.Eval(scriptExtend, []string{fullKey}, lock.uuid, leaseTime.Milliseconds())
				if err != nil || result.(int64) == 0 {
					// Lock no longer exists or not owned
					return
				}
			case <-lock.watchdogDone:
				return
			}
		}
	}()
}

// stopWatchdog stops the watchdog goroutine
func (lock *Lock) stopWatchdog() {
	select {
	case <-lock.watchdogDone:
		return
	default:
		close(lock.watchdogDone)
	}
}

// contextKey type for context values
type contextKey string

const (
	lockKey contextKey = "distributed_lock"
)

// LockWithContext acquires a lock and stores it in context for automatic release.
// The lock will be automatically released when the context is cancelled.
//
// Example:
// ctx, cancel := context.WithCancel(context.Background())
// ctx, err := lock.LockWithContext(ctx, "mykey", 30*time.Second)
//
//	if err != nil {
//	    return err
//	}
//
// defer cancel()
// // do work with lock
func (lock *Lock) LockWithContext(ctx context.Context, key string, leaseTime time.Duration) (context.Context, error) {
	if leaseTime <= 0 {
		leaseTime = DefaultLeaseTime
	}

	err, ok := lock.TryLock(key, leaseTime)
	if err != nil {
		return ctx, err
	}

	if !ok {
		return ctx, fmt.Errorf("failed to acquire lock: %s", key)
	}

	// Create new context with cancel
	ctx, cancel := context.WithCancel(ctx)

	// Store lock info in context
	ctx = context.WithValue(ctx, lockKey, &lockInfo{
		key:    key,
		lock:   lock,
		cancel: cancel,
	})

	// When context is cancelled, release the lock
	go func() {
		select {
		case <-ctx.Done():
			_ = lock.Unlock(key)
		}
	}()

	return ctx, nil
}

// lockInfo holds lock information for context-based locking
type lockInfo struct {
	key    string
	lock   *Lock
	cancel context.CancelFunc
}

// GetLockFromContext retrieves lock info from context
func GetLockFromContext(ctx context.Context) (*lockInfo, bool) {
	info, ok := ctx.Value(lockKey).(*lockInfo)
	return info, ok
}

// Legacy methods for backward compatibility

/*
GetDLock acquires a lock with wait timeout.
Note: This method is deprecated, use TryLockWithWait instead.

Parameters:
  - key: lock key
  - waitSecond: wait timeout in seconds

Returns:
  - error: any error that occurred
  - bool: true if lock acquired, false otherwise
*/
func (lock *Lock) GetDLock(key string, waitSecond int) (error, bool) {
	return lock.TryLockWithWait(key, time.Duration(waitSecond)*time.Second, DefaultLeaseTime)
}

/*
GetDLockWithTtl acquires a lock with TTL and wait timeout.
Note: This method is deprecated, use TryLockWithWait instead.

Parameters:
  - key: lock key
  - waitSecond: wait timeout in seconds
  - duration: lock TTL

Returns:
  - error: any error that occurred
  - bool: true if lock acquired, false otherwise
*/
func (lock *Lock) GetDLockWithTtl(key string, waitSecond int, duration time.Duration) (error, bool) {
	return lock.TryLockWithWait(key, time.Duration(waitSecond)*time.Second, duration)
}

/*
DisposeLock releases the lock.
Note: This method is deprecated, use Unlock instead.

Parameters:
  - key: lock key

Returns:
  - error: any error that occurred
  - bool: true if lock released, false otherwise
*/
func (lock *Lock) DisposeLock(key string) (error, bool) {
	err := lock.Unlock(key)
	return err, err == nil
}

// SetIfAbsent sets the value if the key does not exist (alias for SetNX)
func (lock *Lock) SetIfAbsent(key string, value interface{}) (bool, error) {
	return lock.ops.SetIfAbsent(key, value)
}

// SetIfAbsentWithTTL sets the value with TTL if the key does not exist
func (lock *Lock) SetIfAbsentWithTTL(key string, value interface{}, duration time.Duration) (bool, error) {
	return lock.ops.SetIfAbsentWithTTL(key, value, duration)
}
