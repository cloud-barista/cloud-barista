package topic

import (
	"fmt"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/rest"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AddDeleteTopicToQueue(c echo.Context) error {
	delTopic := c.Param("topic")
	if err := util.RingQueuePut(types.TopicDel, delTopic); err != nil {
		return c.JSON(http.StatusInternalServerError, rest.SetMessage(fmt.Sprintf("failed to update metadata, error=%s", err)))
	}
	return c.JSON(http.StatusOK, delTopic)
}
