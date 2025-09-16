module github.com/benkoppe/bear-trak-backend/go-server

go 1.24.1

toolchain go1.24.2

require (
	github.com/PuerkitoBio/goquery v1.10.3
	github.com/amit7itz/goset v1.2.1
	github.com/bmizerany/pq v0.0.0-20131128184720-da2b95e392c1
	github.com/chromedp/cdproto v0.0.0-20250509201441-70372ae9ef75
	github.com/chromedp/chromedp v0.13.6
	github.com/emersion/go-imap v1.2.1
	github.com/go-faster/errors v0.7.1
	github.com/go-faster/jx v1.1.0
	github.com/jackc/pgx/v5 v5.7.4
	github.com/jamespfennell/gtfs v0.1.24
	github.com/jhillyerd/enmime v1.3.0
	github.com/ogen-go/ogen v1.12.0
	github.com/pressly/goose/v3 v3.24.3
	github.com/revrost/go-openrouter v0.0.0-20250414052218-c9123df8a97e
	github.com/twpayne/go-polyline v1.1.1
	go.opentelemetry.io/otel v1.35.0
	go.opentelemetry.io/otel/metric v1.35.0
	go.opentelemetry.io/otel/trace v1.35.0
	golang.org/x/sync v0.14.0
	golang.org/x/text v0.25.0
)

require (
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/cention-sany/utf7 v0.0.0-20170124080048-26cad61bd60a // indirect
	github.com/chromedp/sysutil v1.1.0 // indirect
	github.com/dlclark/regexp2 v1.11.5 // indirect
	github.com/emersion/go-sasl v0.0.0-20241020182733-b788ff22d5a6 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-faster/yaml v0.4.6 // indirect
	github.com/go-json-experiment/json v0.0.0-20250417205406-170dfdcf87d1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.4.0 // indirect
	github.com/gogs/chardet v0.0.0-20211120154057-b7413eaefb8f // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jaytaylor/html2text v0.0.0-20230321000545-74c2419ad056 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/olekukonko/tablewriter v1.0.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	github.com/ssor/bom v0.0.0-20170718123548-6386211fdfcf // indirect
	go.opencensus.io v0.22.3 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/exp v0.0.0-20250506013437-ce4c2cf36ca6 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	googlemaps.github.io/maps v1.7.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

// temporary fix: html2text currently broken
// https://github.com/jaytaylor/html2text/issues/67
replace github.com/olekukonko/tablewriter => github.com/olekukonko/tablewriter v0.0.5
