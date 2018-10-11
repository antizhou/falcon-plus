package utils

import "github.com/open-falcon/falcon-plus/modules/agent/log_collector/common/scheme"

func DeepCopyStrategy(ori *scheme.Strategy) *scheme.Strategy {
	ret := &scheme.Strategy{
		ID:         ori.ID,
		Name:       ori.Name,
		FilePath:   ori.FilePath,
		TimeFormat: ori.TimeFormat,
		Pattern:    ori.Pattern,
		Interval:   ori.Interval,
		Tags:       DeepCopyStringMap(ori.Tags),
		Func:       ori.Func,
		Degree:     ori.Degree,
		Comment:    ori.Comment,
	}
	return ret
}

func DeepCopyStringSlice(params []string) []string {
	ret := make([]string, len(params))
	for i, str := range params {
		ret[i] = str
	}
	return ret
}

func DeepCopyStringMap(params map[string]string) map[string]string {
	ret := make(map[string]string, len(params))
	for k, v := range params {
		ret[k] = v
	}
	return ret
}
