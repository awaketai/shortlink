package idgen

import "context"

// Generator 定义了短码生成器的能力
// 实现应尽可能保证生成的短码在一定概率下是唯一的，并符合业务对长度、字符集的要求
type Generator interface {
	// GenerateShortCode 为给定的输入(通常是长URL)，生成一个短码
	GenerateShortCode(ctx context.Context, input string) (string, error)
}
