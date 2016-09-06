package handler

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/zhutingle/gotrix/global"
)

type handleRedis struct {
	redisClient *redis.Pool
}

func (this *handleRedis) init() *handleRedis {
	// 建立连接池
	this.redisClient = &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle:     global.Config.Redis.MaxIdle,
		MaxActive:   global.Config.Redis.MaxActive,
		IdleTimeout: global.Config.Redis.IdleTimeout * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", global.Config.Redis.Ip+":"+global.Config.Redis.Host, redis.DialPassword(global.Config.Redis.Pass))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
	return this
}

func (this *handleRedis) handle(job *Job, cp *global.CheckedParams) (result interface{}, gErr *global.GotrixError) {

	funcStrs := funcReg.FindAllStringSubmatch(job.Job, -1)
	funcName := funcStrs[0][1]
	funcPara := funcStrs[0][2]

	args := make([]interface{}, 0)
	strs := argsReg.FindAllStringSubmatch(funcPara, -1)
	for j := 0; j < len(strs); j++ {
		str := strs[j][1]
		if str[0] == '$' {
			value := cp.V[str[2:len(str)-1]]
			switch value.(type) {
			case string:
				args = append(args, value.(string))
				break
			case float64:
				args = append(args, strconv.FormatFloat(value.(float64), 'g', 'e', 64))
				break
			case int64:
				args = append(args, strconv.FormatInt(value.(int64), 10))
				break
			case bool:
				args = append(args, strconv.FormatBool(value.(bool)))
				break
			case map[string]interface{}:
				json, err := json.Marshal(value)
				if err != nil {
					gErr = global.REDIS_HANDLE_JSON_ERROR
					return
				}
				args = append(args, string(json))
				break
			case []interface{}:
				json, err := json.Marshal(value)
				if err != nil {
					gErr = global.REDIS_HANDLE_JSON_ERROR
					return
				}
				args = append(args, string(json))
				break
			default:
				break
			}
		} else if str[0] == '\'' {
			args = append(args, str[1:len(str)-1])
		} else {
			args = append(args, str)
		}
	}

	c := this.redisClient.Get()
	defer c.Close()

	v, err := c.Do(funcName, args...)
	if err != nil {
		gErr = global.REDIS_CONNECT_ERROR
		return
	}
	switch v.(type) {
	case []byte:
		var vJson map[string]interface{} = make(map[string]interface{})
		err = json.Unmarshal(v.([]byte), &vJson)
		if err == nil {
			result = vJson
		} else {
			result = string(v.([]byte))
		}
		break
	case int64:
		result = strconv.Itoa(int(v.(int64)))
		break
	case redis.Error:
		log.Printf("执行时[%s]出现异常[%s]", funcName, v.(redis.Error))
		log.Println(args)
		gErr = global.REDIS_EXEC_ERROR
		return
	default:
		break
	}
	return
}
