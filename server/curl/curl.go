package curl

import (
	"github.com/gin-gonic/gin"

	"github.com/Depado/goploader/server/utils"
)

// WriteTutorial will write to a gin.Context's writer the full curl tutorial
func WriteTutorial(c *gin.Context) {
	w := c.Writer
	Header(w, header)
	Title(w, "Usage\n")
	Standard(w, "Upload from stdin")
	Command(w, "$ tree | curl -F file=@- "+utils.ServerURI(c))
	Standard(w, "Simple file upload")
	Command(w, "$ curl -F file=@myfile.txt "+utils.ServerURI(c))
}
