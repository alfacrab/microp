build:
	go build -o microp -buildmode=exe microp.go
	tar cvJf microp.tar.xz microp && rm -rf microp
