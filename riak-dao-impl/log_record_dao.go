package riakdaoimpl

import (
	riak "github.com/basho/riak-go-client"
)

type LogRecordRiakDao struct {
	Cluster   *riak.Cluster
	LogBucket string
}

func (dao *LogRecordRiakDao) SaveLogRecord(logId string, logRecord string) error {
	value := &riak.Object{
		ContentType:     "text/plain",
		Charset:         "utf-8",
		ContentEncoding: "utf-8",
		Value:           []byte(logRecord),
	}

	cmd, _ := riak.NewStoreValueCommandBuilder().
		WithKey(logId).
		WithBucket(dao.LogBucket).
		WithContent(value).
		Build()
	return dao.Cluster.Execute(cmd)
}

func (dao *LogRecordRiakDao) GetLogRecord(logId string) (string, error) {
	cmd, _ := riak.NewFetchValueCommandBuilder().
		WithBucket(dao.LogBucket).
		WithNotFoundOk(true).
		WithKey(logId).
		Build()
	if err := dao.Cluster.Execute(cmd); err != nil {
		logErr("Fetching log record for "+logId, err)
		return "", err
	}
	fvc := cmd.(*riak.FetchValueCommand)
	if fvc.Response == nil || fvc.Response.IsNotFound || len(fvc.Response.Values) == 0 {
		return "", nil
	}
	return string(fvc.Response.Values[0].Value), nil
}

func (dao *LogRecordRiakDao) DeleteLogRecord(logId string) error {
	cmd, _ := riak.NewDeleteValueCommandBuilder().
		WithBucket(dao.LogBucket).
		WithKey(logId).
		Build()
	return dao.Cluster.Execute(cmd)
}

func (dao *LogRecordRiakDao) Ping() (bool, error) {
	client, error := riak.NewClient(&riak.NewClientOptions{dao.Cluster, 0, nil})
	if error != nil {
		return false, error
	}
	return client.Ping()
}

func NewLogRecordRiakDao(cluster *riak.Cluster, logBucket string) *LogRecordRiakDao {
	return &LogRecordRiakDao{
		Cluster:   cluster,
		LogBucket: logBucket,
	}
}
