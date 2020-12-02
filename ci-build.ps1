GOOS=linux 
go build ./chapter/chapter.go
go build ./list/list.go
go build ./chapterlist/chapterlist.go 
mkdir outputs
zip ./chapter.zip chapter
zip ./chapterlist.zip chapterlist
zip ./list.zip list