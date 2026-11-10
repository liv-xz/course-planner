-- ============================================================
-- 选课规划智能体 - 初始化建库脚本
-- 数据库：MySQL 8.0+
-- 说明：18 张表，表名采用单数（与设计文档一致，GORM 用 SingularTable）
-- ============================================================

CREATE DATABASE IF NOT EXISTS course_planner
  DEFAULT CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE course_planner;

-- 1. 角色
CREATE TABLE IF NOT EXISTS `role` (
  id          BIGINT PRIMARY KEY AUTO_INCREMENT,
  name        VARCHAR(32) NOT NULL COMMENT 'student / admin',
  description VARCHAR(128) DEFAULT NULL,
  UNIQUE KEY uk_role_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色';

-- 2. 专业
CREATE TABLE IF NOT EXISTS major (
  id   BIGINT PRIMARY KEY AUTO_INCREMENT,
  code VARCHAR(32) NOT NULL,
  name VARCHAR(64) NOT NULL,
  UNIQUE KEY uk_major_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='专业';

-- 2.1 学期字典
CREATE TABLE IF NOT EXISTS semester (
  id            BIGINT PRIMARY KEY AUTO_INCREMENT,
  code          VARCHAR(16) NOT NULL COMMENT '如 2023-2024-1',
  academic_year VARCHAR(9) NOT NULL,
  term_no       TINYINT NOT NULL COMMENT '1/2',
  sequence_no   INT NOT NULL,
  display_name  VARCHAR(32) NOT NULL,
  UNIQUE KEY uk_semester_code (code),
  KEY idx_semester_seq (sequence_no)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='学期字典';

-- 3. 培养方案
CREATE TABLE IF NOT EXISTS training_plan (
  id            BIGINT PRIMARY KEY AUTO_INCREMENT,
  major_id      BIGINT NOT NULL,
  grade         VARCHAR(16) NOT NULL COMMENT '年级，如 2023',
  name          VARCHAR(128) NOT NULL,
  total_credits DECIMAL(5,1) DEFAULT NULL,
  created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_plan (major_id, grade),
  KEY idx_plan_major (major_id),
  CONSTRAINT fk_plan_major FOREIGN KEY (major_id) REFERENCES major(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='培养方案';

-- 4. 用户（只放通用字段）
CREATE TABLE IF NOT EXISTS `user` (
  id            BIGINT PRIMARY KEY AUTO_INCREMENT,
  username      VARCHAR(64) NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  role_id       BIGINT NOT NULL,
  name          VARCHAR(32) DEFAULT NULL,
  status        TINYINT NOT NULL DEFAULT 1 COMMENT '1启用 0禁用',
  is_deleted    TINYINT NOT NULL DEFAULT 0,
  created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at    DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_user_username (username),
  KEY idx_user_role (role_id),
  CONSTRAINT fk_user_role FOREIGN KEY (role_id) REFERENCES `role`(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户';

-- 5. 学生档案（学生专属字段，1:1 user）
CREATE TABLE IF NOT EXISTS student_profile (
  user_id    BIGINT PRIMARY KEY,
  major_id   BIGINT NOT NULL,
  student_no VARCHAR(32) NOT NULL,
  grade      VARCHAR(16) NOT NULL COMMENT '年级，如 2023',
  UNIQUE KEY uk_sp_student_no (student_no),
  KEY idx_sp_major (major_id),
  CONSTRAINT fk_sp_user FOREIGN KEY (user_id) REFERENCES `user`(id),
  CONSTRAINT fk_sp_major FOREIGN KEY (major_id) REFERENCES major(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='学生档案';

-- 6. 课程类别（含子模块）
CREATE TABLE IF NOT EXISTS category (
  id               BIGINT PRIMARY KEY AUTO_INCREMENT,
  plan_id          BIGINT NOT NULL,
  parent_id        BIGINT DEFAULT NULL,
  name             VARCHAR(64) NOT NULL,
  required_credits DECIMAL(5,1) NOT NULL DEFAULT 0,
  sort             INT NOT NULL DEFAULT 0,
  KEY idx_cat_plan (plan_id),
  KEY idx_cat_parent (parent_id),
  CONSTRAINT fk_cat_plan FOREIGN KEY (plan_id) REFERENCES training_plan(id),
  CONSTRAINT fk_cat_parent FOREIGN KEY (parent_id) REFERENCES category(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='课程类别';

-- 7. 课程库
CREATE TABLE IF NOT EXISTS course (
  id         BIGINT PRIMARY KEY AUTO_INCREMENT,
  code       VARCHAR(32) NOT NULL,
  name       VARCHAR(128) NOT NULL,
  credits    DECIMAL(4,1) NOT NULL,
  nature     TINYINT NOT NULL DEFAULT 1 COMMENT '1必修 2选修',
  is_deleted TINYINT NOT NULL DEFAULT 0,
  UNIQUE KEY uk_course_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='课程库';

-- 8. 课程-类别映射（N:M 中间表）
CREATE TABLE IF NOT EXISTS course_category (
  id          BIGINT PRIMARY KEY AUTO_INCREMENT,
  course_id   BIGINT NOT NULL,
  category_id BIGINT NOT NULL,
  UNIQUE KEY uk_cc (course_id, category_id),
  KEY idx_cc_category (category_id),
  CONSTRAINT fk_cc_course FOREIGN KEY (course_id) REFERENCES course(id),
  CONSTRAINT fk_cc_category FOREIGN KEY (category_id) REFERENCES category(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='课程-类别映射';

-- 9. 先修关系（课程自关联 N:M 中间表）
CREATE TABLE IF NOT EXISTS prerequisite (
  id               BIGINT PRIMARY KEY AUTO_INCREMENT,
  course_id        BIGINT NOT NULL COMMENT '本课',
  prereq_course_id BIGINT NOT NULL COMMENT '要先通过的课',
  UNIQUE KEY uk_prereq (course_id, prereq_course_id),
  KEY idx_prereq_course (course_id),
  CONSTRAINT fk_prereq_course FOREIGN KEY (course_id) REFERENCES course(id),
  CONSTRAINT fk_prereq_pre FOREIGN KEY (prereq_course_id) REFERENCES course(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='先修关系';

-- 10. 开课信息
CREATE TABLE IF NOT EXISTS course_offering (
  id        BIGINT PRIMARY KEY AUTO_INCREMENT,
  course_id BIGINT NOT NULL,
  semester  VARCHAR(16) NOT NULL COMMENT '日历学期，如 2023-2024-1',
  UNIQUE KEY uk_offering_course_semester (course_id, semester),
  KEY idx_offering_course (course_id),
  CONSTRAINT fk_offering_course FOREIGN KEY (course_id) REFERENCES course(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='开课信息';

-- 11. 学生修读记录
CREATE TABLE IF NOT EXISTS student_record (
  id         BIGINT PRIMARY KEY AUTO_INCREMENT,
  student_id BIGINT NOT NULL,
  course_id  BIGINT NOT NULL,
  score      DECIMAL(5,1) DEFAULT NULL,
  semester   VARCHAR(16) DEFAULT NULL,
  status     TINYINT NOT NULL DEFAULT 0 COMMENT '0在读 1通过 2未通过',
  credits    DECIMAL(4,1) DEFAULT NULL,
  KEY idx_record_student (student_id),
  KEY idx_record_course (course_id),
  CONSTRAINT fk_record_student FOREIGN KEY (student_id) REFERENCES `user`(id),
  CONSTRAINT fk_record_course FOREIGN KEY (course_id) REFERENCES course(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='学生修读记录';

-- 12. 学分审核结果
CREATE TABLE IF NOT EXISTS credit_audit (
  id                BIGINT PRIMARY KEY AUTO_INCREMENT,
  student_id        BIGINT NOT NULL,
  plan_id           BIGINT NOT NULL,
  category_id       BIGINT NOT NULL,
  completed_credits DECIMAL(5,1) NOT NULL DEFAULT 0,
  gap_credits       DECIMAL(5,1) NOT NULL DEFAULT 0,
  audited_at        DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_audit (student_id, plan_id, category_id),
  KEY idx_audit_student (student_id),
  CONSTRAINT fk_audit_student FOREIGN KEY (student_id) REFERENCES `user`(id),
  CONSTRAINT fk_audit_plan FOREIGN KEY (plan_id) REFERENCES training_plan(id),
  CONSTRAINT fk_audit_cat FOREIGN KEY (category_id) REFERENCES category(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='学分审核结果';

-- 13. 会话
CREATE TABLE IF NOT EXISTS chat_session (
  id         BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id    BIGINT NOT NULL,
  title      VARCHAR(128) DEFAULT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  KEY idx_session_user (user_id),
  CONSTRAINT fk_session_user FOREIGN KEY (user_id) REFERENCES `user`(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='智能体会话';

-- 14. 消息
CREATE TABLE IF NOT EXISTS chat_message (
  id         BIGINT PRIMARY KEY AUTO_INCREMENT,
  session_id BIGINT NOT NULL,
  role       VARCHAR(16) NOT NULL COMMENT 'user/assistant/tool',
  content    TEXT,
  tool_calls JSON DEFAULT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  KEY idx_msg_session (session_id),
  CONSTRAINT fk_msg_session FOREIGN KEY (session_id) REFERENCES chat_session(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会话消息';

-- 15. 拦截日志
CREATE TABLE IF NOT EXISTS security_log (
  id         BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id    BIGINT DEFAULT NULL,
  input      TEXT,
  reason     VARCHAR(255) DEFAULT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  KEY idx_seclog_user (user_id),
  CONSTRAINT fk_seclog_user FOREIGN KEY (user_id) REFERENCES `user`(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='安全拦截日志';

-- 16. 知识文档（RAG 元数据）
CREATE TABLE IF NOT EXISTS knowledge_doc (
  id          BIGINT PRIMARY KEY AUTO_INCREMENT,
  plan_id     BIGINT DEFAULT NULL,
  title       VARCHAR(128) NOT NULL,
  source      VARCHAR(255) DEFAULT NULL,
  chunk_count INT NOT NULL DEFAULT 0,
  created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
  KEY idx_kdoc_plan (plan_id),
  CONSTRAINT fk_kdoc_plan FOREIGN KEY (plan_id) REFERENCES training_plan(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='知识文档';

-- 17. 知识切片（RAG 检索）
CREATE TABLE IF NOT EXISTS knowledge_chunk (
  id          BIGINT PRIMARY KEY AUTO_INCREMENT,
  doc_id      BIGINT NOT NULL,
  chunk_index INT NOT NULL,
  title       VARCHAR(128) NOT NULL,
  content     TEXT NOT NULL,
  category    VARCHAR(64) DEFAULT NULL,
  source      VARCHAR(255) DEFAULT NULL,
  keywords    VARCHAR(255) DEFAULT NULL,
  created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
  KEY idx_chunk_doc (doc_id),
  CONSTRAINT fk_chunk_doc FOREIGN KEY (doc_id) REFERENCES knowledge_doc(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='知识切片';
