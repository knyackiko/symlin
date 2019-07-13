# Symlin
Symlin manages symbolic links.

# Installation
```
$ go get github.com/kyklades/symlin
```

# How to use
## List up symbolic links
```
$ symlin list
toDoc1 -> dir1/dir1-1/dir1-1-1/doc1.md
toDoc4 -> dir1/dir1-2/doc4.md
toDoc5 -> dir2/dir2-1/doc5.md
toDoc7 -> dir3/doc7.md
```

## Create new symbolic links
```
$ symlin create path/to/target/file path/to/new/symbolic/link
New symbolic link has been created!
toDoc7 -> testdata/dir3/doc7.md
```

## Unlink existing symbolic links
```
$ symlin unlink toDoc7
toDoc7 has been successfully unlinked!
```