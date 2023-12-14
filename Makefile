run: tmpl tailwind
	TAAVI_ENV=dev go run cmd/http/main.go

tmpl:
	templ generate

tailwind:
	npx tailwindcss -i input.css -o public/style.css
