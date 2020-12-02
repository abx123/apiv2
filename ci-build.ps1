GOOS=linux 
go build ./chapter/chapter.go
go build ./list/list.go
go build ./chapterlist/chapterlist.go 
mkdir outputs
zip ./chapter/chapter.zip chapter
zip ./chapter/chapterlist.zip chapterlist
zip ./chapter/list.zip list