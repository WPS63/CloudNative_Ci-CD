//Damian Hupka
//ECGR-4090: Cloud Native Application Architecture Spring 20201
//02/15/2021
//Lab #2 - Interfaces (Implementing get() and put() of LRUCache)
package lru

import ( 
	"errors"
	"fmt"
)

type Cacher interface {
	Get(interface{}) (interface{}, error)
	Put(interface{}, interface{}) error
}

type lruCache struct {
	size      int
	remaining int
	cache     map[string]string
	queue     []string
}

func NewCache(size int) Cacher {
	return &lruCache{size: size, remaining: size, cache: make(map[string]string), queue: make([]string, size)}
}


//We have to do this remaining stuff because append _, _, _, key, key, key
func (lru *lruCache) Get(key interface{}) (interface{}, error) {
    //Looking if key is found, returns key or -1 if not found
    if _, ok := lru.cache[key.(string)]; ok {
        if lru.remaining == 0 {
            lru.queue[lru.size-1] = key.(string) //add to queue as highest available index
        } else {
            lru.queue[lru.size-lru.remaining] = key.(string) //add to queue as current index
        }
        return lru.cache[key.(string)], nil //return the value at the given key
    }

    return "-1", errors.New("error")
}

func (lru *lruCache) Put(key, val interface{}) error {
    if lru.remaining < 0 {
        return errors.New("shoulnd't be here")
    }

    //If size = remaining, deletes queue element 0 from cache and queue, adds 
    if lru.remaining == 0 {
        delete(lru.cache, lru.queue[0])   //delete the LRU from the cache
        lru.qDel(lru.queue[0])            //delete the LRU(head) from queue (which reduces queue slice size by one)
        lru.queue = append(lru.queue, "") //append empty string to queue to make it the original size again
        fmt.Print("Empty queue index amended: ")
        fmt.Println(lru.queue)
        lru.queue[lru.size-1] = key.(string) //now add the key to the tail of queue
        fmt.Print("Now the queue is: ")
        fmt.Println(lru.queue)
    } else {
        lru.queue[lru.size-lru.remaining] = key.(string) //if capacity isn't max, just add to slice
    }
    if lru.remaining > 0 { //Doing this so wont reach the error throwing, ensures non-negative here
        lru.remaining--
    }
    lru.cache[key.(string)] = val.(string) //insert into cache

    return nil
}


// Delete element from queue
func (lru *lruCache) qDel(ele string) {
	for i := 0; i < len(lru.queue); i++ {
		if lru.queue[i] == ele {
			oldlen := len(lru.queue)
			copy(lru.queue[i:], lru.queue[i+1:])
			lru.queue = lru.queue[:oldlen-1]
			break
		}
	}
}
