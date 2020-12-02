GOOS=linux 
go build -o output/novels ./chapter/novels.go
go build -o output/chapters ./list/chapters.go
go build -o output/chapter ./chapterlist/chapter.go 
mkdir zip
zip ./zip/novels.zip ./output/novels
zip ./zip/chapters.zip ./output/chapters
zip ./zip/chapter.zip ./output/chapter
