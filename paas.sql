-- -------------------------------------------------------------
-- TablePlus 6.2.1(578)
--
-- https://tableplus.com/
--
-- Database: paas
-- Generation Time: 2025-02-25 10:42:49.6420
-- -------------------------------------------------------------


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


DROP TABLE IF EXISTS `developer`;
CREATE TABLE `developer` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `uri` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '开发者编码',
  `user_id` bigint NOT NULL DEFAULT '0' COMMENT '用户id',
  `state` tinyint NOT NULL DEFAULT '1' COMMENT '开发者状态1未认证2已认证3认证失败',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `file`;
CREATE TABLE `file` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `uri` varchar(64) NOT NULL DEFAULT '',
  `url` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'url',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '文件名',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `game`;
CREATE TABLE `game` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `uri` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '游戏编码',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '游戏名称',
  `description` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '游戏描述',
  `logo` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '游戏logo',
  `user_id` bigint NOT NULL DEFAULT '0' COMMENT '用户id',
  `developer_id` bigint NOT NULL DEFAULT '0' COMMENT '开发者id',
  `app_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'app_id',
  `secret_key` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'app_secret',
  `runtime_config` varchar(1024) NOT NULL DEFAULT '' COMMENT '运行时配置',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_name` (`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `game_version`;
CREATE TABLE `game_version` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `uri` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '版本编码',
  `game_uri` varchar(32) NOT NULL DEFAULT '' COMMENT '游戏URI',
  `game_id` int NOT NULL DEFAULT '0' COMMENT '游戏id',
  `user_id` bigint NOT NULL DEFAULT '0' COMMENT '用户id',
  `developer_id` bigint NOT NULL DEFAULT '0' COMMENT '开发者id',
  `version` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '版本号',
  `change_log` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '变更日志',
  `file_url` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '游戏包文件地址',
  `backend_image_url` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '后端镜像地址',
  `socket_image_url` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'socket镜像地址',
  `frontend_image_url` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '前端镜像地址',
  `check_state` tinyint NOT NULL DEFAULT '1' COMMENT '审核状态1未提交审核2审核中3审核成功4审核失败',
  `check_failed_reason` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '审核失败原因',
  `deploy_state` tinyint NOT NULL DEFAULT '1' COMMENT '部署状态1未提交部署2部署中3部署成功4部署失败',
  `build_config` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '打包镜像配置',
  `deployed_at` datetime DEFAULT NULL COMMENT '部署时间',
  `checked_at` datetime DEFAULT NULL COMMENT '审核时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_game_uri_version` (`game_uri`,`version`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `user_account`;
CREATE TABLE `user_account` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` bigint NOT NULL DEFAULT '0' COMMENT '用户id',
  `openid` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'openid',
  `platform_id` tinyint NOT NULL DEFAULT '1' COMMENT '注册来源1手机号',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `user_private_info`;
CREATE TABLE `user_private_info` (
  `user_id` bigint NOT NULL DEFAULT '0' COMMENT '用户id',
  `tel` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '手机号',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `user_profile`;
CREATE TABLE `user_profile` (
  `user_id` bigint NOT NULL DEFAULT '0' COMMENT '用户id',
  `uri` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '用户编码',
  `nickname` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '昵称',
  `headimgurl` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '头像',
  `sex` tinyint NOT NULL DEFAULT '1' COMMENT '性别',
  `region` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '区域',
  `country` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '国家',
  `province` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '省份',
  `city` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '城市',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `file` (`id`, `uri`, `url`, `name`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1, '2502231948cerNiu', 'https://ikit-blog.oss-cn-hangzhou.aliyuncs.com/uploads/2025/02/23/2502231948cerNiu.pdf', '游戏服务时序图.pdf', '2025-02-23 19:48:43', '2025-02-23 19:48:43', NULL),
(2, '2502231950TIYHgA', '/uploads/2025/02/23/2502231950TIYHgA.pdf', '游戏服务时序图.pdf', '2025-02-23 19:50:51', '2025-02-23 19:50:51', NULL),
(3, '2502240930kr9qSF', '/uploads/2025/02/24/2502240930kr9qSF.png', 'WX20250125-202622.png', '2025-02-24 09:30:03', '2025-02-24 09:30:03', NULL),
(4, '25022409318bvD0B', '/uploads/2025/02/24/25022409318bvD0B.jpg', 'vladislav_klapin_pvr9_gsj93_pc_unsplash_1_wps9dzapaa.jpg', '2025-02-24 09:31:11', '2025-02-24 09:31:11', NULL),
(5, '2502241511s3tyf6', '/uploads/2025/02/24/2502241511s3tyf6.png', '2502240930kr9qSF.png', '2025-02-24 15:11:01', '2025-02-24 15:11:01', NULL),
(6, '2502241512i351Br', '/uploads/2025/02/24/2502241512i351Br.png', '2502240930kr9qSF.png', '2025-02-24 15:12:05', '2025-02-24 15:12:05', NULL),
(7, '2502241521pbsiMJ', '/uploads/2025/02/24/2502241521pbsiMJ.png', '2502240930kr9qSF.png', '2025-02-24 15:21:42', '2025-02-24 15:21:42', NULL),
(8, '2502241522QV0ZIJ', '/uploads/2025/02/24/2502241522QV0ZIJ.png', '2502240930kr9qSF.png', '2025-02-24 15:22:14', '2025-02-24 15:22:14', NULL),
(9, '2502241523gvEvaP', '/uploads/2025/02/24/2502241523gvEvaP.png', 'WX20250125-202622.png', '2025-02-24 15:23:52', '2025-02-24 15:23:52', NULL),
(10, '2502241524oBK7yh', '/uploads/2025/02/24/2502241524oBK7yh.png', '2502240930kr9qSF.png', '2025-02-24 15:24:29', '2025-02-24 15:24:29', NULL),
(11, '2502241533OZFOiz', '/uploads/2025/02/24/2502241533OZFOiz.png', '2502240930kr9qSF.png', '2025-02-24 15:33:58', '2025-02-24 15:33:58', NULL),
(12, '25022416562nu1kp', '/uploads/2025/02/24/25022416562nu1kp.png', 'WX20250121-094105.png', '2025-02-24 16:56:10', '2025-02-24 16:56:10', NULL),
(13, '2502241700XDI7Su', '/uploads/2025/02/24/2502241700XDI7Su.gif', 'node.gif', '2025-02-24 17:00:09', '2025-02-24 17:00:09', NULL),
(14, '2502241700gwwSiQ', '/uploads/2025/02/24/2502241700gwwSiQ.png', 'WX20250121-094105.png', '2025-02-24 17:00:32', '2025-02-24 17:00:32', NULL),
(15, '25022417018nh1yC', '/uploads/2025/02/24/25022417018nh1yC.png', 'WX20250121-094143.png', '2025-02-24 17:01:52', '2025-02-24 17:01:52', NULL),
(16, '2502241703DkpEGc', '/uploads/2025/02/24/2502241703DkpEGc.png', '2502240930kr9qSF.png', '2025-02-24 17:03:40', '2025-02-24 17:03:40', NULL),
(17, '2502241704xFIEYa', '/uploads/2025/02/24/2502241704xFIEYa.png', 'WX20250121-094143.png', '2025-02-24 17:04:01', '2025-02-24 17:04:01', NULL),
(18, '2502241723cT79cE', '/uploads/2025/02/24/2502241723cT79cE.png', '2502240930kr9qSF.png', '2025-02-24 17:23:28', '2025-02-24 17:23:28', NULL),
(19, '2502250917garbz2', '/uploads/2025/02/25/2502250917garbz2.zip', '代码交付.zip', '2025-02-25 09:17:24', '2025-02-25 09:17:24', NULL),
(20, '2502250932YO2czx', '/uploads/2025/02/25/2502250932YO2czx.zip', '代码交付.zip', '2025-02-25 09:32:21', '2025-02-25 09:32:21', NULL),
(21, '2502250934xYKugd', '/uploads/2025/02/25/2502250934xYKugd.zip', '评分卡模型说明及代码例子.zip', '2025-02-25 09:34:04', '2025-02-25 09:34:04', NULL),
(22, '2502250936B2w2T3', '/uploads/2025/02/25/2502250936B2w2T3.zip', '2i2jxv6u81hlwaq-profile4.certimate.fun.zip', '2025-02-25 09:36:49', '2025-02-25 09:36:49', NULL),
(23, '2502250938eFvV9o', '/uploads/2025/02/25/2502250938eFvV9o.zip', '评分卡模型说明及代码例子.zip', '2025-02-25 09:38:40', '2025-02-25 09:38:40', NULL),
(24, '2502250940Ao8l77', '/uploads/2025/02/25/2502250940Ao8l77.zip', '评分卡模型说明及代码例子.zip', '2025-02-25 09:40:06', '2025-02-25 09:40:06', NULL),
(25, '2502250955rzIvf2', '/uploads/2025/02/25/2502250955rzIvf2.zip', '评分卡模型说明及代码例子.zip', '2025-02-25 09:55:50', '2025-02-25 09:55:50', NULL),
(26, '2502251008e8ZtA4', '/uploads/2025/02/25/2502251008e8ZtA4.zip', '评分卡模型说明及代码例子.zip', '2025-02-25 10:08:48', '2025-02-25 10:08:48', NULL),
(27, '2502251018QtL4mE', '/uploads/2025/02/25/2502251018QtL4mE.zip', '评分卡模型说明及代码例子.zip', '2025-02-25 10:18:04', '2025-02-25 10:18:04', NULL);

INSERT INTO `game` (`id`, `uri`, `name`, `description`, `logo`, `user_id`, `developer_id`, `app_id`, `secret_key`, `runtime_config`, `created_at`, `updated_at`, `deleted_at`) VALUES
(4, '2502211130Da3tA6', '', '', '/uploads/2025/02/24/25022409318bvD0B.jpg', 1, 1, 'yu3a0da6daa5ff5bba2ad', 'mzedopm05noY41wqxuHRUyimWIM5G339', '', '2025-02-21 11:30:10', '2025-02-24 07:54:41', '2025-02-24 15:54:41'),
(8, '2502211315oQaryn', '哈哈哈', 'sssssss', '/uploads/2025/02/24/25022409318bvD0B.jpg', 1, 1, 'yu333747b54670a1b6e6e', 'l7Ci9mTKG0GvJs4QkJjHNeg0w3GGpRYY', '', '2025-02-21 13:15:31', '2025-02-24 01:34:24', '2025-02-21 15:22:59'),
(12, '2502211657OkslEo', '测试3', '', '/uploads/2025/02/24/25022409318bvD0B.jpg', 1, 1, 'yu316d9937f82009ececb', '2SocSvtM2TEuIU5LBMTMLsl60y6Jc6h4', '', '2025-02-21 16:57:49', '2025-02-24 09:45:25', '2025-02-24 17:45:26'),
(13, '25022116583JLh1L', '测试4', '这是一款非常拉风的游戏', '/uploads/2025/02/24/25022409318bvD0B.jpg', 1, 1, 'yu36d47573f1b9a923c8a', '4BEKUECJMn1aD6pp4356G9DVMF0ZRHt3', '', '2025-02-21 16:58:04', '2025-02-24 17:06:03', NULL),
(14, '2502241535Q8DTiJ', '测试6', '', '', 1, 1, 'yu3f414fc7d03363fcbea', 'RKkSMJlbcYzj5ZeSkNfUbG2zxzvGTP3t', '', '2025-02-24 15:35:28', '2025-02-24 07:54:45', '2025-02-24 15:54:46'),
(15, '2502241538MfcnJg', 'test5', '1234234234', '/uploads/2025/02/24/2502241723cT79cE.png', 1, 1, 'yu31e3491a7a0c5ca0e87', '91hWgo760XZfBiIoiHC8qLGDs51VhAsZ', '', '2025-02-24 15:38:35', '2025-02-24 17:23:30', NULL);

INSERT INTO `game_version` (`id`, `uri`, `game_uri`, `game_id`, `user_id`, `developer_id`, `version`, `change_log`, `file_url`, `backend_image_url`, `socket_image_url`, `frontend_image_url`, `check_state`, `check_failed_reason`, `deploy_state`, `build_config`, `deployed_at`, `checked_at`, `created_at`, `updated_at`, `deleted_at`) VALUES
(2, '2502231729F3Vu8q', '25022116583JLh1L', 13, 1, 1, '1.0.0', '第一个版本1', 'http://baidu.com', '', '', '', 1, '', 1, 'null', NULL, NULL, '2025-02-23 17:29:43', '2025-02-23 09:42:18', '2025-02-23 17:42:18'),
(3, '25022317313G1Nam', '25022116583JLh1L', 13, 1, 1, '1.0.1', '第一个版本', 'http://baidu.com', 'xxxxx', 'xxxxx', 'xxxxxx', 1, '', 1, 'null', NULL, NULL, '2025-02-23 17:31:56', '2025-02-24 13:40:16', NULL),
(4, '25022317312yg26f', '25022116583JLh1L', 13, 1, 1, '1.0.2', '第一个版本', 'http://baidu.com', '', '', '', 1, '', 1, 'null', NULL, NULL, '2025-02-23 17:31:58', '2025-02-23 09:33:42', NULL),
(6, '2502250940vzpPTG', '2502241538MfcnJg', 15, 1, 1, '1.0.1', '第一个版本', 'https://ikit-blog.oss-cn-hangzhou.aliyuncs.com/uploads/2025/02/25/2502250940Ao8l77.zip', '', '', '', 1, '', 1, '{\"backend\":{\"workDir\":\"/app/backend\",\"cmd\":\"npm start\"},\"frontend\":{\"workDir\":\"/app/frontent\",\"cmd\":\"npm start\"},\"socket\":{\"workDir\":\"/app/server\",\"cmd\":\"npm start\"}}', NULL, NULL, '2025-02-25 09:40:21', '2025-02-25 02:07:29', '2025-02-25 10:07:29'),
(7, '2502250956gNVc5y', '2502241538MfcnJg', 15, 1, 1, '1.0.2', 'ssssss', 'https://ikit-blog.oss-cn-hangzhou.aliyuncs.com/uploads/2025/02/25/2502250955rzIvf2.zip', '', '', '', 1, '', 1, '{\"backend\":{\"workDir\":\"/backend\",\"cmd\":\"npm start\"},\"frontend\":{\"workDir\":\"/frontend\",\"cmd\":\"npm start\"},\"socket\":{\"workDir\":\"/socket\",\"cmd\":\"npm start\"}}', NULL, NULL, '2025-02-25 09:56:16', '2025-02-25 02:09:26', '2025-02-25 10:09:27'),
(8, '2502251009geDCyW', '2502241538MfcnJg', 15, 1, 1, '1.0.3', '第3个版本', 'https://ikit-blog.oss-cn-hangzhou.aliyuncs.com/uploads/2025/02/25/2502251008e8ZtA4.zip', '', '', '', 1, '', 1, '{\"backend\":{\"workDir\":\"/backend\",\"cmd\":\"npm start\"},\"frontend\":{\"workDir\":\"/frontend\",\"cmd\":\"npm start\"},\"socket\":{\"workDir\":\"/socket\",\"cmd\":\"npm start\"}}', NULL, NULL, '2025-02-25 10:09:17', '2025-02-25 02:18:43', '2025-02-25 10:18:44'),
(10, '2502251018sKJcBM', '2502241538MfcnJg', 15, 1, 1, '1.0.4', '第4个版本', 'https://ikit-blog.oss-cn-hangzhou.aliyuncs.com/uploads/2025/02/25/2502251018QtL4mE.zip', '', '', '', 1, '', 1, '{\"backend\":{\"workDir\":\"/backend\",\"cmd\":\"npm run start\"},\"frontend\":{\"workDir\":\"/frontend\",\"cmd\":\"npm run start\"},\"socket\":{\"workDir\":\"/socket\",\"cmd\":\"npm run start\"}}', NULL, NULL, '2025-02-25 10:18:38', '2025-02-25 10:18:38', NULL);



/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;