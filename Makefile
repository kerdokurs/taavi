run: tailwind
	TAAVI_ENV=dev go run cmd/http/main.go

tailwind:
	npx tailwindcss -i input.css -o public/style.css
