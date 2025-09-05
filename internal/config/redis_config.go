/*

The redis config for the  connection
*/

package config

type RedisDBConfig struct {
	Host string
	Port int
}

var RedisDefaultConfig = RedisDBConfig{
	Host: "localhost",
	Port: 6379,
}

var UsedHashSet = "used_hash_set"

var UniqueHashSet = "frest_hash_set"

//unique_hash_set := "frest_hash_set"
