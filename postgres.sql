create table Question (
	Id integer primary key,
	Title text not null,
	CreationDate integer not null,
	ClosedDate integer not null,
	Owner integer not null,
	Tags text
);
create table Answer (
	Id integer primary key,
	Owner integer not null,
	Question integer references Question(id),
	CreationDate integer not null,
	Accepted boolean not null,
	Score integer not null,
	SkillRep integer not null,
);
create index answer_question on Answer(Question);
create index answer_owner_accepted on Answer(Owner,Accepted);
create table Player (
	Id integer primary key,
	Name text,
	Profile text,
	SkillRep integer not null
);
create index player_skillrep on Player(SkillRep);

