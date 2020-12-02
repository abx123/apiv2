GOOS=linux 
go build -o output/novels ./novels/novels.go
go build -o output/chapters ./chapters/chapters.go
go build -o output/chapter ./chapter/chapter.go 
mkdir zip
zip ./zip/novels.zip ./output/novels
zip ./zip/chapters.zip ./output/chapters
zip ./zip/chapter.zip ./output/chapter
