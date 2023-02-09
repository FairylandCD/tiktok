package controller

import (
	"dousheng/data"
	"dousheng/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type FeedResponse struct {
	data.Response
	NextTime  int64        `json:"next_time,omitempty"`
	VideoList []data.Video `json:"video_list,omitempty"` //默认空
}
type VideoListResponse struct {
	data.Response
	VideoList []data.Video `json:"video_list"`
}

// Feed:"douyin/feed/接口"
func Feed(c *gin.Context) {
	strLatestTime := c.Query("latest_time")
	log.Print("Received latest_time: " + strLatestTime)
	var latestTime time.Time
	if strLatestTime != "" {
		int64Time, _ := strconv.ParseInt(strLatestTime, 10, 64)
		latestTime = time.Unix(int64Time/1000, 0)
	} else {
		latestTime = time.Now()
	}
	log.Printf("latestTime UTS:%v", latestTime)

	// TODO:user_id JWT
	strToken := c.Query("token")
	userToken, _ := strconv.ParseInt(strToken, 10, 64)
	userID := userToken

	var videoService service.VideoServiceImpl

	videoList, nextTime, err := videoService.Feed(latestTime, userID)
	if err != nil {
		log.Printf("failed with videoService.Feed(latestTime, userID) : %v", err)
		c.JSON(http.StatusOK, FeedResponse{
			Response: data.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  data.Response{StatusCode: 0, StatusMsg: "succeed"},
		NextTime:  nextTime,
		VideoList: videoList,
	})
}

// Publish check token then save upload file to public directory
//func Publish(c *gin.Context) {
//	token := c.PostForm("token")
//
//	if _, exist := usersLoginInfo[token]; !exist {
//		c.JSON(http.StatusOK, data.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
//		return
//	}
//
//	data1, err := c.FormFile("data")
//	if err != nil {
//		c.JSON(http.StatusOK, data.Response{
//			StatusCode: 1,
//			StatusMsg:  err.Error(),
//		})
//		return
//	}
//
//	filename := filepath.Base(data1.Filename)
//	user := usersLoginInfo[token]
//	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
//	saveFile := filepath.Join("./public/", finalName)
//	if err := c.SaveUploadedFile(data, saveFile); err != nil {
//		c.JSON(http.StatusOK, data.Response{
//			StatusCode: 1,
//			StatusMsg:  err.Error(),
//		})
//		return
//	}
//
//	c.JSON(http.StatusOK, data.Response{
//		StatusCode: 0,
//		StatusMsg:  finalName + " uploaded successfully",
//	})
//}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {

	strUserId := c.Query("user_id")
	userId, _ := strconv.ParseInt(strUserId, 10, 64)
	log.Printf("user_id : %v", userId)

	var videoService service.VideoServiceImpl
	publishList, err := videoService.PubList(userId)
	if err != nil {
		log.Printf("failed with videoService.PubList(userId) : %v", err)
		c.JSON(http.StatusOK, FeedResponse{
			Response: data.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: data.Response{
			StatusCode: 0, StatusMsg: "succeed",
		},
		VideoList: publishList,
	})
}