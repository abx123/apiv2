GOOS=linux 
go build -o chapter ./chapter/chapter.go
go build -o list ./list/list.go
go build -o chapterlist ./chapterlist/chapterlist.go 
mkdir outputs
zip ./chapter/chapter.zip chapter
zip ./chapterlist/chapterlist.zip chapterlist
zip ./list/list.zip list