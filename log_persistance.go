package riaktimeline

import (
	riak "github.com/basho/riak-go-client"
)

func storeLog(logRecord string) (string, error) {
	value := &riak.Object{
		ContentType:     "text/plain",
		Charset:         "utf-8",
		ContentEncoding: "utf-8",
		Value:           []byte(logRecord),
	}

	cmd, _ := riak.NewStoreValueCommandBuilder().
		WithBucket(Config.LogBucket).
		WithContent(value).
		Build()
	if err := RiakCluster.Execute(cmd); err != nil {
		logErr("Saving log record ", err)
		return "", err
	}
	svc := cmd.(*riak.StoreValueCommand)
	return svc.Response.GeneratedKey, nil
}

func fetchLog(logId string) (string, error) {
	cmd, _ := riak.NewFetchValueCommandBuilder().
		WithBucket(Config.LogBucket).
		WithNotFoundOk(true).
		WithKey(logId).
		Build()
	if err := RiakCluster.Execute(cmd); err != nil {
		logErr("Fetching log record for "+logId, err)
		return "", err
	}
	fvc := cmd.(*riak.FetchValueCommand)
	if fvc.Response == nil || fvc.Response.IsNotFound || len(fvc.Response.Values) == 0 {
		return "", nil
	}
	return string(fvc.Response.Values[0].Value), nil
}

func deleteLog(logId string) error {
	cmd, _ := riak.NewDeleteValueCommandBuilder().
		WithBucket(Config.LogBucket).
		WithKey(logId).
		Build()
	if err := RiakCluster.Execute(cmd); err != nil {
		logErr("Deleting log record for "+logId, err)
		return err
	}
	return nil
}
