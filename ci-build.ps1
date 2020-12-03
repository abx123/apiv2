GOOS=linux 
go build -o output/novels ./novels/novels.go
go build -o output/chapters ./chapters/chapters.go
go build -o output/chapter ./chapter/chapter.go 
cd output
# zip ./novels.zip ./novels
# zip ./chapters.zip ./chapters
# zip ./chapter.zip ./chapter
zip ./package.zip ./chapter ./chapters ./novels
