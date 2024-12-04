.PHONY: setup build

BINARY_NAME=x-straight-check

build:
	npm i && \
	go mod tidy && \
	go generate && \
	go build -ldflags="-w -s" -o $(BINARY_NAME)

generate:
	templ generate

dev: generate clean build tailwindcss-build
	./$(BINARY_NAME)

setup: tailwindcss-init
	mkdir -p assets/js/htmx &&\
	mkdir -p assets/css &&\
	wget https://unpkg.com/htmx-ext-loading-states@2.0.0/loading-states.js -O assets/js/htmx/loading-states.js &&\
	wget https://unpkg.com/htmx.org@2.0.3/dist/htmx.min.js -O assets/js/htmx/htmx.min.js

tailwindcss-init:
	npm i &&\
	npx tailwindcss build -i assets/css/style.css -o assets/css/tailwind.css -m

tailwindcss-build:
	npx tailwindcss build -i assets/css/style.css -o assets/css/tailwind.css -m

clean:
	go clean
	rm -rf $(BINARY_NAMR)

