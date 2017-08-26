package curl

import "io"

// WriteTutorial will write to an io.Writer the full curl tutorial
func WriteTutorial(w io.Writer) {
	Header(w, header)
	Title(w, "Usage\n")
	Standard(w, "Upload from stdin")
	Command(w, "$ tree | curl -F file=@- https://gpldr.in/\n")
	Standard(w, "Simple file upload")
	Command(w, "$ curl -F file=@myfile.txt https://gpldr.in/")
}
