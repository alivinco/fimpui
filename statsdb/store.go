package statsdb

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/protobuf"
	metrics "github.com/rcrowley/go-metrics"
	log "github.com/Sirupsen/logrus"
	"time"
)

type StatsStore struct {
	db              *storm.DB
	dbFileName      string
	metricsRegistry metrics.Registry
	maxEventLogSize int
	eventLogSize    int
}

func NewStatsStore(storeFile string) *StatsStore {
	store := StatsStore{dbFileName: storeFile}
	store.Init()
	store.metricsRegistry = metrics.NewRegistry()
	return &store
}

func (st *StatsStore) Init(){
	st.ConnectDB()
	st.maxEventLogSize = 10000
	st.eventLogSize = st.CountEvents()
	log.Info("<Stats> Error log size = ",st.eventLogSize)

}

func (st *StatsStore) ConnectDB() error {
	var err error
	st.db, err = storm.Open("stats.db", storm.Codec(protobuf.Codec))
	if err != nil {
		log.Error("<Stats> Can't open DB file . Error : ", err)
		return err
	}
	err = st.db.Init(&EventRec{})
	if err != nil {
		log.Error("<Stats> Can't Init Error . Error : ", err)
		return err
	}
	return nil
}


func (st *StatsStore) DisconnectDB() {
	st.db.Close()
}

func (st *StatsStore) AddEvent(errMsg *EventRec) {
	errMsg.Timestamp = time.Now()
	err := st.db.Save(errMsg)
	if err != nil {
		log.Error("<Stats> Can't register event. Error : ", err)
	} else {
		st.eventLogSize +=1
		st.ApplyRetentionPolicy()
		log.Debug("<Stats> Event registered")
	}
}

func (st *StatsStore) GetAllEvents() ([]EventRec, error) {
	var fimpErrors []EventRec
	err := st.db.Select().Find(&fimpErrors)
	return fimpErrors, err
}

func (st *StatsStore) ApplyRetentionPolicy() {
	delCounter := 0
	delBatchSize := 100
	if st.eventLogSize > st.maxEventLogSize {
		evtRecords, err := st.GetAllEvents()
		if err == nil {
			log.Error("<Stats> Rotating logs ")
			for i := range evtRecords {
				if delCounter > delBatchSize {
					st.eventLogSize = st.eventLogSize - delCounter
					break
				}
				st.db.DeleteStruct(&evtRecords[i])
				delCounter+=1
			}
		}else {
			log.Error("<Stats> Failed to rotate event log. Error : ", err)
		}


	}
}

func (st *StatsStore) CountEvents() int {
	var fimpEvents []EventRec
	st.db.All(&fimpEvents)
	return len(fimpEvents)
}

func (st *StatsStore) GetEvents(pageSize int , page int) ([]EventRec, error) {
	var fimpEvents []EventRec
	err := st.db.Select().Limit(pageSize).Skip(pageSize*page).Reverse().Find(&fimpEvents)
	return fimpEvents, err
}