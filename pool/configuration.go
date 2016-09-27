package pool

import (
	"time"
)

// Configuration структура настроек пулла соединений к редису
type Configuration struct {
	WaitConnection           bool          // Ожидание свободного подключения при достижении ActivePoolSize
	MaxIdleConnectionCount   int           // Количество соеднинений в режиме ожидания. Если 0, то в пулле не сохраняется соединение
	MaxActiveConnectionCount int           // Максимальное количество соединений. Если 0, то неограниченно
	IdleConnectionTimeout    time.Duration // Время хранения соединения в пулле
	CheckConnectionFrequency time.Duration // Таймаут проверки доступности редиса
}
