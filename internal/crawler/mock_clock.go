package crawler

import "time"

// MockClock 用于在测试中模拟时间流逝
type MockClock struct {
	now time.Time
}

// Now 返回当前模拟时间
func (m *MockClock) Now() time.Time {
	return m.now
}

// Advance 推进模拟时钟的时间
func (m *MockClock) Advance(d time.Duration) {
	m.now = m.now.Add(d)
}