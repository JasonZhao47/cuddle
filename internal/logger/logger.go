package logger

import "go.uber.org/zap"

// 这是一个没有意义的解耦
// 因为使用的时候同样要初始化这个Logger
// 所以依然耦合
var _ *zap.Logger
