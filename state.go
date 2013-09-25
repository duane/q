package main

import (
  "labix.org/v2/mgo/bson"
)

func (state *Queue) Save() (data []byte, err error) {
  l := state.List
  array := []interface{}{}
  for e := l.Front(); e != nil; e = e.Next() {
    array = append(array, *e.Value.(*interface{}))
  }
  data, err = bson.Marshal(bson.M{"Queue": array})
  return
}

func (state *Queue) Recovery(data []byte) (err error) {
  var unmarshalled bson.M
  err = bson.Unmarshal(data, &unmarshalled)
  if err != nil {
    return
  }
  array := unmarshalled["Queue"].([]interface{})

  for _, e := range array {
    state.Enqueue(e)
  }
  return
}
