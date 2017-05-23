package protocol

const (
	PING             = "ping"
	PONG             = "pong"
	LIST_OF_METRICS  = "list-of-metrics"
	LIST_OF_PLUGINS  = "list-of-plugins"
	GET_METRIC_VALUE = "get-metric-value"
	EMPTY            = "zalua-empty-message"

	UNKNOWN_METRIC  = "unknown-metric"
	UNKNOWN_COMMAND = "unknown-command"

	COMMAND_ERROR = "command-error"
	COMMAND_KILL  = "command-kill"
)
