
build:
	go build -o prayertimes 

run: build
	./prayertimes

reset:
	rm ~/.config/prayertimes/config.yaml
	rm ~/.cache/prayertimes/prayertimes.sqlite
