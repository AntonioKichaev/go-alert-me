package grabbers

type Option func(rac *racoon)

func SetAllowMetrics(metrics map[string]struct{}) Option {
	return func(rac *racoon) {
		rac.allowMetrics = metrics
	}
}
