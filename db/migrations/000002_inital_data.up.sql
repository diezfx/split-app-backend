INSERT INTO projects (id, name)
VALUES ('f380b8ad-b0b2-4387-960b-0c107ce7f37e', 'testProject');

insert into members (id)
VALUES('user1'),('user2');

insert into project_memberships (project_id,user_id)
values('f380b8ad-b0b2-4387-960b-0c107ce7f37e','user1'),
('f380b8ad-b0b2-4387-960b-0c107ce7f37e','user2');

insert into transactions (id,name,amount,source_id,transaction_type,project_id)
values('79281b4f-b0a1-44cf-a2a0-2cf6cfa41faa','transaction1',146,'user1','Expense','f380b8ad-b0b2-4387-960b-0c107ce7f37e');