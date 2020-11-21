/*
Navicat MySQL Data Transfer

Source Server         : MonDB
Source Server Version : 50556
Source Host           : 
Source Database       : mon

Target Server Type    : MYSQL
Target Server Version : 50556
File Encoding         : 65001

Date: 2018-10-29 10:34:01
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for t_telegraf_traceroute
-- ----------------------------
DROP TABLE IF EXISTS `t_telegraf_traceroute`;
CREATE TABLE `t_telegraf_traceroute` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `hostname` char(50) NOT NULL,
  `target` char(50) NOT NULL,
  `type` char(50) NOT NULL,
  `enable` int(11) NOT NULL,
  `update` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;
