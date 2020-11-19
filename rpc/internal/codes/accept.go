package codes

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Acceptable 根据rpc状态码判断错误是否可接受
func Acceptable(err error) bool {
	switch status.Code(err) {
	case codes.DeadlineExceeded, codes.Internal, codes.Unavailable, codes.DataLoss:
		return false
	default:
		return true
	}
}
