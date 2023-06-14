package quota

import (
	infraFirestore "github.com/Nekodigi/charge-framework-backend/infrastructure/firestore"
	"github.com/gin-gonic/gin"
)

type (
	Quota struct {
		Fs *infraFirestore.Firestore
	}
)

func (q *Quota) Handle(e *gin.Engine) {
	q.HandleCheckQuota(e)
	q.HandleUseQuota(e)
	q.HandleAddQuota(e)
}
