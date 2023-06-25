package config

type Option func(*Agent)

func SetLoggingLevel(level string) Option {
	return func(agent *Agent) {
		agent.LoggingLevel = level
	}
}
