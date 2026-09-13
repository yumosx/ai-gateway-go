-- intent_rules: 意图识别规则表
-- MySQL DDL

CREATE TABLE IF NOT EXISTS `intent_rules` (
    `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `rule_id`      VARCHAR(64)     NOT NULL COMMENT '规则唯一ID',
    `intent_name`  VARCHAR(128)    NOT NULL COMMENT '意图名称',
    `patterns`     JSON            NOT NULL COMMENT '匹配模式列表, e.g. ["天气","下雨"]',
    `match_type`   VARCHAR(16)     NOT NULL DEFAULT 'contains' COMMENT 'exact/prefix/contains/regex',
    `params`       JSON            DEFAULT NULL COMMENT '规则携带的默认参数',
    `priority`     INT             NOT NULL DEFAULT 0 COMMENT '优先级, 越大越优先',
    `confidence`   DOUBLE          NOT NULL DEFAULT 1.0 COMMENT '该规则的固定置信度',
    `status`       TINYINT         NOT NULL DEFAULT 1 COMMENT '1=启用 0=禁用',
    `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_rule_id` (`rule_id`),
    KEY `idx_intent_name` (`intent_name`),
    KEY `idx_status_priority` (`status`, `priority` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='意图识别规则';

-- 种子数据: 常见意图规则
INSERT INTO `intent_rules` (`rule_id`, `intent_name`, `patterns`, `match_type`, `params`, `priority`, `confidence`, `status`) VALUES
('r001', 'greeting',       '["你好","hello","hi","嗨","早上好","下午好"]',                     'contains', NULL,                     10, 1.0,   1),
('r002', 'weather_query',  '["天气怎么样","今天天气","明天天气","天气预报"]',                    'contains', '{"date":"today"}',        20, 0.98,  1),
('r003', 'code_help',      '["帮我写代码","写一个","实现一个","代码实现"]',                      'contains', NULL,                     15, 0.95,  1),
('r004', 'translation',    '["翻译成","translate to","翻译一下"]',                             'contains', '{"target_lang":"auto"}',  15, 0.97,  1),
('r005', 'knowledge_qa',   '["什么是","解释一下","介绍一下","怎么理解"]',                       'contains', NULL,                     12, 0.90,  1),
('r006', 'date_query',     '["\\\\d{4}-\\\\d{2}-\\\\d{2}"]',                                  'regex',    NULL,                     25, 0.99,  1),
('r007', 'math_calculate', '["计算","算一下","等于多少","求和","加减乘除"]',                     'contains', NULL,                     18, 0.96,  1),
('r008', 'music_play',     '["播放音乐","放一首","听歌","来首歌"]',                             'contains', NULL,                     14, 0.95,  1),
('r009', 'navigation',     '["导航到","怎么去","路线","开车去","步行到"]',                      'contains', '{"mode":"driving"}',      22, 0.97,  1),
('r010', 'alarm_set',      '["设个闹钟","提醒我","定时","几分钟后提醒"]',                       'contains', NULL,                     16, 0.94,  1)
ON DUPLICATE KEY UPDATE `updated_at` = CURRENT_TIMESTAMP;
