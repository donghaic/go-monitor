package dao

import (
	"gt-monitor/common/mongo"
	"gt-monitor/models"
	"gt-monitor/utils"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"strings"
	"time"
)

const (
	CLICK_LOG_DB   = "ClickLog"
	CONV_LOG_DB    = "EventLog"
	CONV_LOG_TABLE = "event"

	GT_REPORT_DB    = "GtStats"
	GT_REPORT_TABLE = "report"
)

type MongoDao struct {
	logMgoClient    *mongo.MongoPool
	reportMgoClient *mongo.MongoPool
}

func NewMongoDao(logMgoClient, reportMgoClient *mongo.MongoPool) *MongoDao {
	return &MongoDao{logMgoClient, reportMgoClient}
}

func (m *MongoDao) SaveConv(conv *models.ConvEntity) error {
	session := m.logMgoClient.NewSession()
	defer session.Close()
	err := session.DB(CONV_LOG_DB).C(CONV_LOG_TABLE).Insert(conv)
	return err
}

func (m *MongoDao) FindConvByImpId(impId string) (*models.ConvEntity, error) {
	query := bson.M{"impid": impId}

	session := m.logMgoClient.NewSession()
	defer session.Close()

	result := models.ConvEntity{}

	err := session.DB(CONV_LOG_DB).C(CONV_LOG_TABLE).Find(query).One(&result)

	return &result, err
}

func (m *MongoDao) FindClickByImpId(impId string, clickTsStr string) (*models.ClickEntity, error) {
	clickTs, _ := strconv.ParseInt(clickTsStr, 10, 64)
	collName := time.Unix(clickTs/1000000000, 0).Format("2006-01-02")
	query := bson.M{"impid": impId}
	println(collName, query)

	session := m.logMgoClient.NewSession()
	defer session.Close()

	result := models.ClickEntity{}
	err := session.DB(CLICK_LOG_DB).C(collName).Find(query).One(&result)

	return &result, err
}

func (m *MongoDao) UpsertReport(clickEntity *models.ClickEntity, event *models.ConvReqParams) error {
	timeclk, now := time.Unix(event.EventTs/1000000000, 0), time.Now()

	timeDate := time.Date(timeclk.Year(), timeclk.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())

	timeStamp := strconv.FormatInt(timeDate.Unix()*1000, 10)

	if 0 == len(clickEntity.TagId) {
		clickEntity.TagId = "0"
	}

	keyStr := strings.Join([]string{clickEntity.AdGroupId, clickEntity.MediaId, clickEntity.CreativeId,
		clickEntity.AdvId, clickEntity.CampaignId, timeStamp,
		clickEntity.ProductId, clickEntity.AdPosId, clickEntity.AdvLinkTagId, clickEntity.UserId, clickEntity.UserBudgetId, clickEntity.TagId}, "__")
	keyHash := utils.StringHashCode(keyStr)

	dstrtmp := time.Unix(event.EventTs/1000000000, 0).Format("2006-01-02 15")
	selector := bson.M{
		"keyStr":       keyStr,
		"keyHash":      keyHash,
		"dateTimeStr":  dstrtmp,
		"date":         strings.Split(dstrtmp, " ")[0],
		"time":         strings.Split(dstrtmp, " ")[1],
		"dateTime":     timeDate.Add(time.Hour * 8),
		"timeStamp":    timeDate.Unix() * 1000,
		"adverId":      clickEntity.AdvId,
		"productId":    clickEntity.ProductId,
		"campaignId":   clickEntity.CampaignId,
		"adgroupId":    clickEntity.AdGroupId,
		"mediaId":      clickEntity.MediaId,
		"creativeId":   clickEntity.CreativeId,
		"adPos":        clickEntity.AdPosId,
		"advLinkTagId": clickEntity.AdvLinkTagId,
		"userId":       clickEntity.UserId,
		"userBudgetId": clickEntity.UserBudgetId,
		"tagId":        clickEntity.TagId,
	}

	update := bson.M{
		"totalCvs": 1,
		"idpCvs":   0,
		"dupCvs":   0,
	}
	if event.EventType == "duplicated" {
		update["dupCvs"] = 1
	}
	if event.EventType == "normal" {
		update["idpCvs"] = 1
	}

	session := m.reportMgoClient.NewSession()
	defer session.Close()

	_, err := session.DB(GT_REPORT_DB).C(GT_REPORT_TABLE).Upsert(selector, bson.M{"$inc": update})

	return err
}
