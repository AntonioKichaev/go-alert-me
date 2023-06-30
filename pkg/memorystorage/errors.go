package memorystorage

import "errors"

var ErrorNotExistMetric = errors.New("doesn't exist metric")
var ErrorFailedLoad = errors.New("couldn't load")
var ErrorFailedSave = errors.New("couldn't save")
