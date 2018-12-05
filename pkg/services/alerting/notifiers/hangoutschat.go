package notifiers

import (
	"encoding/json"

	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/log"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/alerting"
)

func init() {
	alerting.RegisterNotifier(&alerting.NotifierPlugin{
		Type:        "hangoutschat",
		Name:        "Google Hangouts Chat",
		Description: "Sends a specified message to a given Hangouts Chat URL",
		Factory:     NewHangoutChatNotifier,
		OptionsTemplate: `
      <h3 class="page-heading">Webhook settings</h3>
      <div class="gf-form">
        <span class="gf-form-label width-10">Url</span>
        <input type="text" required class="gf-form-input max-width-26" ng-model="ctrl.model.settings.url"></input>
      </div>
    `,
	})

}

func NewHangoutChatNotifier(model *m.AlertNotification) (alerting.Notifier, error) {
	url := model.Settings.Get("url").MustString()
	if url == "" {
		return nil, alerting.ValidationError{Reason: "Could not find url property in settings"}
	}

	return &HangoutChatNotifier{
		NotifierBase: NewNotifierBase(model),
		Url:          url,
		log:          log.New("alerting.notifier.webhook"),
	}, nil
}

type HangoutChatNotifier struct {
	NotifierBase
	Url string
	log log.Logger
}

func (this *HangoutChatNotifier) Notify(evalContext *alerting.EvalContext) error {
	this.log.Info("Sending Hangout Chat Notification")

	bodyJSON := simplejson.New()

	if evalContext.Rule.Message != "" {
		bodyJSON.Set("text", evalContext.Rule.Message)
	}

	body, _ := bodyJSON.MarshalJSON()

	if evalContext.ImagePublicUrl != "" {
		bodyMap := map[string]interface{}{
			"cards": []map[string]interface{}{
				{
					"sections": []map[string]interface{}{
						{
							"sections": []map[string]interface{}{
								{
									"image": map[string]interface{}{
										"imageUrl": evalContext.ImagePublicUrl,
									},
								},
							},
						},
					},
					"text": evalContext.Rule.Message,
				},
			},
		}
		card, _ := json.Marshal(&bodyMap)

		body = card

	}

	cmd := &m.SendWebhookSync{
		Url:        this.Url,
		Body:       string(body),
		HttpMethod: "POST",
	}

	if err := bus.DispatchCtx(evalContext.Ctx, cmd); err != nil {
		this.log.Error("Failed to send chat notification", "error", err, "webhook", this.Name)
		return err
	}

	return nil
}
