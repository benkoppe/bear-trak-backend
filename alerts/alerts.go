package alerts

import (
	"github.com/benkoppe/bear-trak-backend/alerts/static"
	backend "github.com/benkoppe/bear-trak-backend/backend"
)

func Get() ([]backend.Alert, error) {
	staticAlerts := static.GetAlerts()
	var alerts []backend.Alert

	for _, staticAlert := range staticAlerts {
		alerts = append(alerts, convertStatic(staticAlert))
	}

	return alerts, nil
}

func convertStatic(static static.Alert) backend.Alert {
	button := backend.NilAlertButton{Null: true}

	if static.Button != nil {
		button = backend.NewNilAlertButton(backend.AlertButton{
			Title: static.Button.Title,
			URL:   *static.Button.URL.URL,
		})
	}

	return backend.Alert{
		ID:       static.ID,
		Title:    static.Title,
		Message:  static.Message,
		Enabled:  static.Enabled,
		ShowOnce: static.ShowOnce,
		Button:   button,
	}
}
