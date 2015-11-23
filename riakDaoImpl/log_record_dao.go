package riakdaoimpl

import (
	riak "github.com/basho/riak-go-client"
)

type LogRecordDao struct {
	Cluster   riak.Cluster
	LogBucket string
}

func (dao *LogRecordDao) SaveLogRecord(logId string, logRecord string) error {
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

func (dao *LogRecordDao) GetLogRecord(logId string) (string, error) {
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

func (dao *LogRecordDao) DeleteLogRecord(logId string) error {
	cmd, _ := riak.NewDeleteValueCommandBuilder().
		WithBucket(dao.LogBucket).
		WithKey(logId).
		Build()
	return dao.Cluster.Execute(cmd)
}