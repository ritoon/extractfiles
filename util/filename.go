package util

import (
	"fmt"
	"regexp"
	"roumet/logger"
)

func ExtractLotNumber(url string) string {

	// Le regex capture le nombre avant ".jpg" ou "-X.jpg"
	re := regexp.MustCompile(`/(\d+)(?:-\d+)?\.(?:jpg|jpeg|png)$`)
	match := re.FindStringSubmatch(url)
	if match == nil {
		logger.Error("util", fmt.Errorf("Aucun numéro de lot trouvé dans l'URL : %s", url))
		return ""
	}
	return match[1]
}
