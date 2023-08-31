package test

import (
	"math/rand"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/m1al04949/avito-tech-service/internal/http-server/handlers/addtouser"
	"github.com/m1al04949/avito-tech-service/internal/http-server/handlers/adduser"
	"github.com/m1al04949/avito-tech-service/internal/http-server/handlers/createsegment"
	"github.com/m1al04949/avito-tech-service/internal/model"
)

const (
	host = "localhost:8082"
)

func TestAvitoService_EasyHappy(t *testing.T) {

	chars := []rune("ABCDEFGIHJKLMNOPQRSTVUWXYZ" +
		"0123456789_")

	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	user := rand.Intn(1000)

	// Create some user
	e.POST("/users").
		WithJSON(adduser.Request{
			UserID: user,
		}).
		WithBasicAuth("myuser", "mypass").
		Expect().
		Status(200).
		JSON().
		Object()

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	segm := make([]rune, 10)

	for i := range segm {
		segm[i] = chars[rnd.Intn(len(chars))]
	}
	segment := string(segm)

	// Create some segment
	e.POST("/segments").
		WithJSON(createsegment.Request{
			Slug: segment,
		}).
		WithBasicAuth("myuser", "mypass").
		Expect().
		Status(200).
		JSON().
		Object()

	var segments []model.Segment
	newsegment := model.Segment{
		Slug: segment,
	}
	segments = append(segments, newsegment)
	id := strconv.Itoa(user)

	// Add this segment to user
	e.POST("/users/id="+id).
		WithJSON(addtouser.Request{
			Segments: segments,
		}).
		WithBasicAuth("myuser", "mypass").
		Expect().
		Status(200).
		JSON().
		Object()

	// Get segment from user
	e.GET("/users/id="+id).
		WithBasicAuth("myuser", "mypass").
		Expect().
		Status(200).
		JSON().
		Object()
}
