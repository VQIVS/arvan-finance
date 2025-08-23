package rabbit

import "billing-service/pkg/constants"

func GetQueueName(key string) string {
	return constants.ServiceName + "_" + key
}
