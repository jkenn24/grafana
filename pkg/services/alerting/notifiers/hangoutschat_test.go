package notifiers

import (
	"testing"

	"github.com/grafana/grafana/pkg/components/simplejson"
	m "github.com/grafana/grafana/pkg/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHangoutsChatNotifier(t *testing.T) {
	Convey("Hangouts Chat notifier tests", t, func() {

		Convey("Parsing alert notification from settings", func() {
			Convey("empty settings should return error", func() {
				json := `{ }`

				settingsJSON, _ := simplejson.NewJson([]byte(json))
				model := &m.AlertNotification{
					Name:     "ops",
					Type:     "hangoutschat",
					Settings: settingsJSON,
				}

				_, err := NewHangoutChatNotifier(model)
				So(err, ShouldNotBeNil)
			})

			Convey("from settings", func() {
				json := `
				{
          "url": "http://google.com"
				}`

				settingsJSON, _ := simplejson.NewJson([]byte(json))
				model := &m.AlertNotification{
					Name:     "ops",
					Type:     "hangoutschat",
					Settings: settingsJSON,
				}

				not, err := NewHangoutChatNotifier(model)
				hangoutChatNotifier := not.(*HangoutChatNotifier)

				So(err, ShouldBeNil)
				So(hangoutChatNotifier.Name, ShouldEqual, "ops")
				So(hangoutChatNotifier.Type, ShouldEqual, "hangoutschat")
				So(hangoutChatNotifier.Url, ShouldEqual, "http://google.com")
			})
		})
	})
}
