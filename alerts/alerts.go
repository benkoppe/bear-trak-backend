package alerts

import (
	"github.com/benkoppe/bear-trak-backend/alerts/static"
	"github.com/benkoppe/bear-trak-backend/api"
)

func Get() ([]api.Alert, error) {
	staticAlerts := static.GetAlerts()
	var alerts []api.Alert

	for _, staticAlert := range staticAlerts {
		// filter disabled alerts
		if !staticAlert.Enabled {
			continue
		}

		alerts = append(alerts, convertStatic(staticAlert))
	}

	return alerts, nil
}

func convertStatic(static static.Alert) api.Alert {
	button := api.NilAlertButton{Null: true}
	if static.Button != nil {
		button = api.NewNilAlertButton(api.AlertButton{
			Title: static.Button.Title,
			URL:   *static.Button.URL.URL,
		})
	}

	maxBuild := api.NilInt{Null: true}
	if static.MaxBuild != nil {
		maxBuild = api.NewNilInt(*static.MaxBuild)
	}

	return api.Alert{
		ID:       static.ID,
		Title:    static.Title,
		Message:  static.Message,
		Enabled:  static.Enabled,
		ShowOnce: static.ShowOnce,
		Button:   button,
		MaxBuild: maxBuild,
	}
}
