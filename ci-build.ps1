GOOS=linux 
go build -o output/chapter ./chapter/chapter.go
go build -o output/list ./list/list.go
go build -o output/chapterlist ./chapterlist/chapterlist.go 
mkdir zip
zip ./zip/chapter.zip ./output/chapter
zip ./zip/chapterlist.zip ./output/chapterlist
zip ./zip/list.zip ./output/list
