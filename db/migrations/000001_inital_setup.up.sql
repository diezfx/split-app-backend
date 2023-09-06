create table if not exists projects(
    id UUID primary key,
    name text not null
);
create table if not exists members(
    id text primary key
);

create table if not exists project_memberships(
    project_id UUID ,
    user_id text,
    constraint fk_project_id
        foreign key(project_id) 
	    	references projects(id),
	constraint fk_user_id
    	foreign key(user_id) 
	    	references members(id),
    primary key(project_id, user_id)
);

create type transaction_type as enum (
  'Expense',
  'Transfer'
);

create table if not exists transactions(
	id UUID primary key,
	name text not null,
	amount INTEGER not null,
	source_id text not null,
	transaction_type transaction_type not null,
    project_id UUID not null,
    constraint fk_project_id
      foreign key(project_id) 
	  references projects(id),
    constraint fk_source_id
      foreign key(source_id) 
	  references members(id)
);

create table if not exists transaction_targets(
		transaction_id UUID,
		user_id text not null,
	  constraint fk_transaction_id
        foreign key(transaction_id) 
	    	references transactions(id),
	  constraint fk_user_id
        foreign key(user_id) 
	    	references members(id),
    primary key(transaction_id, user_id)
);
