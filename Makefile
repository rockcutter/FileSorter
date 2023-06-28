FileSorter.exe: main.go FileSorter/FileSorter.go
	go build -o FileSorter.exe main.go

clean:; rm FileSorter.exe