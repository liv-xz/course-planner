-- ============================================================
-- 基础数据种子（不含 knowledge_doc）
-- 说明：参考数据用固定 ID + INSERT IGNORE，可重复执行
-- ============================================================

USE course_planner;

-- 0. 角色（student / admin）
INSERT IGNORE INTO `role` (name, description) VALUES
('student', '学生'),
('admin', '教务管理员');

-- 0.1 管理员用户（密码 123456，只写通用字段）
SET @admin_role = (SELECT id FROM `role` WHERE name = 'admin');
INSERT IGNORE INTO `user` (username, password_hash, role_id, name, status, is_deleted) VALUES
('admin', '$2a$10$8/DvdB6Sj.wT6yi9cZDHi.KJWbrj/z.JOnGPRVXbETA9xGRnKAYeW', @admin_role, '教务管理员', 1, 0);

-- 1. 专业
INSERT IGNORE INTO major (id, code, name) VALUES
(1, 'CS', '计算机科学与技术');

-- 1.1 学期字典（2022-2027）
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

-- 2. 培养方案
INSERT IGNORE INTO training_plan (id, major_id, grade, name, total_credits) VALUES
(1, 1, '2023', '2023 级计算机科学与技术培养方案', 165.0);

-- 3. 课程类别（4 大类 + 12 子模块）
INSERT IGNORE INTO category (id, plan_id, parent_id, name, required_credits, sort) VALUES
(1,  1, NULL, '通识教育课程',     51.0, 1),
(2,  1, 1,    '必修',             39.0, 1),
(3,  1, 1,    '体育选项课',         3.0, 2),
(4,  1, 1,    '综合素质选修课',      9.0, 3),
(5,  1, 4,    '创新创业类选修课',    1.0, 1),
(6,  1, 4,    '公共艺术类选修课',    2.0, 2),
(7,  1, 4,    '核心素养类选修课',    6.0, 3),
(8,  1, NULL, '学科基础课程',     43.0, 2),
(9,  1, 8,    '公共基础课',       27.0, 1),
(10, 1, 8,    '专业基础必修',     16.0, 2),
(11, 1, NULL, '专业课程',         39.0, 3),
(12, 1, 11,   '必修',             11.0, 1),
(13, 1, 11,   '选修',             28.0, 2),
(14, 1, NULL, '集中实践环节',     32.0, 4),
(15, 1, 14,   '综合实践',          6.0, 1),
(16, 1, 14,   '专业实践',         26.0, 2);

-- 4. 课程库（16 门样例）
INSERT IGNORE INTO course (id, code, name, credits, nature) VALUES
(1,  'CS101', '程序设计基础',   4.0, 1),
(2,  'CS102', '数据结构',       4.0, 1),
(3,  'CS103', '计算机组成原理', 4.0, 1),
(4,  'CS104', '操作系统',       4.0, 1),
(5,  'CS105', '计算机网络',     3.0, 1),
(6,  'CS201', '数据库原理',     3.0, 1),
(7,  'CS202', '软件工程',       3.0, 1),
(8,  'MA101', '高等数学',       5.0, 1),
(9,  'MA102', '线性代数',       3.0, 1),
(10, 'EN101', '大学英语',       4.0, 1),
(11, 'PE101', '体育',           1.0, 1),
(12, 'GE101', '创新创业基础',   1.0, 2),
(13, 'GE102', '公共艺术鉴赏',   2.0, 2),
(14, 'GE103', '核心素养导论',   2.0, 2),
(15, 'PR101', '专业实践',       2.0, 1),
(16, 'PR102', '综合实践',       1.0, 1);

-- 5. 学生用户（密码 123456）——只写通用字段
SET @student_role = (SELECT id FROM role WHERE name = 'student');
INSERT IGNORE INTO user (username, password_hash, role_id, name, status, is_deleted) VALUES
('2323500', '$2a$10$8/DvdB6Sj.wT6yi9cZDHi.KJWbrj/z.JOnGPRVXbETA9xGRnKAYeW', @student_role, '徐志强', 1, 0);

-- 6. 学生档案（学号、专业、年级）
INSERT IGNORE INTO student_profile (user_id, major_id, student_no, grade)
SELECT id, 1, '2323500', '2023' FROM user WHERE username = '2323500';

-- 7. 课程-类别映射
INSERT IGNORE INTO course_category (course_id, category_id) VALUES
(1, 10), (2, 10), (3, 10),
(4, 12), (5, 12), (6, 12), (7, 12),
(8, 9),  (9, 9),  (10, 9),
(11, 3),
(12, 5), (13, 6), (14, 7),
(15, 16), (16, 15);

-- 8. 先修关系
INSERT IGNORE INTO prerequisite (course_id, prereq_course_id) VALUES
(2, 1),
(4, 2),
(4, 3),
(6, 2);

-- 9. 开课信息（先清本 seed 范围，避免重复）
DELETE FROM course_offering WHERE course_id BETWEEN 1 AND 16;
INSERT INTO course_offering (course_id, semester) VALUES
(1,  '2023-2024-1'),
(8,  '2023-2024-1'),
(10, '2023-2024-1'),
(2,  '2024-2025-1'),
(3,  '2024-2025-1'),
(4,  '2024-2025-2'),
(5,  '2025-2026-1'),
(6,  '2024-2025-2'),
(7,  '2025-2026-1'),
(15, '2025-2026-2'),
(16, '2026-2027-1');

-- 10. 学生修读记录（先清该学生，避免重复）
SET @sid = (SELECT user_id FROM student_profile WHERE student_no = '2323500');
DELETE FROM student_record WHERE student_id = @sid;
INSERT INTO student_record (student_id, course_id, score, semester, status, credits) VALUES
(@sid, 1,  85.0, '2023-2024-1', 1, 4.0),
(@sid, 8,  90.0, '2023-2024-1', 1, 5.0),
(@sid, 9,  88.0, '2023-2024-2', 1, 3.0),
(@sid, 10, 82.0, '2023-2024-1', 1, 4.0),
(@sid, 11, 95.0, '2024-2025-1', 1, 1.0),
(@sid, 12, 90.0, '2024-2025-1', 1, 1.0),
(@sid, 13, 92.0, '2024-2025-2', 1, 2.0),
(@sid, 14, 89.0, '2024-2025-2', 1, 2.0);
