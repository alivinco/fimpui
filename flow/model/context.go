package model

import (
	"time"
	//"io/ioutil"
	//"encoding/json"
	//"path"
	"bytes"
	"github.com/boltdb/bolt"
	log "github.com/Sirupsen/logrus"
	"encoding/gob"
)

type Variable struct {
	Value     interface{}
	ValueType string
}

type ContextRecord struct {
	Name string
	Description string
	UpdatedAt time.Time
	Variable Variable
}

type Context struct {
	storageLocation string
	db *bolt.DB
	inMemoryStore map[string][]ContextRecord
}

func NewContextDB(storageLocation string) (*Context , error) {
	var err error
	ctx := Context{}
	ctx.inMemoryStore = make(map[string][]ContextRecord)
	ctx.db, err = bolt.Open(storageLocation, 0600, nil)
	if err != nil {
		log.Error(err)
		return nil ,err
	}
	ctx.RegisterFlow("global")
	return &ctx,nil
}
func (ctx *Context) Close() {
	ctx.db.Close()
}
func (ctx *Context) RegisterFlow(flowId string ) error {
	ctx.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(flowId))
		if err != nil {
			log.Errorf("<ctx> Can't create bucket %s . Error: %s",flowId, err)
			return err
		}
		return nil
	})
	log.Infof("<ctx> Flow %s is registered in store.",flowId)
	return nil
}

func (ctx *Context) UnregisterFlow(flowId string) error {
	ctx.db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(flowId))
		if err != nil {
			log.Errorf("<ctx> Can't delete bucket %s . Error: %s",flowId, err)
			return err
		}
		return nil
	})
	log.Info("<ctx> Flow %s is deleted .",flowId)
	return nil
}

//func (ctx *Context) SaveToStorage() error {
	//fullPath := path.Join(ctx.storageLocation,ctx.flowName,"json")
	//log.Debugf("<FlMan> Saving context to file %s ",ctx.flowName)
	//binBody , err := json.Marshal(ctx)
	//if err != nil {
	//	log.Error("<Ctx> Can't serialize context . Error : ",err)
	//	return err
	//}
	//err = ioutil.WriteFile(fullPath, binBody, 0644)
	//if err != nil {
	//	log.Error("<Ctx> Can't save context to file . Error : ",err)
	//	return err
	//}
	//return nil
//}

//func (ctx *Context) LoadFromStorage() {
	//fullPath := path.Join(ctx.storageLocation,ctx.flowName,"json")
	//contextFileBody, err := ioutil.ReadFile(fullPath)
	//err = json.Unmarshal(contextFileBody, ctx)
	//if err != nil {
	//	log.Error("")
	//}
//}

func (ctx *Context) SetVariable(name string,valueType string,value interface{},description string,flowId string,inMemory bool ) error {
	rec := ContextRecord{Name:name,UpdatedAt:time.Now(),Description:description,Variable: Variable{ValueType:valueType,Value:value}}
	return ctx.PutRecord(&rec,flowId,inMemory)
}

func (ctx *Context) PutRecord(rec *ContextRecord,flowId string,inMemory bool ) error {
	if inMemory {
		//ctx.inMemoryStore[flowId] = *rec
	} else {
		err := ctx.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(flowId))
			data , err := ctx.encodeRecord(rec)
			if err != nil {
				return err
			}
			err = b.Put([]byte(rec.Name), data)
			return err
		})
		return err
	}
	return nil
}


func (ctx *Context) GetVariable(name string,flowId string) (Variable,error) {
	rec , err := ctx.GetRecord(name,flowId)
	return rec.Variable , err
}

func (ctx *Context) GetRecord(name string,flowId string) (*ContextRecord,error) {
	//rec , ok := ctx.inMemoryRecords[name]
	//if ok {
	//		return &rec,nil
	//}

	var ctxRec *ContextRecord
	var err error
	ctx.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(flowId))
			data := b.Get([]byte(name))
			ctxRec,err = ctx.decodeRecord(data)
			return nil
	})
	if err != nil {
		return nil , err
	}
	return ctxRec, err
}



func (ctx *Context) GetRecords(flowId string) []ContextRecord  {
	result := []ContextRecord{}

	ctx.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(flowId))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			rec ,err := ctx.decodeRecord(v)
			if err == nil {
				result = append(result,*rec)
			}else {
				log.Errorf("Can't decode record = %s , %s",k,err)
			}

		}

		return nil
	})

	return result
}


func (ctx *Context) encodeRecord(rec *ContextRecord) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(rec)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(),nil
}

func (ctx *Context) decodeRecord(data []byte) (*ContextRecord, error) {
	ctxRec := ContextRecord{}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&ctxRec)
	return &ctxRec,err
}


//func (ctx *Context) GetRecord(name string) (*ContextRecord,error) {
//	rec , ok := ctx.records[name]
//	if ok {
//		return &rec,nil
//	}
//	return &ContextRecord{},errors.New("Variable doesn't exist")
//}
//

