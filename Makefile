test:
	go test -v -race .

bench:
	go test -run - -bench=. -benchmem
