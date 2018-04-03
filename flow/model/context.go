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
	"github.com/pkg/errors"
)

type Variable struct {
	Value     interface{}
	ValueType string
}

func (vrbl * Variable) IsNumber() bool {
	if vrbl.ValueType == "int" || vrbl.ValueType == "float" {
		return true
	}else {
		return false
	}
}

func (vrbl *Variable) IsEqual(var2 *Variable) (bool,error) {
	if vrbl.ValueType == var2.ValueType {
		switch vrbl.ValueType {
		case "string":
			v1,ok1 := vrbl.Value.(string)
			v2,ok2 := var2.Value.(string)
			if ok1 && ok2 {
				return v1==v2,nil
			}else {
				return false , errors.New("Can't cast var to string")
			}
		case "int","float":
			v1,ok1 := vrbl.ToNumber()
			v2,ok2 := var2.ToNumber()
			if ok1 == nil && ok2==nil {
				return v1==v2,nil
			}else {
				return false , errors.New("Can't cast var to number")
			}
		case "bool":
			v1,ok1 := vrbl.Value.(bool)
			v2,ok2 := var2.Value.(bool)
			if ok1 && ok2 {
				return v1==v2,nil
			}else {
				return false , errors.New("Can't cast var to bool")
			}

		}
	}else {
		return false , errors.New("Types are different")
	}
	return false , nil
}

func (vrbl * Variable)ToNumber()(float64,error) {
	switch v := vrbl.Value.(type) {
	case int :
		return float64(v),nil
	case int32 :
		return float64(v),nil
	case int64 :
		return float64(v),nil
	case float32 :
		return float64(v),nil
	case float64 :
		return float64(v),nil
	default:
		return 0 , errors.New("Can't convert into float")
	}
	return 0 , errors.New("Not numeric value type")
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
	gob.Register(map[string]interface{}{})
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

func (ctx *Context) DeleteRecord(name string,flowId string,inMemory bool ) error {
	if inMemory {
		//ctx.inMemoryStore[flowId] = *rec
	} else {
		err := ctx.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(flowId))
			return b.Delete([]byte(name))
		})
		return err
	}
	return nil
}


func (ctx *Context) GetVariable(name string,flowId string) (Variable,error) {
	rec , err := ctx.GetRecord(name,flowId)
	if err == nil {
		return rec.Variable , err
	}else {
		return Variable{},err
	}

}

func (ctx *Context) GetVariableType(name string ,flowId string) (string,error) {
	varb,err := ctx.GetVariable(name,flowId)
	if err == nil {
		return varb.ValueType , err
	}else {
		return "",err
	}
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
			if b == nil {
				err = errors.New("Flow doesn't exist")
				return nil
			}
			data := b.Get([]byte(name))
		    if data == nil {
				err = errors.New("Not Found")
			}else {
				ctxRec,err = ctx.decodeRecord(data)
			}
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
		if b == nil {
			return nil
		}
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

