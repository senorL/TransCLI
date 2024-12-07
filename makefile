# Build for all platforms
build:
	goreleaser release --clean --skip=publish
