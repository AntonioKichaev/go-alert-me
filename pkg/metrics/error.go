package metrics

import "errors"

var ErrorName = errors.New("counter: error metric name")
var ErrorUnknownMetricType = errors.New("counter: unknown metric type")
var ErrorBadValue = errors.New("bad value")
