package internal

import "google.golang.org/grpc"

// WithUnaryServerInterceptors 用于扩展链式一元拦截器。
// 第一个拦截器在最外层，最后一个拦截器是包围实际调用的最内部的包装器。
func WithUnaryServerInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(interceptors...)
}

// WithStreamServerInterceptors 用于扩展链式流拦截器。
// 第一个拦截器在最外层，最后一个拦截器是包围实际调用的最内部的包装器。
func WithStreamServerInterceptors(interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	return grpc.ChainStreamInterceptor(interceptors...)
}
