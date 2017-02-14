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

	// Title(w, "Introduction")
	// Standard(w, intro)
	// Title(w, "Usage and Examples\n")
	// Command(w, "$ cat myfile.txt | curl -F file=@- https://gpldr.in/")
	// Explanation(w, "Reads data from stdin\n")
	// Command(w, "$ curl -F file=@myfile.txt https://gpldr.in/")
	// Explanation(w, "Your file will be named myfile.txt\n")
	// Command(w, "$ curl -F name=\"myamazingfile!\" -F file=@myfile.txt https://gpldr.in/")
	// Explanation(w, "Your file will be named myamazingfile!\n")
}
