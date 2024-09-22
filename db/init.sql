create table messages
(
    id         bigserial
        primary key,
    message_id uuid
        constraint uni_messages_message_id
            unique,
    content    varchar(50)      not null,
    "to"       varchar(25)      not null,
    status     bigint default 0 not null,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

-- insert multiple messages
insert into messages (message_id, content, "to", status) values
(null, 'Lorem ipsum dolor sit amet, consectetur adipiscing', '+90505323111', 0),
(null, 'Elit sed do eiusmod tempor incididunt ut labore et', '+90505323112', 0),
(null, 'Dolore magna aliqua. Ut enim ad minim veniam quis', '+90505323113', 0),
(null, 'Nostrud exercitation ullamco laboris nisi ut aliqu', '+90505323114', 0),
(null, 'Ex ea commodo consequat. Duis aute irure dolor in', '+90505323115', 0),
(null, 'Reprehenderit in voluptate velit esse cillum dolor', '+90505323116', 0),
(null, 'Eu fugiat nulla pariatur. Excepteur sint occaecat', '+90505323117', 0),
(null, 'Cupidatat non proident, sunt in culpa qui officia', '+90505323118', 0);




