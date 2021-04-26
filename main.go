package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strings"
)
var ctx = context.TODO()

func main() {
	args := getInput()
	args.validateArgs()

	if *args.cluster {
		fmt.Printf("Connected to cluster! \n")
		clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: strings.Split(*args.hosts, ";"),
		})
		deleteKeyCluster(*clusterClient, *args.key)
		return
	}

	redisClient := redis.NewClient(&redis.Options{Addr: *args.hosts})
	deleteKey(*redisClient, *args.key)
}

func deleteKeyCluster(clusterClient redis.ClusterClient, key string) {
	err := clusterClient.ForEachMaster(ctx, func(ctx context.Context, client *redis.Client) error {
		fmt.Printf("Connected to the %s node \n", client.Options().Addr)
		deleteKey(*client, key)
		return nil
	})
	if err != nil {
		panic (err)
	}
}

func deleteKey(client redis.Client, key string) {
	var cursor uint64

	iter := client.Scan(context.TODO(), cursor, key, 0).Iterator()

	logCount := 0
	for iter.Next(ctx) {
		client.Unlink(ctx, iter.Val())
		if logCount > 100 {
			logCount = 0
			fmt.Printf("Info: Key %s deleted \n", iter.Val())
		}
		logCount++
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
}

type Args struct {
	cluster *bool
	hosts *string
	key *string
}
func (a Args) validateArgs() {
	if *a.hosts == "" {
		panic("You must inform a host to connect.")
	}
	if *a.key == "" {
		panic("You must inform a key to delete.")
	}
}

func getInput() Args {
	args := Args{
		cluster: flag.Bool("cluster", false, "Connect to a cluster?"),
		hosts: flag.String("hosts", "", "Redis host. (For cluster, separate with ;)"),
		key: flag.String("key", "", "Key to delete"),
	}
	flag.Parse()

	return args
}