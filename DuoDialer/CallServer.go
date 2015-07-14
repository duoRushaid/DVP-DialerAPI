package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func RegisterCallServer(serverId string) CallServerInfo {
	defCallServerInfo := CallServerInfo{}
	if serverId == "*" {
		//pick callserver
	}
	//Get CallServer info
	cs := CallServerInfo{}
	cs.CallServerId = serverId
	cs.MaxChannelCount = 100
	cs.Url = callServer

	callServerKey := fmt.Sprintf("CallServer:%s", serverId)
	callServerjson, _ := json.Marshal(cs)
	addResult := RedisAdd(callServerKey, string(callServerjson))

	if addResult == "OK" {
		csck := fmt.Sprintf("CallServerConcurrentCalls:%s", serverId)
		csmcl := fmt.Sprintf("CallServerMaxCallLimit:%s", serverId)
		RedisSet(csck, "0")
		RedisIncrBy(csmcl, cs.MaxChannelCount)
		return cs
	}
	return defCallServerInfo
}

func GetCallServerInfo(serverId string) CallServerInfo {
	callServerKey := fmt.Sprintf("CallServer:%s", serverId)
	csString := RedisGet(callServerKey)
	if csString != "" {
		var callServerInfo CallServerInfo
		json.Unmarshal([]byte(csString), &callServerInfo)
		return callServerInfo
	} else {
		return RegisterCallServer(serverId)
	}
}

func GetConcurrentChannelCount(serverId, campaignId string) (concurrentOnServer, concurrentOnCamp int) {
	csckC := fmt.Sprintf("CallServerConcurrentCalls:%s:%s", serverId, campaignId)
	csck := fmt.Sprintf("CallServerConcurrentCalls:%s", serverId)
	channelCountC := RedisGet(csckC)
	fmt.Println("RedisGet channelCountC: ", channelCountC)

	if channelCountC == "" {
		RedisSet(csckC, "0")
		channelCountC = "0"
	}

	channelCountS := RedisGet(csck)
	fmt.Println("RedisGet channelCountS: ", channelCountS)
	valueC, err := strconv.Atoi(channelCountC)
	valueS, _ := strconv.Atoi(channelCountS)
	if err != nil {
		fmt.Println(err.Error())
		return 0, 0
	} else {
		return valueS, valueC
	}
}

func IncrConcurrentChannelCount(serverId, campaignId string) {
	csckC := fmt.Sprintf("CallServerConcurrentCalls:%s:%s", serverId, campaignId)
	csck := fmt.Sprintf("CallServerConcurrentCalls:%s", serverId)
	RedisIncr(csckC)
	RedisIncr(csck)
}

func DecrConcurrentChannelCount(serverId, campaignId string) {
	csckC := fmt.Sprintf("CallServerConcurrentCalls:%s:%s", serverId, campaignId)
	csck := fmt.Sprintf("CallServerConcurrentCalls:%s", serverId)
	decValueC := RedisIncrBy(csckC, -1)
	RedisIncrBy(csck, -1)
	if decValueC < 0 {
		RedisSet(csck, "0")
	}
}

func IncrMaxLimit(serverId string) {
	callServerKey := fmt.Sprintf("CallServer:%s", serverId)
	csString := RedisGet(callServerKey)
	if csString == "" {
		RegisterCallServer(serverId)
	}

	csmcl := fmt.Sprintf("CallServerMaxCallLimit:%s", serverId)
	RedisIncr(csmcl)
}

func DecrMaxLimit(serverId string) {
	csmcl := fmt.Sprintf("CallServerMaxCallLimit:%s", serverId)
	decValue := RedisIncrBy(csmcl, -1)
	if decValue < 0 {
		RedisSet(csmcl, "0")
	}
}

func GetMaxChannelLimit(serverId string) int {
	csmcl := fmt.Sprintf("CallServerMaxCallLimit:%s", serverId)
	maxChannelCount := RedisGet(csmcl)
	value, err := strconv.Atoi(maxChannelCount)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	} else {
		return value
	}
}