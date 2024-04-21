package logging

import (
	"context"
	"encoding/json"
	"fmt"
)

type JSONLogger struct {
}

func (log JSONLogger) Log(msg string, priority int, datas ...map[string]interface{}) {
	providedDatas := mergeMaps(datas...)
	totalDatas := mergeMaps(providedDatas, map[string]interface{}{
		"priority": priority,
		"message":  msg,
	})
	encoded, err := json.Marshal(totalDatas)
	if err != nil {
		Error(context.Background(), "Could not Marshal log", map[string]interface{}{"originalmessage": msg, "marshalerror": err.Error()})
		return
	}
	fmt.Printf("%s\n", string(encoded))
}
