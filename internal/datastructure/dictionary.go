package datastructure

import "time"

type Obj struct {
	Value interface{}
}

type Dict struct {
	dictStore        map[string]*Obj
	expiredDictStore map[string]uint64
}

func (d *Dict) Set(key string, value interface{}, expir uint64) {
	d.dictStore[key] = &Obj{Value: value}
	d.expiredDictStore[key] = expir
}

func (d *Dict) Get(key string) (Obj, bool) {
	obj, exist := d.dictStore[key]
	if !exist {
		return Obj{}, false
	}

	if expiredAt, hasExpired := d.expiredDictStore[key]; hasExpired {
		if uint64(time.Now().UnixMilli()) > expiredAt {
			delete(d.dictStore, key)
			delete(d.expiredDictStore, key)
			return Obj{}, false
		}
	}
	return *obj, true
}

func (d *Dict) Ttl(key string) (uint64, bool) {
	expir, exist := d.expiredDictStore[key]
	if !exist {
		return 0, false
	}
	return expir, true
}

func (d *Dict) Expire(key string, expr uint64) (int, bool) {
	_, ok := d.expiredDictStore[key]
	if !ok {
		return 0, false
	}
	d.expiredDictStore[key] = expr
	return 1, true
}

func (d *Dict) Del(keys []string) (int, bool) {
	cnt := 0
	for _, k := range keys {
		if _, ok := d.dictStore[k]; !ok {
			continue
		}
		delete(d.dictStore, k)
		delete(d.expiredDictStore, k)
		cnt++
	}
	return cnt, true
}

func (d *Dict) Exist(keys []string) (int, bool) {
	cnt := 0
	expiredList := make([]string, 0)
	for _, k := range keys {
		if _, ok := d.dictStore[k]; !ok {
			continue
		}
		if hadExpr, ok := d.expiredDictStore[k]; ok && hadExpr < uint64(time.Now().UnixMilli()) {
			expiredList = append(expiredList, k)
			continue
		}
		cnt++
	}
	d.Del(expiredList)
	return cnt, true
}
