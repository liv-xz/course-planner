-- ============================================================
-- 学期字典与历史数据转换（已有数据库升级用）
-- 新库执行 001 + 002 即可；本文件可重复执行
-- ============================================================

USE course_planner;

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

INSERT IGNORE INTO semester (code, academic_year, term_no, sequence_no, display_name) VALUES
('2022-2023-1','2022-2023',1,1,'2022-2023学年第一学期'),
('2022-2023-2','2022-2023',2,2,'2022-2023学年第二学期'),
('2023-2024-1','2023-2024',1,3,'2023-2024学年第一学期'),
('2023-2024-2','2023-2024',2,4,'2023-2024学年第二学期'),
('2024-2025-1','2024-2025',1,5,'2024-2025学年第一学期'),
('2024-2025-2','2024-2025',2,6,'2024-2025学年第二学期'),
('2025-2026-1','2025-2026',1,7,'2025-2026学年第一学期'),
('2025-2026-2','2025-2026',2,8,'2025-2026学年第二学期'),
('2026-2027-1','2026-2027',1,9,'2026-2027学年第一学期'),
('2026-2027-2','2026-2027',2,10,'2026-2027学年第二学期'),
('2027-2028-1','2027-2028',1,11,'2027-2028学年第一学期'),
('2027-2028-2','2027-2028',2,12,'2027-2028学年第二学期');

UPDATE course_offering SET semester = CASE semester
  WHEN '第1学期' THEN '2023-2024-1'
  WHEN '第2学期' THEN '2023-2024-2'
  WHEN '第3学期' THEN '2024-2025-1'
  WHEN '第4学期' THEN '2024-2025-2'
  WHEN '第5学期' THEN '2025-2026-1'
  WHEN '第6学期' THEN '2025-2026-2'
  WHEN '第7学期' THEN '2026-2027-1'
  WHEN '第8学期' THEN '2026-2027-2'
  ELSE semester
END
WHERE semester LIKE '第_学期';

UPDATE student_record SET semester = CASE semester
  WHEN '第1学期' THEN '2023-2024-1'
  WHEN '第2学期' THEN '2023-2024-2'
  WHEN '第3学期' THEN '2024-2025-1'
  WHEN '第4学期' THEN '2024-2025-2'
  WHEN '第5学期' THEN '2025-2026-1'
  WHEN '第6学期' THEN '2025-2026-2'
  WHEN '第7学期' THEN '2026-2027-1'
  WHEN '第8学期' THEN '2026-2027-2'
  ELSE semester
END
WHERE semester LIKE '第_学期';
