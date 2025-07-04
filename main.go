package sakha
import (
	"github.com/sakhaclothing/route"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("sakha", route.URL)
}
