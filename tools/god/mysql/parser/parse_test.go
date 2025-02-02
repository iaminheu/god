package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParsePlainText(t *testing.T) {
	_, err := Parse("plain text")
	assert.NotNil(t, err)
}

func TestParseCreateTable(t *testing.T) {
	table, err := Parse("CREATE TABLE `user_snake` (\n  `id` bigint(10) NOT NULL AUTO_INCREMENT,\n  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名称',\n  `password` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户密码',\n  `mobile` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机号',\n  `gender` char(5) COLLATE utf8mb4_general_ci NOT NULL COMMENT '男｜女｜未公开',\n  `nickname` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '用户昵称',\n  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,\n  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,\n  PRIMARY KEY (`id`),\n  UNIQUE KEY `name_index` (`name`),\n  KEY `mobile_index` (`mobile`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;")
	assert.Nil(t, err)
	assert.Equal(t, "user_snake", table.Name.Source())
	assert.Equal(t, "id", table.PrimaryKey.Name.Source())
	assert.Equal(t, true, table.ContainsTime())
}

func TestParseSelect(t *testing.T) {
	_, err := Parse("select * from user")
	assert.Equal(t, unSupportDDL, err)
}
