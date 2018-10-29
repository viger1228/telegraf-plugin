/*
Navicat MySQL Data Transfer

Source Server         : wz-MonDB
Source Server Version : 50556
Source Host           : 119.28.66.203:3306
Source Database       : mon

Target Server Type    : MYSQL
Target Server Version : 50556
File Encoding         : 65001

Date: 2018-10-29 08:10:28
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for t_telegraf_tcping
-- ----------------------------
DROP TABLE IF EXISTS `t_telegraf_tcping`;
CREATE TABLE `t_telegraf_tcping` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `hostname` char(50) DEFAULT NULL,
  `target` char(50) DEFAULT NULL,
  `port` int(11) DEFAULT NULL,
  `enable` int(11) DEFAULT NULL,
  `update` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8;
