package eventsourcing

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type DataStorage struct {
	data map[string]interface{}
}

func NewEmptyDataStorage() *DataStorage {
	return &DataStorage{data: make(map[string]interface{})}
}

func NewDataStorage(data map[string]interface{}) *DataStorage {
	return &DataStorage{data: data}
}

func (storage *DataStorage) String() string {
	return fmt.Sprintf("%v", storage.data)
}

func (storage *DataStorage) MarshalJSON() string {
	jsonBytes, _ := json.Marshal(storage.ToSimpleData())
	jsonString := string(jsonBytes)

	return jsonString
}

func (storage *DataStorage) UnmarshalJSON(jsonString string) {
	simpleData := make(map[string]interface{})
	json.Unmarshal([]byte(jsonString), &simpleData)

	storage.FromSimpleData(simpleData)
}

func (storage *DataStorage) ToSimpleData() *map[string]interface{} {
	simpleData := make(map[string]interface{})

	for key, value := range storage.data {
		valueType := strings.Trim(reflect.TypeOf(value).String(), "*") // cut pointer symbol

		if valueType == reflect.TypeOf(DataStorage{}).String() {
			dataStorage := value.(*DataStorage)
			simpleData[key] = dataStorage.ToSimpleData()
		} else {
			simpleData[key] = value
		}
	}

	return &simpleData
}

func (storage *DataStorage) FromSimpleData(simpleData map[string]interface{}) {
	storage.data = make(map[string]interface{})

	for key, value := range simpleData {
		valueType := strings.Trim(reflect.TypeOf(value).String(), "*") // cut pointer symbol

		if valueType == "map[string]interface {}" {
			dataStorage := NewEmptyDataStorage()
			dataStorage.FromSimpleData(value.(map[string]interface{}))

			storage.data[key] = dataStorage
		} else {
			storage.data[key] = value
		}
	}
}

func (storage *DataStorage) Set(key string, value interface{}) {
	pieces := strings.Split(key, "/")
	lastPiece, pieces := pieces[len(pieces)-1], pieces[:len(pieces)-1]
	storagePointer := storage

	for _, piece := range pieces {
		nestedStoragePointer, isExists := storagePointer.getStorage(piece)

		if !isExists {
			nestedStoragePointer = NewEmptyDataStorage()
			storagePointer.setStorage(piece, nestedStoragePointer)
		}

		storagePointer = nestedStoragePointer
	}

	storagePointer.data[lastPiece] = value
}

func (storage *DataStorage) Get(key string) interface{} {
	pieces := strings.Split(key, "/")
	lastPiece, pieces := pieces[len(pieces)-1], pieces[:len(pieces)-1]
	storagePointer := storage

	for _, piece := range pieces {
		nestedStoragePointer, isExists := storagePointer.getStorage(piece)

		if !isExists {
			return nil
		}
		storagePointer = nestedStoragePointer
	}

	return storagePointer.data[lastPiece]
}

func (storage *DataStorage) Delete(key string) {

}

func (storage *DataStorage) merge(key string, value interface{}) {

}

func (storage *DataStorage) setStorage(key string, nestedStorage *DataStorage) {
	storage.data[key] = nestedStorage
}

func (storage *DataStorage) getStorage(key string) (*DataStorage, bool) {
	if value, ok := storage.data[key]; ok {
		dataStorage := value.(*DataStorage)
		return dataStorage, true
	}

	return nil, false
}

//<?php

////    public static function merge(array $array1, array $array2)
////    {
////        $merged = $array1;
////
////        foreach ($array2 as $key => & $value) {
////            if (is_array($value) && isset($merged[$key]) && is_array($merged[$key])) {
////                $merged[$key] = $this->merge($merged[$key], $value);
////            } else {
////                if (is_numeric($key)) {
////                    if (!in_array($value, $merged)) {
////                        $merged[] = $value;
////                    }
////                } else {
////                    $merged[$key] = $value;
////                }
////            }
////        }
////
////        return $merged;
////    }
//
//    protected function setByPath($path, $value)
//    {
//        $pieces = explode('/', $path);
//        $lastPiece = array_pop($pieces);
//        $data = &$this->data;
//
//        foreach ($pieces as $piece) {
//            if (!isset($data[$piece]) || !is_array($data[$piece])) {
//                $data[$piece] = [];
//            }
//            $data = &$data[$piece];
//        }
//        $data[$lastPiece] = $value;
//    }
//
//    protected function getStorage($path = null)
//    {
//        $pieces = $path ? explode('/', $path) : array();
//        $data = &$this->data;
//
//        foreach ($pieces as $piece) {
//            if (!isset($data[$piece])) {
//                return null;
//            }
//            $data = &$data[$piece];
//        }
//
//        return $data;
//    }
//
//    protected function delByPath($path)
//    {
//        $pieces = explode('/', $path);
//        $key = array_pop($pieces);
//        $data = &$this->data;
//
//        foreach ($pieces as $piece) {
//            if (!isset($data[$piece])) {
//                return;
//            }
//            $data = &$data[$piece];
//        }
//
//        unset($data[$key]);
//    }
//}
